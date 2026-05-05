package job

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/rahulkumarpahwa/traer/config"
	"github.com/rahulkumarpahwa/traer/types"
	"github.com/rahulkumarpahwa/traer/utils"
)

type JobWorker struct {
	startOnce       sync.Once
	Config          *config.Config
	ExecutablePaths map[string]string
}

var jobQueue = make(chan *types.Job, 100)
var progressRegex = regexp.MustCompile(`(\d+(?:\.\d+)?)%`)

var Jobs = make(map[string]*types.Job) // to store the jobs and fetch later with the id.
var JobsMu sync.RWMutex

func (jw *JobWorker) AddJob(url string, contentType types.ContentType) *types.Job {
	job := &types.Job{
		ID:     utils.GenerateID(),
		URL:    url,
		Status: types.StatusQueued,
		Type:   contentType,
	}
	JobsMu.Lock()
	Jobs[job.ID] = job
	JobsMu.Unlock()

	// selection in case queue fulls
	select {
	case jobQueue <- job:
		// ok
	default:
		job.Status = types.StatusFailed
		job.Error = "queue full"
	}

	return job
}

func GetJob(id string) *types.Job {
	JobsMu.RLock()
	job := Jobs[id]
	JobsMu.RUnlock()
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

	cmd := jw.buildCommand(job)
	if cmd == nil {
		job.MU.Lock()
		job.Status = types.StatusFailed
		job.Error = "invalid job type"
		job.MU.Unlock()
		return
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		job.MU.Lock()
		job.Status = types.StatusFailed
		job.Error = err.Error()
		job.MU.Unlock()
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		job.MU.Lock()
		job.Status = types.StatusFailed
		job.Error = err.Error()
		job.MU.Unlock()
		return
	}

	if err := cmd.Start(); err != nil {
		job.MU.Lock()
		job.Status = types.StatusFailed
		job.Error = err.Error()
		job.MU.Unlock()
		return
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

	err = cmd.Wait()
	if err != nil {
		job.MU.Lock()
		job.Status = types.StatusFailed
		job.Error = err.Error()
		job.MU.Unlock()
		return
	}

	outputWg.Wait()

	if resolvedPath, err := ResolveOutputPath(job.Output); err == nil && resolvedPath != "" {
		job.MU.Lock()
		job.Output = resolvedPath
		job.MU.Unlock()
	}

	job.MU.Lock()
	job.Progress = 100
	job.Status = types.StatusDone
	fmt.Printf("[job] Job %s done. Output: %s\n", job.ID, job.Output)
	job.MU.Unlock()
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

	if job.Type == types.AudioContent {
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
