package logger

import (
	"log/slog"
	"os"
)

func InitLogger() *slog.Logger {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})

	logger := slog.New(logHandler)

	return logger
}
