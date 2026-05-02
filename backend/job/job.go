package job

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/rahulkumarpahwa/traer/types"
	"github.com/rahulkumarpahwa/traer/utils"
)

var jobQueue = make(chan *types.Job, 100)
var progressRegex = regexp.MustCompile(`(\d+\.\d+)%`)

func AddJob(url string, contentType types.ContentType) *types.Job {
	job := &types.Job{
		ID:     utils.GenerateID(),
		URL:    url,
		Status: types.StatusQueued,
		Type:   contentType,
	}

	jobQueue <- job
	return job
}

func StartWorkers(n int) {
	for i := 0; i < n; i++ {
		go worker()
	}
}

func worker() {
	for job := range jobQueue {
		processJob(job)
	}
}

func processJob(job *types.Job) {
	job.Status = types.StatusRunning

	cmd := buildCommand(job)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		job.Status = types.StatusFailed
		job.Error = err.Error()
		return
	}

	// Read both stdout & stderr
	go parseOutput(stdout, job)
	go parseOutput(stderr, job)

	err := cmd.Wait()
	if err != nil {
		job.Status = types.StatusFailed
		job.Error = err.Error()
		return
	}

	job.Progress = 100
	job.Status = types.StatusDone
}

func buildCommand(job *types.Job) *exec.Cmd {
	baseArgs := []string{
		"--newline",
		"--progress",
		"--print", "after_move:filepath",
		"-o", "downloads/%(title)s.%(ext)s",
	}

	if job.Type == types.AudioContent {
		args := append([]string{
			"-f", "bestaudio",
			"-x",
			"--audio-format", "mp3",
			"--embed-thumbnail",
			"--add-metadata",
		}, baseArgs...)

		args = append(args, job.URL)
		return exec.Command("yt-dlp", args...)
	}

	if job.Type == types.VideoContent {
		args := append([]string{
			"-f", "bv*+ba/b",
			"--merge-output-format", "mp4",
			"--embed-thumbnail",
			"--add-metadata",
		}, baseArgs...)

		args = append(args, job.URL)
		return exec.Command("yt-dlp", args...)
	}

	return nil
}

func parseOutput(pipe io.ReadCloser, job *types.Job) {
	scanner := bufio.NewScanner(pipe)

	for scanner.Scan() {
		line := scanner.Text()

		match := progressRegex.FindStringSubmatch(line)
		if len(match) > 1 {
			progress, err := strconv.ParseFloat(match[1], 64)
			if err == nil {
				job.MU.Lock()
				job.Progress = progress
				job.MU.Unlock()
			}
		}
	}


    if err := scanner.Err(); err != nil {
        // handle or log the error
        fmt.Println("Error reading output:", err)
    }
}
