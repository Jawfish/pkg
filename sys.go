package main

import (
	"log/slog"
	"os"
	"os/exec"
	"runtime"
)

// getArch returns Fedora's architecture string for the current system's architecture.
//
// Since Fedora packages are only available for a limited number of architectures, this
// function returns an error if Fedora packages are not available for the current
// architecture.
func getArch() (string, error) {
	slog.Debug("getting architecture")

	// maps the architecture string from Go to the Fedora architecture string
	archMap := map[string]string{
		"i386":  "i686",
		"amd64": "x86_64",
		"arm64": "aarch64",
		"arm":   "armv7hl",
	}

	arch := runtime.GOARCH

	if archMap[arch] == "" {
		return "", &ErrUnsupportedArch{Arch: arch}
	}

	slog.Debug("found supported architecture", "goArch", arch, "fedoraArch", archMap[arch])

	return archMap[arch], nil
}

// getPackageManager returns the package manager to use for the current system by
// checking if dnf or dnf5 is available in the PATH. dnf5 is preferred over dnf if both
// are available. If neither are available, an error is returned.
func getPackageManager() (string, error) {
	slog.Debug("getting package manager")
	_, errDnf := exec.LookPath("dnf")
	_, errDnf5 := exec.LookPath("dnf5")

	if errDnf5 == nil {
		slog.Debug("found dnf5")
		return "dnf5", nil
	} else if errDnf == nil {
		slog.Debug("found dnf")
		return "dnf", nil
	} else {
		return "", &ErrPackageManagerNotFound{}
	}
}

// checkCache returns nil if the dnf package database cache exists and is a valid
// package database with the expected schema, otherwise an error is returned.
func checkCache(path string) error {
	slog.Debug("checking if cache exists")
	_, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			slog.Warn("cache does not exist")
			return &ErrCacheNotFound{Path: path}

		} else {
			slog.Error("error checking if cache exists", "err", err)
			return err
		}
	}
	slog.Debug("cache found, checking schema")

	err = ensureValidSchema(path)
	if err != nil {
		slog.Error("error ensuring valid schema", "err", err)
		return err
	}

	slog.Debug("cache has valid schema")

	return nil
}

// generateCache updates the dnf package database cache using the dnf update --refresh
// command.
func generateCache(dnfBinary string) error {
	slog.Debug("generating cache")
	out, err := exec.Command(dnfBinary, "makecache").Output()

	if err != nil {
		slog.Error("error generating cache", "err", err)
		return err
	}

	slog.Info("cache generated", "output", string(out))
	return nil
}
