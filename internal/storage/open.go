package storage

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/rahulkumarpahwa/traer/config"
	_ "modernc.org/sqlite"
)

type OpenDB struct {
	Config *config.Config
	Logger *slog.Logger
}

func (o *OpenDB) Open() (*sql.DB, error) {
	if o.Config.DatabasePath == "" {
		return nil, fmt.Errorf("Database Path not found!")
	}

	dsn := o.Config.DatabasePath
	// Ensure every pooled connection gets the same pragmas.
	// modernc.org/sqlite supports _pragma query params.
	sep := "?"
	if strings.Contains(dsn, "?") {
		sep = "&"
	}
	dsn += sep + "_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)&_pragma=synchronous(NORMAL)"

	DB, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	// SQLite has a single writer: serialize access through one connection.
	// 10 conns + busy_timeout=0 is what produced SQLITE_BUSY on Prepare.
	DB.SetMaxOpenConns(1)
	DB.SetMaxIdleConns(1)

	// ping the DB
	err = DB.Ping()
	if err != nil {
		return nil, err
	}

	// Belt-and-braces: also set on this connection in case DSN param
	// is ignored by driver version. Must be separate Execs - the driver
	// does not reliably run multi-statement strings in one Exec, and
	// the old code ignored the error return.
	for _, pragma := range []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA busy_timeout=5000",
		"PRAGMA foreign_keys=ON",
		"PRAGMA synchronous=NORMAL",
	} {
		if _, err := DB.Exec(pragma); err != nil {
			return nil, fmt.Errorf("failed to set %q: %w", pragma, err)
		}
	}

	// Opening the TABLE
	userStmt, err := DB.Prepare(`CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,                         
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		return nil, err
	}

	_, err = userStmt.Exec()
	_ = userStmt.Close()
	if err != nil {
		return nil, err
	}

	slog.Info("user table created successfully.")

	processStmt, err := DB.Prepare(`CREATE TABLE IF NOT EXISTS process (
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

	_, err = processStmt.Exec()
	_ = processStmt.Close()
	if err != nil {
		return nil, err
	}
	o.Logger.Info("process table created successfully.")

	return DB, nil
}
