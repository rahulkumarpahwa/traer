package transcription

import "time"

type TranscriptionStatus string

const (
	StatusQueued  TranscriptionStatus = "queued"
	StatusRunning TranscriptionStatus = "running"
	StatusDone    TranscriptionStatus = "done"
	StatusFailed  TranscriptionStatus = "failed"
)

type TranscriptionJobs struct {
	ID             int                 `json:"id"`
	Status         TranscriptionStatus `json:"status"`
	SourceUrl      string              `json:"source_url"`
	AudioPath      string              `json:"audio_path"`
	TranscriptPath string              `json:"transcript_path"`
	CreatedAt      time.Time           `json:"created_at"`
	CompletedAt    time.Time           `json:"completed_at"`
}

type TranscriptSegments struct {
	ID        int       `json:"id"`
	JobID     int       `json:"job_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Text      string    `json:"text"`
}
