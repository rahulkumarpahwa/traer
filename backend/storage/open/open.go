package open

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
	if o.Config.Database.Path == "" {
		return nil, fmt.Errorf("Database Path not found!")
	}

	DB, err := sql.Open("sqlite", o.Config.Database.Path)
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
    id INTEGER PRIMARY KEY AUTOINCREMENT,  
	username TEXT NOT NULL UNIQUE,             
	email TEXT NOT NULL UNIQUE,                
	password_hash TEXT NOT NULL,           
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME   
	); `)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}
	slog.Info("User Table Created Successfully!")

	return DB, nil
}
