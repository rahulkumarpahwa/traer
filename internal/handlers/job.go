package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rahulkumarpahwa/traer/internal/auth"
	"github.com/rahulkumarpahwa/traer/internal/dto"
	"github.com/rahulkumarpahwa/traer/internal/process"
	"github.com/rahulkumarpahwa/traer/internal/service"
	"github.com/rahulkumarpahwa/traer/internal/storage"
	"github.com/rahulkumarpahwa/traer/internal/types"
	"github.com/rahulkumarpahwa/traer/internal/utils"
	"github.com/rahulkumarpahwa/traer/internal/validators"
)

type ServiceHandler struct {
	JW          *process.JobWorker
	PL          *service.PlaylistWorker
	UserStorage storage.UserStorageInterface
}

func (hs *ServiceHandler) HandleCreateJobs(w http.ResponseWriter, r *http.Request) {
	var createJobReq dto.CreateJobRequest

	userID, err := auth.GetUserId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&createJobReq); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	urls, err := validators.FilterAndValidateUrls(createJobReq.URLs)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	var contentType types.ContentType

	switch createJobReq.ContentType {
	case "audio":
		contentType = types.ContentAudio
	case "video":
		contentType = types.ContentVideo
	default:
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid content_type\n"))
		return
	}

	var processType types.ProcessType

	switch createJobReq.ProcessType {
	case "download":
		processType = types.ProcessDownload
	case "upload":
		processType = types.ProcessUpload
	default:
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid process_type\n"))
		return
	}

	var audioQuality types.AudioQuality
	if createJobReq.AudioQuality != "" {
		switch createJobReq.AudioQuality {
		case "mp3":
			audioQuality = types.AudioMp3
		case "m4a":
			audioQuality = types.AudioM4A

		default:
			utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid audio_quality\n"))
			return
		}
	}

	var videoQuality types.VideoQuality
	if createJobReq.VideoQuality != "" {

		switch createJobReq.VideoQuality {
		case "360p":
			videoQuality = types.Video360p
		case "480p":
			videoQuality = types.Video480p
		case "720p":
			videoQuality = types.Video720p
		case "1080p":
			videoQuality = types.Video1080p
		case "1440p":
			videoQuality = types.Video1440p
		case "2160p":
			videoQuality = types.Video2160p
		default:
			utils.ErrorResponse(
				w,
				http.StatusBadRequest,
				fmt.Errorf("invalid video_quality"),
			)
			return
		}
	}

	for _, u := range urls {

		ok, _ := validators.IsYouTubePlaylist(u)

		if ok {
			playlist := &types.Playlist{
				PlaylistUrl: u,
				ProcessType: processType,
			}
			if err := hs.PL.RunPlayList(playlist); err != nil {
				utils.ErrorResponse(w, http.StatusBadGateway, err)
				return
			}

			for _, entry := range playlist.Entries {
				if err := hs.JW.AddNewJob(entry.Url, contentType, audioQuality, videoQuality, processType, userID); err != nil {
					utils.ErrorResponse(w, http.StatusBadRequest, err)
					return
				}
			}

		} else {
			if err := hs.JW.AddNewJob(u, contentType, audioQuality, videoQuality, processType, userID); err != nil {
				utils.ErrorResponse(w, http.StatusBadRequest, err)
				return
			}
		}

	}

	utils.JSONResponse(w, http.StatusOK, "process has been created\n")
}

func (hs *ServiceHandler) HandleJobs(w http.ResponseWriter, r *http.Request) {
	var jobs []dto.GetProcess

	jobs, err := hs.JW.ProcessStorage.GetProcesses()
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, jobs)
}

// func (hs *ServiceHandler) HandleGetInstances(w http.ResponseWriter, r *http.Request) {
// 	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
// 		"instances": hs.JW.Config.WorkerInstances,
// 	})
// }

// func (hs *ServiceHandler) HandleJobStatus(w http.ResponseWriter, r *http.Request) {
// 	id := r.URL.Query().Get("id")

// 	if id == "" {
// 		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("id is required"))
// 		return
// 	}

// 	j := hs.JW.GetJob(id)
// 	if j == nil {
// 		utils.ErrorResponse(w, http.StatusNotFound, fmt.Errorf("job not found"))
// 		return
// 	}

// 	j.MU.Lock()
// 	response := map[string]interface{}{
// 		"id":          j.ID,
// 		"url":         j.URL,
// 		"type":        j.Type,
// 		"status":      j.Status,
// 		"progress":    j.Progress,
// 		"output":      j.Output,
// 		"error":       j.Error,
// 		"uploaded":    j.Uploaded,
// 		"uploadError": j.UploadError,
// 	}
// 	j.MU.Unlock()

// 	utils.JSONResponse(w, http.StatusOK, response)
// }

// // HandleRetryJob re-queues a previously failed job. If the media file is
// // already on disk (download succeeded but the upload failed), the job resumes
// // at the upload step instead of re-downloading.
// func (hs *ServiceHandler) HandleRetryJob(w http.ResponseWriter, r *http.Request) {
// 	id := r.URL.Query().Get("id")
// 	if id == "" {
// 		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("id is required"))
// 		return
// 	}

// 	j, err := hs.JW.RetryJob(id)
// 	if err != nil {
// 		utils.ErrorResponse(w, http.StatusBadRequest, err)
// 		return
// 	}

// 	utils.JSONResponse(w, http.StatusOK, j)
// }

// func (hs *ServiceHandler) HandleDownload(w http.ResponseWriter, r *http.Request) {
// 	id := r.URL.Query().Get("id")

// 	if id == "" {
// 		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("id is required"))
// 		return
// 	}

// 	fmt.Printf("[download] Attempting to download job with id: %s\n", id)

// 	j := hs.JW.GetJob(id)
// 	if j == nil {
// 		fmt.Printf("[download] Job not found with id: %s\n", id)
// 		hs.JW.JobsMu.RLock()
// 		fmt.Printf("[download] Available jobs: %v\n", len(process.Jobs))
// 		for jid := range process.Jobs {
// 			fmt.Printf("[download]   - %s\n", jid)
// 		}
// 		hs.JW.JobsMu.RUnlock()
// 		utils.ErrorResponse(w, http.StatusNotFound, fmt.Errorf("job not found"))
// 		return
// 	}

// 	j.MU.Lock()
// 	status := j.Status
// 	outputPath := j.Output
// 	j.MU.Unlock()

// 	fmt.Printf("[download] Job found: status=%s, output=%s\n", status, outputPath)

// 	if status != types.StatusDone {
// 		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("file not ready"))
// 		return
// 	}

// 	if outputPath == "" {
// 		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("output path missing"))
// 		return
// 	}

// 	resolvedPath, err := process.ResolveOutputPath(outputPath)
// 	if err == nil && resolvedPath != "" {
// 		outputPath = resolvedPath
// 		j.MU.Lock()
// 		j.Output = resolvedPath
// 		j.MU.Unlock()
// 	}

// 	if _, err := os.Stat(outputPath); err != nil {
// 		fmt.Printf("[download] File stat error: %v\n", err)
// 		utils.ErrorResponse(w, http.StatusNotFound, fmt.Errorf("file not found: %v", err))
// 		return
// 	}

// 	// Serve file
// 	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(outputPath))
// 	w.Header().Set("Content-Type", "application/octet-stream")

// 	http.ServeFile(w, r, outputPath)
// }
