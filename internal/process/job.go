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
)

type JobWorker struct {
	startOnce       sync.Once
	ProcessStorage  storage.ProcessStorageInterface
	Config          *config.Config
	ExecutablePaths map[string]string
	Uploader        *nextcloud.Client
}

var jobQueue = make(chan *types.Job, 200)

func (jw *JobWorker) StartWorkers() {
	jw.startOnce.Do(func() {
		for i := 0; i < jw.Config.WorkerInstances; i++ {
			go jw.worker()
		}
	})
}

func (jw *JobWorker) worker() {
	for job := range jobQueue {
		if job.ProcessType != types.UploadProcess {
		// 	jw.processUploadJob(job)
		// } else {
			jw.processDownloadJob(job)
		}
	}
}

// adds the job to the channel
func (jw *JobWorker) AddNewJob(url string, processType types.ProcessType, userId uuid.UUID) error {
	job := &types.Job{
		ID:          uuid.New(),
		URL:         url,
		Status:      types.StatusQueued,
		ProcessType: processType,
	}

	if err := jw.ProcessStorage.CreateProcess(job.ID, userId, job.URL, job.Status, job.ProcessType); err != nil {
		return err
	}

	// selection in case queue fulls
	select {
	case jobQueue <- job:
		fmt.Printf("new process added in queue of type: %s\n", processType)
	default:
		job.Status = types.StatusFailed
		job.Error = append(job.Error, "queue full")
		if err := jw.ProcessStorage.UpdateProcessByID(job.ID, types.StatusFailed, "", false, 0); err != nil {
			return err
		}
	}

	return nil
}

// downloads the process job temporarily
func (jw *JobWorker) processDownloadJob(job *types.Job) {
	if job.ProcessType == types.UploadProcess {
		return
	}

	job.Status = types.StatusRunning
	if err := jw.ProcessStorage.UpdateProcessByID(job.ID, types.StatusRunning, "", false, 0); err != nil {
		fmt.Println(err)
		return
	}

	if err := jw.runDownload(job); err != nil {
		job.Status = types.StatusFailed
		job.Error = append(job.Error, err.Error())
		return
	}

	job.Status = types.StatusCompleted
	if err := jw.ProcessStorage.UpdateProcessByID(job.ID, types.StatusCompleted, "", false, 0); err != nil {
		fmt.Println(err)
		return
	}
	slog.Info("job has been completed", "job_id", job.ID.String(), "process_type", job.ProcessType)

	// adding the job again to upload:
	job.ProcessType = types.UploadProcess
	jobQueue <- job
	fmt.Printf("new upload process added in queue of type: %s\n", job.ProcessType)
}
