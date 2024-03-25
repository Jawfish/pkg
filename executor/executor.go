package executor

import (
	"io"
	"log/slog"
	"os/exec"
	"strings"
)

type ExecutorIo struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type Executor interface {
	Execute(io ExecutorIo, args ...string) error
}

type Executable struct {
	Name BinaryName
}

func NewExecutable(name BinaryName) *Executable {
	return &Executable{Name: name}
}

// TODO: clean up how many arguments this function takes
func (e *Executable) Execute(io ExecutorIo, args ...string) error {
	slog.Debug("executing command", "binary", e.Name, "args", strings.Trim(strings.Join(args, " "), " "))

	bin := string(e.Name)
	cmd := exec.Command(bin, args...)
	cmd.Stdin = io.Stdin
	cmd.Stdout = io.Stdout
	cmd.Stderr = io.Stderr

	err := cmd.Run()
	if err != nil {
		// if the command exits with a non-zero status, leave it up to the process
		// that we executed to decide what to say or do about it
		if _, ok := err.(*exec.ExitError); ok {
			return nil
		}
		slog.Error("error executing command", "binary", e.Name, "cmd", strings.Join(args, " "), "err", err)
		return &ErrCmdExec{
			Binary: e.Name,
			Cmd:    strings.Join(args, " "),
			Err:    err,
		}
	}

	return nil
}
