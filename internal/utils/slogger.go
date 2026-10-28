package utils

import (
	"log/slog"
	"os"
)

func NewSlogger(writer *os.File) *slog.Logger {

	slogHandler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		// AddSource: true,
		Level:     slog.LevelInfo, // current and upper level will be logged.
		// Normally we use the level debug only when we need some debuging on the project and info level is set for the production.
	})
	Logger := slog.New(slogHandler)
	slog.SetDefault(Logger)
	return Logger
}
