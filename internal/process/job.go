package process

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rahulkumarpahwa/traer/config"
	"github.com/rahulkumarpahwa/traer/internal/nextcloud"
	"github.com/rahulkumarpahwa/traer/internal/types"
	"github.com/rahulkumarpahwa/traer/internal/utils"
)

type JobWorker struct {
	startOnce       sync.Once
	Config          *config.Config
	ExecutablePaths map[string]string
	// Uploader streams finished music files to Nextcloud over WebDAV. It is nil
	// when Nextcloud credentials are not configured, in which case music jobs
	// fail fast with a clear error.
	Uploader *nextcloud.Client
	JobsMu   sync.RWMutex
}

var jobQueue = make(chan *types.Job, 100)
var progressRegex = regexp.MustCompile(`(\d+(?:\.\d+)?)%`)

var Jobs = make(map[string]*types.Job) // to store the jobs and fetch later with the id.

func (jw *JobWorker) AddJob(url string, contentType types.ContentType) *types.Job {
	// Prevent accidental duplicate downloads: if an identical URL + type is
	// already queued, running, uploading, or done, reuse that job instead of
	// starting a new one. Failed jobs are excluded so a retry can be issued.
	jw.JobsMu.RLock()
	for _, existing := range Jobs {
		existing.MU.Lock()
		dup := existing.URL == url && existing.Type == contentType &&
			(existing.Status == types.StatusQueued ||
				existing.Status == types.StatusRunning ||
				existing.Status == types.StatusUploading ||
				existing.Status == types.StatusDone)
		existing.MU.Unlock()
		if dup {
			jw.JobsMu.RUnlock()
			return existing
		}
	}
	jw.JobsMu.RUnlock()

	job := &types.Job{
		ID:     utils.GenerateID(),
		URL:    url,
		Status: types.StatusQueued,
		Type:   contentType,
	}
	jw.JobsMu.Lock()
	Jobs[job.ID] = job
	jw.JobsMu.Unlock()

	// selection in case queue fulls
	select {
	case jobQueue <- job:
		// ok
	default:
		job.MU.Lock()
		job.Status = types.StatusFailed
		job.Error = "queue full"
		job.MU.Unlock()
	}

	return job
}

// RetryJob re-queues an existing failed job. If the media file is already on
// disk (download succeeded but the upload failed), the job resumes at the
// upload step instead of re-downloading.
func (jw *JobWorker) RetryJob(id string) (*types.Job, error) {
	jw.JobsMu.RLock()
	job := Jobs[id]
	jw.JobsMu.RUnlock()

	if job == nil {
		return nil, fmt.Errorf("job not found")
	}

	job.MU.Lock()
	// Only failed (or queued-but-stuck) jobs are eligible for retry.
	if job.Status != types.StatusFailed {
		status := job.Status
		job.MU.Unlock()
		if status == types.StatusDone || status == types.StatusUploading {
			return nil, fmt.Errorf("job is already %s", status)
		}
		return nil, fmt.Errorf("job is %s and cannot be retried", status)
	}
	job.Error = ""
	job.UploadError = ""
	job.Status = types.StatusQueued
	job.MU.Unlock()

	select {
	case jobQueue <- job:
		// ok
	default:
		job.MU.Lock()
		job.Status = types.StatusFailed
		job.Error = "queue full"
		job.MU.Unlock()
		return nil, fmt.Errorf("queue full")
	}

	return job, nil
}

func (jw *JobWorker) GetJob(id string) *types.Job {
	jw.JobsMu.RLock()
	job := Jobs[id]
	jw.JobsMu.RUnlock()
	return job
}

func (jw *JobWorker) StartWorkers() {
	jw.startOnce.Do(func() {
		if err := os.MkdirAll("downloads", os.ModePerm); err != nil {
			panic(err)
		}

		for i := 0; i < jw.Config.Workers.Instances; i++ {
			go jw.worker()
		}
	})
}

