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

// Example logs
// slog.Info("User login successful",
//     slog.String("user_id", "123"),
//     slog.Int("status", 200),
//     slog.Bool("admin", true),
//     slog.Float64("duration_ms", 12.34),
//     slog.Time("login_time", time.Now()),
//     slog.Group("request",
//         slog.String("method", "POST"),
//         slog.String("path", "/login"),
//     ),
// )
