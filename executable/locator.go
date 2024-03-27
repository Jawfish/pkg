package executable

import (
	"log/slog"
	"os/exec"
)

type ExecLocator interface {
	GetExecutable(binaries []Executable) (Executable, error)
}

type SystemExecLocator struct{}

func (r *SystemExecLocator) findExecutable(name string) (Executable, error) {
	slog.Debug("looking for executable", "name", name)
	path, err := exec.LookPath(name)
	if len(path) > 0 {
		slog.Debug("executable found", "path", path)
	}
	return Executable(name), err
}

func (s *SystemExecLocator) GetExecutable(binaries ...Executable) (Executable, error) {
	for _, executable := range binaries {
		exe := string(executable)
		found, err := s.findExecutable(exe)

		if err == nil {
			return found, nil
		}
		slog.Debug("executable not found on $PATH", "executable", exe)
	}

	return "", &ErrNoValidExecutableFound{Binaries: binaries}
}
