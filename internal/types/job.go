package types

import (
	"sync"
)

type JobStatus string

const (
	StatusQueued    JobStatus = "queued"
	StatusRunning   JobStatus = "running"
	StatusUploading JobStatus = "uploading"
	StatusDone      JobStatus = "done"
	StatusFailed    JobStatus = "failed"
)

type ContentType string

const (
	VideoContent ContentType = "video"
	AudioContent ContentType = "audio"
	MusicContent ContentType = "music"
)

type Job struct {
	ID       string
	URL      string
	Status   JobStatus
	Type     ContentType
	Progress float64
	Output   string
	Error    string
	// Uploaded is true once the finished file has been successfully stored in
	// Nextcloud. UploadError captures any WebDAV failure so the UI can offer a
	// retry without losing the already-downloaded file.
	Uploaded    bool
	UploadError string
	MU          sync.Mutex
}
