package runner

import (
	"errors"
	"log/slog"
	"os/exec"
)

type Runners struct {
	Manager   string
	Selector  string
	Escalator string
}

func getPreferredExecutable(names ...string) (string, error) {
	for _, bin := range names {
		slog.Debug("looking for executable", "executable", bin)
		path, err := exec.LookPath(bin)

		if err == nil && len(path) > 0 {
			slog.Debug("executable found", "path", path)
			return path, nil
		} else {
			slog.Debug("executable not found on $PATH", "executable", bin)
		}
	}

	return "", errors.New("no valid executables found")
}

func GetRunners() (Runners, error) {
	m, err := getPreferredExecutable("dnf")
	if err != nil {
		return Runners{}, err
	}

	s, err := getPreferredExecutable("fzf")
	if err != nil {
		return Runners{}, err
	}

	e, err := getPreferredExecutable("doas", "sudo", "pkexec")
	if err != nil {
		return Runners{}, err
	}

	return Runners{Manager: m, Selector: s, Escalator: e}, nil
}
