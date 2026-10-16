package storage

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/google/uuid"
	"github.com/rahulkumarpahwa/traer/config"
	"github.com/rahulkumarpahwa/traer/internal/auth"
)

func (s *UserStorage) CreateAdminUser(ctx context.Context, cfg *config.Config) error {

	id := uuid.New()
	email := cfg.AdminEmail
	password := cfg.AdminPassword
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return err
	}

	row := s.DB.QueryRowContext(ctx, `
    SELECT id FROM users WHERE email = ?`, email)

	if row.Err() != nil && row.Err() != sql.ErrNoRows{
		return row.Err()
	}

	var adminID string
	err = row.Scan(&adminID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if adminID == "" {
		_, err = s.DB.ExecContext(ctx,
			`INSERT INTO users (id, email, password_hash, created_at, updated_at)
     VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
			id, email, hashedPassword,
		)

		if err != nil {
			return err
		}
		slog.Info("admin user has been created.")

		return nil
	}

	slog.Info("admin user already exists.")
	return nil
}
