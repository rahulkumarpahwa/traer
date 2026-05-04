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
	startOnce sync.Once
	Config    *config.Config
}

var jobQueue = make(chan *types.Job, 100)
var progressRegex = regexp.MustCompile(`(\d+(?:\.\d+)?)%`)

var Jobs = make(map[string]*types.Job) // to store the jobs and fetch later with the id.
var JobsMu sync.RWMutex

var WorkerInstances int = 2

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

		for i := 0; i < WorkerInstances; i++ {
			go jw.worker()
		}
	})
}

func (jw * JobWorker) worker() {
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

	// Read both stdout & stderr
	go parseOutput(stdout, job)
	go parseOutput(stderr, job)

	err = cmd.Wait()
	if err != nil {
		job.MU.Lock()
		job.Status = types.StatusFailed
		job.Error = err.Error()
		job.MU.Unlock()
		return
	}

	job.MU.Lock()
	job.Progress = 100
	job.Status = types.StatusDone
	job.MU.Unlock()
}

func (jw * JobWorker) buildCommand(job *types.Job) *exec.Cmd {
	baseArgs := []string{
		"--newline",
		"--progress",
		"--print", "after_move:filepath",
		"-o", "downloads/%(title)s.%(ext)s",
	}

	ytDlpPath := jw.Config.RequiredExecutables["yt-dlp"]
	if ytDlpPath == "" {
		return nil
	}

	ffmpegPath := jw.Config.RequiredExecutables["ffmpeg"]
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
			"-f", "bv*+ba/b",
			"--merge-output-format", "mp4",
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

		if strings.Contains(line, "downloads/") || strings.Contains(line, "downloads\\") {
			job.MU.Lock()
			job.Output = strings.TrimSpace(line)
			job.MU.Unlock()
		}
	}

	if err := scanner.Err(); err != nil {
		// handle or log the error
		fmt.Println("Error reading output:", err)
	}
}