func (jw *JobWorker) worker() {
	for job := range jobQueue {
		jw.processJob(job)
	}
}

func (jw *JobWorker) processJob(job *types.Job) {
	job.MU.Lock()
	job.Status = types.StatusRunning
	job.MU.Unlock()

	// Resume support: for music jobs that already have a downloaded file on
	// disk (e.g. a previous attempt failed during upload), skip the download
	// and go straight to uploading so we never re-download unnecessarily.
	resumeFromDisk := false
	if job.Type == types.MusicContent {
		job.MU.Lock()
		existingOutput := job.Output
		job.MU.Unlock()
		if existingOutput != "" {
			if resolved, err := ResolveOutputPath(existingOutput); err == nil && resolved != "" {
				job.MU.Lock()
				job.Output = resolved
				job.MU.Unlock()
				resumeFromDisk = true
			}
		}
	}

	if !resumeFromDisk {
		if err := jw.runDownload(job); err != nil {
			job.MU.Lock()
			job.Status = types.StatusFailed
			job.Error = err.Error()
			job.MU.Unlock()
			return
		}
	}

	if err := jw.finalizeJob(job); err != nil {
		job.MU.Lock()
		job.Status = types.StatusFailed
		job.Error = err.Error()
		job.MU.Unlock()
		return
	}

	job.MU.Lock()
	job.Progress = 100
	job.Status = types.StatusDone
	fmt.Printf("[job] Job %s done. Output: %s\n", job.ID, job.Output)
	job.MU.Unlock()
}

