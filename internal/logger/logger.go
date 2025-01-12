package customlogger

import (
	"log/slog"
	"os"
)

func NewLogger(logLevel slog.Level, env string) *slog.Logger {
	var logHandler slog.Handler

	switch env {
	case "development":
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		})
	default:
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     logLevel,
		})

	}
	logger := slog.New(logHandler)
	return logger
}
