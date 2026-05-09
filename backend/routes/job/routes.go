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

	var contentType jTypes.ContentType

	switch strings.ToLower(req.Type) {
	case "audio":
		contentType = jTypes.AudioContent
	case "video":
		contentType = jTypes.VideoContent
	default:
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid type"))
		return
	}

	job := hs.JW.AddJob(req.URL, contentType)

	utils.JSONResponse(w, http.StatusOK, job)
}

func (hs *ServiceHandler) HandleActiveJobs(w http.ResponseWriter, r *http.Request) {
	var activeJobs []*jTypes.Job

	hs.JW.JobsMu.RLock()
	for _, j := range job.Jobs {
		j.MU.Lock()
		if j.Status == jTypes.StatusQueued || j.Status == jTypes.StatusRunning {
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
