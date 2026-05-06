package types

import (
	"sync"
	"time"
)

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
	MU       sync.Mutex
}

type User struct {
	ID           int        `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"password_hash"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}
