package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/rahulkumarpahwa/traer/job"
	"github.com/rahulkumarpahwa/traer/types"
	"github.com/rahulkumarpahwa/traer/utils"
)

type ServiceHandler struct {
	JW                  *job.JobWorker
	MissingDependencies []string
}

func (hs *ServiceHandler) HandleCreateJobs(w http.ResponseWriter, r *http.Request) {
	if len(hs.MissingDependencies) > 0 {
		utils.ErrorResponse(w, http.StatusServiceUnavailable, fmt.Errorf("missing dependencies: %s", strings.Join(hs.MissingDependencies, ", ")))
		return
	}

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

	j := job.GetJob(id)
	if j == nil {
		utils.ErrorResponse(w, http.StatusNotFound, fmt.Errorf("job not found"))
		return
	}

	j.MU.Lock()
	defer j.MU.Unlock()

	if j.Status != types.StatusDone {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("file not ready"))
		return
	}

	if j.Output == "" {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("output path missing"))
		return
	}

	// Serve file
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(j.Output))
	w.Header().Set("Content-Type", "application/octet-stream")

	http.ServeFile(w, r, j.Output)
}
