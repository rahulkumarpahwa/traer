package types

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type JobStatus string

const (
	StatusNotStarted    JobStatus = "not_started"
	StatusQueued    JobStatus = "queued"
	StatusRunning   JobStatus = "running"
	StatusUploading JobStatus = "uploading"
	StatusUploaded  JobStatus = "uploaded"
	StatusCompleted JobStatus = "completed"
	StatusFailed    JobStatus = "failed"
)

type ProcessType string

const (
	VideoProcess         ProcessType = "video"
	AudioProcess         ProcessType = "audio"
	PlaylistProcess      ProcessType = "playlist"
	TranscriptionProcess ProcessType = "transcription"
	UploadProcess        ProcessType = "upload"
)

type Job struct {
	ID          uuid.UUID
	URL         string // single link or a playlist link
	Status      JobStatus
	ProcessType ProcessType
	// Progress    float64
	TempPath    string
	Output      string
	Error       []string
	Retries     int
	Uploaded    bool
	UploadError string
	Mutex       sync.Mutex
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