// runDownload executes yt-dlp for the job and records the resolved output path.
// It returns an error if the download fails; the caller records job state.
func (jw *JobWorker) runDownload(job *types.Job) error {
	cmd := jw.buildCommand(job)
	if cmd == nil {
		return fmt.Errorf("invalid job type")
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	var outputWg sync.WaitGroup
	outputWg.Add(2)
	go func() {
		defer outputWg.Done()
		parseOutput(stdout, job)
	}()
	go func() {
		defer outputWg.Done()
		parseOutput(stderr, job)
	}()

	waitErr := cmd.Wait()
	outputWg.Wait()
	if waitErr != nil {
		return waitErr
	}

	if resolvedPath, err := ResolveOutputPath(job.Output); err == nil && resolvedPath != "" {
		job.MU.Lock()
		job.Output = resolvedPath
		job.MU.Unlock()
	}
	return nil
}

// finalizeJob performs post-download work. For music jobs it uploads the MP3 to
// Nextcloud and only deletes the temporary file after a successful upload. For
// audio/video jobs it is a no-op (existing behavior is preserved).
func (jw *JobWorker) finalizeJob(job *types.Job) error {
	if job.Type != types.MusicContent {
		return nil
	}

	if jw.Uploader == nil || !jw.Uploader.Configured() {
		return fmt.Errorf("nextcloud is not configured")
	}

	job.MU.Lock()
	output := job.Output
	job.Status = types.StatusUploading
	job.MU.Unlock()

	if output == "" {
		return fmt.Errorf("no output file to upload")
	}
	if _, err := os.Stat(output); err != nil {
		return fmt.Errorf("output file missing: %w", err)
	}

	fileName := filepath.Base(output)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	if err := jw.Uploader.UploadFile(ctx, output, jw.Config.Nextcloud.BasePath, fileName); err != nil {
		job.MU.Lock()
		job.UploadError = err.Error()
		job.MU.Unlock()
		return fmt.Errorf("upload failed: %w", err)
	}

	// Mark uploaded and delete the temporary file ONLY after a successful
	// upload. If deletion fails we still consider the job successful because the
	// file safely lives in Nextcloud.
	job.MU.Lock()
	job.Uploaded = true
	job.UploadError = ""
	job.MU.Unlock()

	if err := os.Remove(output); err != nil {
		fmt.Printf("[job] uploaded %s but failed to delete temp file %s: %v\n", fileName, output, err)
	} else {
		fmt.Printf("[job] uploaded %s to Nextcloud and removed temp file\n", fileName)
	}
	return nil
}

func (jw *JobWorker) buildCommand(job *types.Job) *exec.Cmd {
	baseArgs := []string{
		"--newline",
		"--progress",
		"--print", "after_move:filepath",
		"-o", "downloads/%(title)s.%(ext)s",
	}

	ytDlpPath := jw.ExecutablePaths["yt-dlp"]
	if ytDlpPath == "" {
		return nil
	}

	ffmpegPath := jw.ExecutablePaths["ffmpeg"]
	if ffmpegPath == "" {
		return nil
	}

	ffmpegLocation := filepath.Dir(ffmpegPath)

	if job.Type == types.AudioContent || job.Type == types.MusicContent {
		args := append([]string{
			"-f", "bestaudio",
			"-x",
			"--audio-format", "mp3",
			"--ffmpeg-location", ffmpegLocation,
			"--embed-thumbnail",
			"--add-metadata",
		}, baseArgs...)

		args = append(args, job.URL)
		return exec.Command(ytDlpPath, args...)
	}

	if job.Type == types.VideoContent {
		args := append([]string{
			"-f", "bv*+ba/best",
			"--merge-output-format", "mp4",
			"--remux-video", "mp4",
			"--ffmpeg-location", ffmpegLocation,
			"--embed-thumbnail",
			"--add-metadata",
		}, baseArgs...)

		args = append(args, job.URL)
		return exec.Command(ytDlpPath, args...)
	}

	return nil
}

func parseOutput(pipe io.ReadCloser, job *types.Job) {
	scanner := bufio.NewScanner(pipe)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	var lastPathLine string // Capture the last path line to ensure we get the final filename

	for scanner.Scan() {
		line := scanner.Text()

		match := progressRegex.FindStringSubmatch(line)
		if len(match) > 1 {
			progress, err := strconv.ParseFloat(match[1], 64)
			if err == nil {
				job.MU.Lock()
				if progress > job.Progress {
					job.Progress = progress
				}
				job.MU.Unlock()
			}
		}

		// Extract filepath - look for the --print after_move:filepath output
		// yt-dlp outputs just the filepath on its own line when using --print after_move:filepath
		trimmedLine := strings.TrimSpace(line)
		if (strings.HasPrefix(trimmedLine, "C:\\") || strings.HasPrefix(trimmedLine, "/")) &&
			strings.Contains(trimmedLine, "downloads") {
			// This looks like a filepath (Windows or Unix absolute path)
			// Store it, we'll use the last one encountered to ensure we get the final filename
			lastPathLine = trimmedLine
			fmt.Printf("[job] Found path line: %s\n", trimmedLine)
		}
	}

	// Set the final output path after processing all output
	if lastPathLine != "" {
		job.MU.Lock()
		job.Output = lastPathLine
		fmt.Printf("[job] Final captured output path: %s\n", lastPathLine)
		job.MU.Unlock()
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading output:", err)
	}
}

func ResolveOutputPath(outputPath string) (string, error) {
	if outputPath == "" {
		return "", nil
	}

	if _, err := os.Stat(outputPath); err == nil {
		return outputPath, nil
	}

	dir := filepath.Dir(outputPath)
	base := filepath.Base(outputPath)
	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)
	stem = strings.TrimSpace(stem)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		candidateName := entry.Name()
		candidateExt := filepath.Ext(candidateName)
		if ext != "" && !strings.EqualFold(candidateExt, ext) {
			continue
		}

		candidateStem := strings.TrimSpace(strings.TrimSuffix(candidateName, candidateExt))
		if strings.HasPrefix(candidateStem, stem) || strings.HasPrefix(stem, candidateStem) {
			return filepath.Join(dir, candidateName), nil
		}
	}

	return "", os.ErrNotExist
}
