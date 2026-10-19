package process

import (
	"fmt"
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
	JobsMu          sync.RWMutex
	JobsWG          sync.WaitGroup
}

var jobQueue = make(chan *types.Job, 10)

func (jw *JobWorker) StartWorkers() {
	jw.startOnce.Do(func() {
		for i := 0; i < jw.Config.WorkerInstances; i++ {
			jw.JobsWG.Add(1)
			go jw.worker()
		}
	})
	jw.JobsWG.Wait()
}

func (jw *JobWorker) worker() {
	defer jw.JobsWG.Done()
	for job := range jobQueue {
		if job.ProcessType == types.UploadProcess {
			jw.processUploadJob(job)
		} else {
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
		fmt.Printf("new process added in queue of type: %s", processType)
	default:
		job.Mutex.Lock()
		job.Status = types.StatusFailed
		job.Error = append(job.Error, "queue full")
		if err := jw.ProcessStorage.UpdateProcessByID(job.ID, types.StatusFailed, "", false, 0); err != nil {
			return err
		}
		job.Mutex.Unlock()
	}

	return nil
}

// downloads the process job temporarily
func (jw *JobWorker) processDownloadJob(job *types.Job) {
	if job.ProcessType == types.UploadProcess {
		return
	}

	job.Mutex.Lock()
	job.Status = types.StatusRunning
	if err := jw.ProcessStorage.UpdateProcessByID(job.ID, types.StatusRunning, "", false, 0); err != nil {
		fmt.Println(err)
		return
	}
	job.Mutex.Unlock()

	if err := jw.runDownload(job); err != nil {
		job.Mutex.Lock()
		job.Status = types.StatusFailed
		job.Error = append(job.Error, err.Error())
		job.Mutex.Unlock()
		return
	}

	job.Mutex.Lock()
	job.Status = types.StatusCompleted
	if err := jw.ProcessStorage.UpdateProcessByID(job.ID, types.StatusCompleted, "", false, 0); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("job [%s] has been completed with process type: %s", job.ID.String(), job.ProcessType)
	job.Mutex.Unlock()

	// // adding the job again to upload: selection in case queue fulls
	// select {
	// case jobQueue <- job:
	// 	fmt.Printf("new upload process added in queue of type: %s", job.ProcessType)
	// default:
	// 	job.Mutex.Lock()
	// 	job.Status = types.StatusFailed
	// 	job.Error = append(job.Error, "queue full")
	// 	if err := jw.ProcessStorage.UpdateProcessByID(job.ID, types.StatusFailed, "", false, 0); err != nil {
	// 		fmt.Printf("updating the process failed: %v with err: %v", job.ID, err)
	// 		return
	// 	}
	// 	job.Mutex.Unlock()
	// }

}
