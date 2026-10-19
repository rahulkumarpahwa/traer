package storage

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rahulkumarpahwa/traer/internal/tdd"
	"github.com/rahulkumarpahwa/traer/internal/types"
	"github.com/rahulkumarpahwa/traer/internal/validators"
)

type ProcessStorage struct {
	DB *sql.DB
}

type ProcessStorageInterface interface {
	GetProcessByID(id uuid.UUID) (tdd.GetProcess, error)
	CreateProcess(id uuid.UUID, user_id uuid.UUID, url string, status types.JobStatus, processType types.ProcessType) error
	UpdateProcessByID(id uuid.UUID, status types.JobStatus, tempPath string, uploaded bool, retries int) error
	// DeleteUserByID(id int) error
}

func (p *ProcessStorage) GetProcessByID(id uuid.UUID) (tdd.GetProcess, error) {
	if id == uuid.Nil {
		return tdd.GetProcess{}, fmt.Errorf("id should not be invalid")
	}

	stmt, err := p.DB.Prepare(`SELECT id, user_id, url, status, process_type, retries, uploaded, created_at, updated_at FROM process WHERE id = ?`)
	if err != nil {
		return tdd.GetProcess{}, fmt.Errorf("failed to prepare select: %v", err)
	}
	defer stmt.Close()

	var process tdd.GetProcess
	err = stmt.QueryRow(id).Scan(&process.ID, &process.UserID, &process.Url, &process.Status, &process.Processtype, &process.Retries, &process.Uploaded, &process.CreatedAt, &process.UpdatedAt)
	if err == sql.ErrNoRows {
		return tdd.GetProcess{}, err // No user found
	} else if err != nil {
		return tdd.GetProcess{}, fmt.Errorf("failed to query user: %v", err)
	}

	return process, nil
}

func (p *ProcessStorage) CreateProcess(id uuid.UUID, user_id uuid.UUID, url string, status types.JobStatus, processType types.ProcessType) error {
	// Validate inputs
	if id == uuid.Nil || user_id == uuid.Nil {
		return fmt.Errorf("id and user_id should not be empty")
	}

	if !validators.IsYouTubeURL(url) {
		return fmt.Errorf("url should not be invalidated")
	}

	// put these in valiadtors

	// Prepare the SQL statement
	stmt, err := p.DB.Prepare(`
		INSERT INTO process (id, user_id, url, status, process_type)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close() // Always close prepared statements

	// Execute the prepared statement
	_, err = stmt.Exec(id.String(), user_id.String(), url, status, processType)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %v", err)
	}

	return nil
}

// upgrade this method to update multiple values.
func (p *ProcessStorage) UpdateProcessByID(id uuid.UUID, status types.JobStatus, tempPath string, uploaded bool, retries int) error {
	if id == uuid.Nil {
		return fmt.Errorf("process id should not be invalid.")
	}

	queryParts := []string{}
	args := []any{}

	if status != "" {
		queryParts = append(queryParts, "status = ?")
		args = append(args, status)
	}

	if tempPath != "" {
		queryParts = append(queryParts, "temp_path = ?")
		args = append(args, tempPath)
	}

	if uploaded {
		queryParts = append(queryParts, "uploaded = ?")
		args = append(args, 1)
	}

	if !(retries < 1) {
		queryParts = append(queryParts, "retries = ?")
		args = append(args, retries)
	}

	// Build final query
	query := fmt.Sprintf("UPDATE process SET %s updated_at = ? WHERE id = ?", strings.Join(queryParts, ", "))

	// Prepare and execute
	stmt, err := p.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare update: %v", err)
	}
	defer stmt.Close()

	args = append(args, time.Now())
	args = append(args, id.String())

	_, err = stmt.Exec(args...)
	if err != nil {
		return fmt.Errorf("failed to execute process update: %v", err)
	}

	return nil
}

// // DeleteUser removes a user by ID
// func (u *UserStorage) DeleteUserByID(id int) error {
// 	if id <= 0 {
// 		return fmt.Errorf("invalid user ID")
// 	}

// 	stmt, err := u.DB.Prepare(`DELETE FROM users WHERE id = ?`)
// 	if err != nil {
// 		return fmt.Errorf("failed to prepare delete: %v", err)
// 	}
// 	defer stmt.Close()

// 	_, err = stmt.Exec(id)
// 	if err != nil {
// 		return fmt.Errorf("failed to execute delete: %v", err)
// 	}

// 	return nil
// }
