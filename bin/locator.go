package bin

import (
	"log/slog"
	"os/exec"
)

type BinaryLocator interface {
	findBinary(file string) (BinaryName, error)
	GetPreferredBinary(binaries []BinaryName) (BinaryName, error)
}

type OSBinaryLocator struct{}

func (r *OSBinaryLocator) findBinary(name string) (string, error) {
	slog.Debug("looking for binary", "name", name)
	bin, err := exec.LookPath(name)
	slog.Debug("binary found", "path", bin)
	return bin, err
}

func (s *OSBinaryLocator) GetPreferredBinary(binaries ...BinaryName) (bin Executor, err error) {
	for _, binary := range binaries {
		b := string(binary)
		found, err := s.findBinary(b)

		if err == nil {
			return NewExecutable(BinaryName(found)), nil
		}
		slog.Debug("binary not found", "binary", b, "err", err)
	}

	return nil, &ErrNoValidBinaryFound{Binaries: binaries}
}
