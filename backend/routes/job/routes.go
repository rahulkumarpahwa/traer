package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/rahulkumarpahwa/traer/job"
	"github.com/rahulkumarpahwa/traer/storage"
	jTypes "github.com/rahulkumarpahwa/traer/types/job"
	"github.com/rahulkumarpahwa/traer/utils"
	"github.com/rahulkumarpahwa/traer/validators"
)

type ServiceHandler struct {
	JW          *job.JobWorker
	UserStorage storage.UserStorageInterface
}

func (hs *ServiceHandler) HandleCreateJobs(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		URL  string   `json:"url"`
		URLs []string `json:"urls"` // optional batch of URLs
		Type string   `json:"type"` // "audio", "video" or "music"
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	// Normalize to a list of URLs, accepting either a single "url" or a batch
	// in "urls". Blank entries are dropped and duplicates within the request
	// are collapsed while preserving order.
	var urls []string
	seen := map[string]bool{}
	for _, u := range append([]string{req.URL}, req.URLs...) {
		u = strings.TrimSpace(u)
		if u == "" || seen[u] {
			continue
		}
		seen[u] = true
		urls = append(urls, u)
	}

	if len(urls) == 0 {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("url is required"))
		return
	}

	var contentType jTypes.ContentType
	switch strings.ToLower(req.Type) {
	case "audio":
		contentType = jTypes.AudioContent
	case "video":
		contentType = jTypes.VideoContent
	case "music":
		contentType = jTypes.MusicContent
		if hs.JW.Uploader == nil || !hs.JW.Uploader.Configured() {
			utils.ErrorResponse(w, http.StatusServiceUnavailable, fmt.Errorf("nextcloud upload is not configured"))
			return
		}
	default:
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid type"))
		return
	}

	for _, u := range urls {
		if !validators.IsYouTubeURL(u) {
			utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("Invalid URL: %s", u))
			return
		}
	}

	created := make([]*jTypes.Job, 0, len(urls))
	for _, u := range urls {
		created = append(created, hs.JW.AddJob(u, contentType))
	}

	// Return the first created job for backward compatibility with existing
	// single-URL clients, plus the full list under "jobs" for batch callers.
	if len(created) == 1 {
		utils.JSONResponse(w, http.StatusOK, created[0])
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"id":     created[0].ID,
		"status": created[0].Status,
		"jobs":   created,
	})
}

func (hs *ServiceHandler) HandleActiveJobs(w http.ResponseWriter, r *http.Request) {
	var activeJobs []*jTypes.Job

	hs.JW.JobsMu.RLock()
	for _, j := range job.Jobs {
		j.MU.Lock()
		if j.Status == jTypes.StatusQueued || j.Status == jTypes.StatusRunning || j.Status == jTypes.StatusUploading {
			activeJobs = append(activeJobs, j)
		}
		j.MU.Unlock()
	}
	hs.JW.JobsMu.RUnlock()

	utils.JSONResponse(w, http.StatusOK, activeJobs)
}

func (hs *ServiceHandler) HandleGetInstances(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"instances": hs.JW.Config.Workers.Instances,
	})
}

func (hs *ServiceHandler) HandleJobStatus(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("id is required"))
		return
	}

	j := hs.JW.GetJob(id)
	if j == nil {
		utils.ErrorResponse(w, http.StatusNotFound, fmt.Errorf("job not found"))
		return
	}

	j.MU.Lock()
	response := map[string]interface{}{
		"id":          j.ID,
		"url":         j.URL,
		"type":        j.Type,
		"status":      j.Status,
		"progress":    j.Progress,
		"output":      j.Output,
		"error":       j.Error,
		"uploaded":    j.Uploaded,
		"uploadError": j.UploadError,
	}
	j.MU.Unlock()

	utils.JSONResponse(w, http.StatusOK, response)
}

// HandleRetryJob re-queues a previously failed job. If the media file is
// already on disk (download succeeded but the upload failed), the job resumes
// at the upload step instead of re-downloading.
func (hs *ServiceHandler) HandleRetryJob(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("id is required"))
		return
	}

	j, err := hs.JW.RetryJob(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, j)
}

func (hs *ServiceHandler) HandleDownload(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("id is required"))
		return
	}

	fmt.Printf("[download] Attempting to download job with id: %s\n", id)

	j := hs.JW.GetJob(id)
	if j == nil {
		fmt.Printf("[download] Job not found with id: %s\n", id)
		hs.JW.JobsMu.RLock()
		fmt.Printf("[download] Available jobs: %v\n", len(job.Jobs))
		for jid := range job.Jobs {
			fmt.Printf("[download]   - %s\n", jid)
		}
		hs.JW.JobsMu.RUnlock()
		utils.ErrorResponse(w, http.StatusNotFound, fmt.Errorf("job not found"))
		return
	}

	j.MU.Lock()
	status := j.Status
	outputPath := j.Output
	j.MU.Unlock()

	fmt.Printf("[download] Job found: status=%s, output=%s\n", status, outputPath)

	if status != jTypes.StatusDone {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("file not ready"))
		return
	}

	if outputPath == "" {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("output path missing"))
		return
	}

	resolvedPath, err := job.ResolveOutputPath(outputPath)
	if err == nil && resolvedPath != "" {
		outputPath = resolvedPath
		j.MU.Lock()
		j.Output = resolvedPath
		j.MU.Unlock()
	}

	if _, err := os.Stat(outputPath); err != nil {
		fmt.Printf("[download] File stat error: %v\n", err)
		utils.ErrorResponse(w, http.StatusNotFound, fmt.Errorf("file not found: %v", err))
		return
	}

	// Serve file
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(outputPath))
	w.Header().Set("Content-Type", "application/octet-stream")

	http.ServeFile(w, r, outputPath)
}
