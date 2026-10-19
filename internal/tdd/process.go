package tdd

import (
	"time"

	"github.com/google/uuid"
	"github.com/rahulkumarpahwa/traer/internal/types"
)

type GetProcess struct {
	ID          uuid.UUID         `json:"id"`
	UserID      uuid.UUID         `json:"user_id"`
	Url         string            `json:"url"`
	Status      types.JobStatus   `json:"status"`
	Processtype types.ProcessType `json:"process_type"`
	Retries     int               `json:"retries"`
	Uploaded    bool              `json:"uploaded"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}
