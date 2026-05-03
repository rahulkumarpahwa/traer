package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rahulkumarpahwa/traer/job"
	"github.com/rahulkumarpahwa/traer/types"
	"github.com/rahulkumarpahwa/traer/utils"
)

type ServiceHandler struct {
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

	job := job.AddJob(req.URL, contentType)

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

func (hs *ServiceHandler) HandleSetInstances(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Count int `json:"count"`
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if req.Count <= 0 {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("count must be > 0"))
		return
	}

	job.WorkerInstances = req.Count

	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"message":   "instances updated",
		"instances": job.WorkerInstances,
	})
}

func (hs *ServiceHandler) HandleGetInstances(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"instances": job.WorkerInstances,
	})
}