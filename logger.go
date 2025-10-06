package main

import (
	"log/slog"
	"os"
)

func newLogger(opts *slog.HandlerOptions) *slog.Logger {
	jsonHandler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(jsonHandler)
	return logger
}
