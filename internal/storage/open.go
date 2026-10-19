package storage

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/rahulkumarpahwa/traer/config"
	_ "modernc.org/sqlite"
)

type OpenDB struct {
	Config *config.Config
}

func (o *OpenDB) Open() (*sql.DB, error) {
	if o.Config.DatabasePath == "" {
		return nil, fmt.Errorf("Database Path not found!")
	}

	DB, err := sql.Open("sqlite", o.Config.DatabasePath)
	if err != nil {
		return nil, err
	}

	// ping the DB
	err = DB.Ping()
	if err != nil {
		return nil, err
	}

	// Opening the TABLE
	stmt, err := DB.Prepare(`CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,                         
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}
	slog.Info("user table created successfully.")

	stmt, err = DB.Prepare(`CREATE TABLE IF NOT EXISTS process (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id),
    url TEXT,
    status TEXT NOT NULL DEFAULT 'not_started'
    CHECK (status IN ('not_started', 'queued', 'running', 'uploading', 'uploaded', 'completed', 'failed')),
    process_type TEXT DEFAULT NULL
    CHECK (process_type IN ('audio', 'video', 'transcription', 'playlist', 'upload') OR process_type IS NULL),
	uploaded INTEGER NOT NULL DEFAULT 0 CHECK (uploaded IN (0, 1)),
	retries INT DEFAULT 0,
	temp_path TEXT DEFAULT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)

	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}
	slog.Info("process table created successfully.")

	_, err = DB.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, err
	}

	return DB, nil
}
