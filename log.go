package main

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func initLogger(verbose bool) {
	level := slog.LevelError

	if verbose {
		level = slog.LevelDebug
	}

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		AddSource: false,
		Level:     level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.Attr{}
			}
			return a
		},
	}))

	slog.SetDefault(logger)
	slog.Debug("logger initialized", "level", level)
}
