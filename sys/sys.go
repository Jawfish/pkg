package sys

import (
	"log/slog"
	"os/exec"
)

type Binary string

// getPreferredBinary returns the first binary found in the PATH from the given list of
// binaries.
func getPreferredBinary(binaries []Binary) (bin Binary, err error) {
	for _, binary := range binaries {
		slog.Debug("looking for binary", "binary", binary)
		b := string(binary)
		found, err := exec.LookPath(b)

		if err == nil {
			slog.Debug("binary found", "binary", found)
			return binary, nil
		}

		slog.Debug("binary not found", "binary", b, "err", err)
	}

	return "", &ErrNoValidBinaryFound{Binaries: binaries}
}

type DnfBinary Binary

const (
	DNF5 DnfBinary = "dnf5"
	DNF  DnfBinary = "dnf"
)

// GetDnfBinary returns the package manager to use for the current system. dnf5 is
// preferred over dnf if both are available. If neither are available, an error is
// returned.
func GetDnfBinary() (bin DnfBinary, err error) {
	binaries := []Binary{"dnf5", "dnf"}
	b, err := getPreferredBinary(binaries)

	return DnfBinary(b), err
}

type FinderBinary Binary

const (
	FZF FinderBinary = "fzf"
)

// GetFinderBinary returns the fzf binary if it is found in the PATH. If it is not
// found, an error is returned.
func GetFinderBinary() (bin FinderBinary, err error) {
	binary := []Binary{"fzf"}
	b, err := getPreferredBinary(binary)

	return FinderBinary(b), err
}

type EscalationBinary Binary

const (
	Pkexec EscalationBinary = "pkexec"
	Sudo   EscalationBinary = "sudo"
	Doas   EscalationBinary = "doas"
)

// GetEscalationBinary returns the first available escalation binary in the PATH. The
// order of preference is pkexec > sudo > doas. If none are found, an error is returned.
func GetEscalationBinary() (bin EscalationBinary, err error) {
	binaries := []Binary{"pkexec", "sudo", "doas"}
	b, err := getPreferredBinary(binaries)

	return EscalationBinary(b), err
}

func GetBinaries() (dnf DnfBinary, fzf FinderBinary, escalation EscalationBinary, err error) {
	if dnf, err = GetDnfBinary(); err != nil {
		return
	}

	if fzf, err = GetFinderBinary(); err != nil {
		return
	}

	if escalation, err = GetEscalationBinary(); err != nil {
		return
	}

	slog.Debug("all binaries found")
	return
}
