package storage

import (
	"context"
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

	_, err = s.DB.ExecContext(ctx,
		`INSERT INTO users (id, email, password_hash, created_at, updated_at)
     VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		id, email, hashedPassword,
	)

	slog.Info("admin user has been created.")
	return nil
}
