package types

import (
	"github.com/google/uuid"
	"time"
)

type AudioQuality string

const (
	AudioMp3 AudioQuality = "mp3"
	AudioM4A AudioQuality = "m4a"
)

type VideoQuality string

const (
	Video360p  VideoQuality = "360p"
	Video480p  VideoQuality = "480p"
	Video720p  VideoQuality = "720p"
	Video1080p VideoQuality = "1080p"
	Video1440p VideoQuality = "1440p"
	Video2160p VideoQuality = "2160p"
)

type JobStatus string

const (
	StatusNotStarted        JobStatus = "not_started"
	StatusDownloadQueued    JobStatus = "download_queued"
	StatusDownloadRunning   JobStatus = "download_running"
	StatusDownloadCompleted JobStatus = "download_completed"
	StatusDownloadFailed    JobStatus = "download_failed"
	StatusUploadQueued      JobStatus = "upload_queued"
	StatusUploadRunning     JobStatus = "upload_running"
	StatusUploadCompleted   JobStatus = "upload_completed"
	StatusUploadFailed      JobStatus = "upload_failed"
)

type ContentType string

const (
	ContentAudio ContentType = "audio"
	ContentVideo ContentType = "video"
)

type ProcessType string

const (
	ProcessDownload ProcessType = "download"
	ProcessUpload   ProcessType = "upload"
)

type Job struct {
	ID           uuid.UUID
	URL          string
	ContentType  ContentType
	Status       JobStatus
	ProcessType  ProcessType
	AudioQuality AudioQuality
	VideoQuality VideoQuality
	OutputPath   string
	Error        string
	Retries      int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
