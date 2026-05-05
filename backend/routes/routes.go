package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/rahulkumarpahwa/traer/job"
	"github.com/rahulkumarpahwa/traer/types"
	"github.com/rahulkumarpahwa/traer/utils"
	"github.com/rahulkumarpahwa/traer/validators"
)

type ServiceHandler struct {
	JW *job.JobWorker
}

func (hs *ServiceHandler) HandleCreateJobs(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		URL  string `json:"url"`
		Type string `json:"type"` // "audio" or "video"
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if req.URL == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("url is required"))
		return
	}

	if !validators.IsYouTubeURL(req.URL) {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("Invalid URL"))
		return
	}

	var contentType types.ContentType

	switch strings.ToLower(req.Type) {
	case "audio":
		contentType = types.AudioContent
	case "video":
		contentType = types.VideoContent
	default:
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid type"))
		return
	}

	job := hs.JW.AddJob(req.URL, contentType)

	utils.JSONResponse(w, http.StatusOK, job)
}

func (hs *ServiceHandler) HandleActiveJobs(w http.ResponseWriter, r *http.Request) {
	var activeJobs []*types.Job

	job.JobsMu.RLock()
	for _, j := range job.Jobs {
		j.MU.Lock()
		if j.Status == types.StatusQueued || j.Status == types.StatusRunning {
			activeJobs = append(activeJobs, j)
		}
		j.MU.Unlock()
	}
	job.JobsMu.RUnlock()

	utils.JSONResponse(w, http.StatusOK, activeJobs)
}

func (hs *ServiceHandler) HandleGetInstances(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"instances": job.WorkerInstances,
	})
}

func (hs *ServiceHandler) HandleJobStatus(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("id is required"))
		return
	}

	j := job.GetJob(id)
	if j == nil {
		utils.ErrorResponse(w, http.StatusNotFound, fmt.Errorf("job not found"))
		return
	}

	j.MU.Lock()
	response := map[string]interface{}{
		"id":       j.ID,
		"url":      j.URL,
		"type":     j.Type,
		"status":   j.Status,
		"progress": j.Progress,
		"output":   j.Output,
		"error":    j.Error,
	}
	j.MU.Unlock()

	utils.JSONResponse(w, http.StatusOK, response)
}

func (hs *ServiceHandler) HandleDownload(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("id is required"))
		return
	}

	fmt.Printf("[download] Attempting to download job with id: %s\n", id)

	j := job.GetJob(id)
	if j == nil {
		fmt.Printf("[download] Job not found with id: %s\n", id)
		job.JobsMu.RLock()
		fmt.Printf("[download] Available jobs: %v\n", len(job.Jobs))
		for jid := range job.Jobs {
			fmt.Printf("[download]   - %s\n", jid)
		}
		job.JobsMu.RUnlock()
		utils.ErrorResponse(w, http.StatusNotFound, fmt.Errorf("job not found"))
		return
	}

	j.MU.Lock()
	status := j.Status
	outputPath := j.Output
	j.MU.Unlock()

	fmt.Printf("[download] Job found: status=%s, output=%s\n", status, outputPath)

	if status != types.StatusDone {
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
