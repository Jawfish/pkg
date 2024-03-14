package main

import (
	"errors"
	"flag"
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

func ensureValidArch() {
	_, err := getArch()

	if err != nil {
		slog.Error("error checking architecture", "err", err)
		os.Exit(1)
	}
}

func ensureCache(dnfBinary string, path string) {
	err := checkCache(path)
	if err == nil {
		return
	}

	var e *ErrCacheNotFound
	if !errors.As(err, &e) {
		slog.Error("error checking cache", "err", err)
		os.Exit(1)
	}

	slog.Warn("cache not found, generating")
	if err = generateCache(dnfBinary); err != nil {
		slog.Error("error generating cache", "err", err)
		os.Exit(1)
	}

	if err = checkCache(path); err != nil {
		if errors.As(err, &e) {
			slog.Error("cache still can't be found after generating. If you are using the -c flag, ensure the path matches the configuration for DNF.")
		} else {
			slog.Error("error checking cache", "err", err)
		}
		os.Exit(1)
	}
}

func setFlags(noConfirm *bool, cachePath *string, verbose *bool) {
	flag.BoolVar(noConfirm, "y", false, "skip the confirmation prompt")
	flag.StringVar(cachePath, "c", "/var/cache/dnf/packages.db", "the path to the DNF cache database")
	flag.BoolVar(verbose, "v", false, "show verbose output")

	flag.Parse()

}

func main() {
	var (
		noConfirm bool
		cachePath string
		verbose   bool
	)
	setFlags(&noConfirm, &cachePath, &verbose)

	initLogger(verbose)

	dnfBinary, err := getPackageManager()
	if err != nil {
		slog.Error("error getting package manager", "err", err)
		os.Exit(1)
	}

	ensureValidArch()
	ensureCache(dnfBinary, cachePath)
}
