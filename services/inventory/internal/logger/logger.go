package logger

import (
	"log/slog"
	"os"
)

func New(serviceName string) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(handler).With(
		slog.String("service", serviceName),
	)

	return logger
}
