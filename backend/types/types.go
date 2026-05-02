package types

import "sync"

type JobStatus string

const (
	StatusQueued  JobStatus = "queued"
	StatusRunning JobStatus = "running"
	StatusDone    JobStatus = "done"
	StatusFailed  JobStatus = "failed"
)

type ContentType string

const (
	VideoContent ContentType = "video"
	AudioContent ContentType = "audio"
)

type Job struct {
	ID       string
	URL      string
	Status   JobStatus
	Type     ContentType
	Progress float64
	Output   string
	Error    string
	mu       sync.Mutex
}
