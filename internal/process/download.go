package process

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rahulkumarpahwa/traer/internal/types"
)

// runDownload executes yt-dlp for the job and records the resolved output path. It returns an error if the download fails; the caller records job state.
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

	outputWg.Wait()
	waitErr := cmd.Wait()
	if waitErr != nil {
		return waitErr
	}

	if resolvedPath, err := ResolveOutputPath(job.Output); err == nil && resolvedPath != "" {
		job.Output = resolvedPath
		job.TempPath = resolvedPath
		err := jw.ProcessStorage.UpdateProcessByID(job.ID, job.Status, job.TempPath, false, 0)
		if err != nil {
			return err
		}

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

	if job.ProcessType == types.AudioProcess {
		args := append([]string{
			"-f", "ba[ext=m4a]",
			"--embed-metadata",
			"--embed-thumbnail",
		}, baseArgs...)

		args = append(args, job.URL)
		return exec.Command(ytDlpPath, args...)
	}

	// mp3 version
	// if job.ProcessType == types.AudioProcess { args := append([]string{ "-f", "bestaudio", "-x", "--audio-format", "mp3", "--ffmpeg-location", ffmpegLocation, "--embed-thumbnail", "--add-metadata", }, baseArgs...) args = append(args, job.URL) return exec.Command(ytDlpPath, args...) }

	if job.ProcessType == types.VideoProcess {
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

// helpers
func parseOutput(pipe io.ReadCloser, job *types.Job) {
	scanner := bufio.NewScanner(pipe)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	var lastPathLine string // Capture the last path line to ensure we get the final filename

	for scanner.Scan() {
		line := scanner.Text()

		// Extract filepath - look for the --print after_move:filepath output
		// yt-dlp outputs just the filepath on its own line when using --print after_move:filepath
		trimmedLine := strings.TrimSpace(line)
		if (strings.HasPrefix(trimmedLine, "C:\\") || strings.HasPrefix(trimmedLine, "/")) &&
			strings.Contains(trimmedLine, "downloads") {
			// This looks like a filepath (Windows or Unix absolute path)
			// Store it, we'll use the last one encountered to ensure we get the final filename
			lastPathLine = trimmedLine
			fmt.Printf("job found path line: %s\n", trimmedLine)
		}
	}

	// Set the final output path after processing all output
	if lastPathLine != "" {
		job.Output = lastPathLine
		fmt.Printf("job Final captured output path: %s\n", lastPathLine)
		// job.Mutex.Unlock()
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading output:", err)
	}
}

// helpers
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
