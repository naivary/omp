package logger

import (
	"log/slog"
	"os"
)

func New(opts *slog.HandlerOptions) *slog.Logger {
	jsonHandler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(jsonHandler)
	return logger
}
