package process

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rahulkumarpahwa/traer/internal/types"
)

// finalizeJob performs post-download work. For music jobs it uploads the MP3 to Nextcloud and only deletes the temporary file after a successful upload. For audio/video jobs it is a no-op (existing behavior is preserved).
func (jw *JobWorker) processUploadJob(job *types.Job) {
	if job.ProcessType != types.UploadProcess {
		return
	}

	if jw.Uploader == nil || !jw.Uploader.Configured() {
		fmt.Println("nextcloud is not configured")
		return
	}

	job.Mutex.Lock()
	output := job.Output
	job.Status = types.StatusUploading
	if err := jw.ProcessStorage.UpdateProcessByID(job.ID, types.StatusUploading, "", false, 0); err != nil {
		fmt.Printf(err.Error())
		return
	}
	job.Mutex.Unlock()

	if output == "" {
		fmt.Printf("no output file to upload")
		return
	}
	if _, err := os.Stat(output); err != nil {
		fmt.Printf("output file missing: %w", err)
		return
	}

	fileName := filepath.Base(output)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	if err := jw.Uploader.UploadFile(ctx, output, jw.Config.NextcloudBasePath, fileName); err != nil {
		job.Mutex.Lock()
		job.UploadError = err.Error()
		job.Mutex.Unlock()

		fmt.Printf("upload failed: %w", err)
		return
	}

	// Mark uploaded and delete the temporary file ONLY after a successful
	// upload. If deletion fails we still consider the job successful because the file safely lives in Nextcloud.
	job.Mutex.Lock()
	job.Uploaded = true
	job.UploadError = ""
	if err := jw.ProcessStorage.UpdateProcessByID(job.ID, types.StatusUploaded, "", true, 0); err != nil {
		fmt.Printf(err.Error())
		return
	}
	job.Mutex.Unlock()

	if err := os.Remove(output); err != nil {
		fmt.Printf("[job] uploaded %s but failed to delete temp file %s: %v\n", fileName, output, err)
	} else {
		fmt.Printf("[job] uploaded %s to Nextcloud and removed temp file\n", fileName)
	}
	return
}
