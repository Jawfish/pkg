package main

import (
	"log/slog"
	"os"
	"os/exec"
	"runtime"
)

// getArch returns Fedora's architecture string for the current system's architecture
// which is used when querying the package database to ensure only viable packages are
// returned.
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
		slog.Error("unsupported architecture", "goArch", arch)
		return "", &ErrUnsupportedArch{Arch: arch}
	} else {
		slog.Info("found supported architecture", "goArch", arch, "fedoraArch", archMap[arch])
	}

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
		slog.Error("no package manager found")
		return "", &ErrPackageManagerNotFound{}
	}
}

// depsInstalled checks if the required dependencies are installed: dnf/dnf5 and fzf.
// If any dependencies are not installed, an error is returned.
func depsInstalled() (bool, error) {
	slog.Debug("checking for dependencies")

	_, err := getPackageManager()
	if err != nil {
		slog.Error("error checking for package manager", "err", err)
		return false, err
	}

	slog.Info("all dependencies found")

	return true, nil
}

// cacheExists returns true if the dnf package database cache exists.
func cacheExists() (bool, error) {
	slog.Debug("checking if cache exists")
	_, err := os.Stat("/var/cache/dnf/packages.db")

	if err != nil {
		if os.IsNotExist(err) {
			slog.Info("cache does not exist")
			return false, nil

		} else {
			slog.Error("error checking if cache exists", "err", err)
			return false, err
		}
	}

	slog.Info("cache found")
	return true, nil
}

// generateCache updates the dnf package database cache using the dnf update --refresh
// command.
func generateCache() {
	slog.Debug("generating cache")
	out, err := exec.Command("dnf", "update", "--refresh").Output()

	if err != nil {
		slog.Error("error generating cache", "err", err)
	}

	slog.Info("cache generated", "output", string(out))
}
