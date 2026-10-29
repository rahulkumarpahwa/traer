package process

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/rahulkumarpahwa/traer/config"
	"github.com/rahulkumarpahwa/traer/internal/nextcloud"
	"github.com/rahulkumarpahwa/traer/internal/storage"
	"github.com/rahulkumarpahwa/traer/internal/types"
	"github.com/rahulkumarpahwa/traer/internal/utils"
)

type JobWorker struct {
	startOnce       sync.Once
	ProcessStorage  storage.ProcessStorageInterface
	Config          *config.Config
	ExecutablePaths map[string]string
	Uploader        *nextcloud.Client
	Logger          *slog.Logger
}

var jobQueue = make(chan *types.Job, utils.GetWorkerQueueSize())

func (jw *JobWorker) StartWorkers() {
	jw.startOnce.Do(func() {
		for i := 0; i < jw.Config.WorkerInstances; i++ {
			go jw.worker()
		}
	})
}

func (jw *JobWorker) worker() {
	for job := range jobQueue {
		if job.ProcessType == types.ProcessUpload {
			jw.processUploadJob(job)
		} else if job.ProcessType == types.ProcessDownload {
			jw.processDownloadJob(job)
		}
	}
}

// adds the job to the channel
func (jw *JobWorker) AddNewJob(url string, contentType types.ContentType, audioQuality types.AudioQuality, videoQuality types.VideoQuality, processType types.ProcessType, userID uuid.UUID) error {

	var status types.JobStatus
	if processType == types.ProcessDownload {
		status = types.StatusDownloadQueued
	} else if processType == types.ProcessUpload {
		status = types.StatusUploadQueued
	}

	job := &types.Job{
		ID:           uuid.New(),
		URL:          url,
		ContentType:  contentType,
		Status:       status,
		ProcessType:  processType,
		AudioQuality: audioQuality,
		VideoQuality: videoQuality,
	}

	if err := jw.ProcessStorage.CreateProcess(job.ID, userID, job.URL, job.Status, job.ProcessType); err != nil {
		return err
	}

	// selection in case queue fulls
	select {
	case jobQueue <- job:
		jw.Logger.Info("new job added ", "id", job.ID)
	default:
		job.Error = "queue full"
		if processType == types.ProcessDownload {
			job.Status = types.StatusDownloadFailed
			if err := jw.ProcessStorage.UpdateProcessByID(job.ID, types.StatusDownloadFailed, 0); err != nil {
				return err
			}
		} else if processType == types.ProcessUpload {
			job.Status = types.StatusUploadFailed
			job.Status = types.StatusDownloadFailed
			if err := jw.ProcessStorage.UpdateProcessByID(job.ID, types.StatusUploadFailed, 0); err != nil {
				return err
			}
		}
	}

	return nil
}

// downloads the process job temporarily
func (jw *JobWorker) processDownloadJob(job *types.Job) {
	if job.ProcessType == types.ProcessUpload {
		return
	}

	job.Status = types.StatusDownloadRunning
	if err := jw.ProcessStorage.UpdateProcessByID(job.ID, types.StatusDownloadRunning, 0); err != nil {
		fmt.Println(err)
		return
	}

	if err := jw.runDownload(job); err != nil {
		job.Status = types.StatusDownloadFailed
		job.Error = err.Error()
		return
	}

	job.Status = types.StatusDownloadCompleted
	if err := jw.ProcessStorage.UpdateProcessByID(job.ID, types.StatusDownloadCompleted, 0); err != nil {
		fmt.Println(err)
		return
	}
	slog.Info("job has been completed", "job_id", job.ID.String(), "process_type", job.ProcessType)

	// adding the job again to upload:
	job.ProcessType = types.ProcessUpload
	job.Status = types.StatusUploadQueued
	jobQueue <- job
	if err := jw.ProcessStorage.UpdateProcessByID(job.ID, types.StatusUploadQueued, 0); err != nil {
		fmt.Printf(err.Error())
		return
	}
	fmt.Printf("new upload process added in queue of type: %s\n", job.ProcessType)
}
