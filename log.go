package main

import (
	"log/slog"
	"os"
)

func initLogger(verbose bool) {
	level := slog.LevelError

	if verbose {
		level = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	slog.SetDefault(logger)
}
