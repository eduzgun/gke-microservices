package logger

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	// Use JSON handler for structured logging
	// Logs are diagnoistics so I've chosen to log to Stderr
	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return slog.New(handler)
}

func NewDevLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}
