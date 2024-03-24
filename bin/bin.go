package bin

import (
	"io"
	"log/slog"
	"os/exec"
	"strings"
)

type BinaryName string

type Executor interface {
	Execute(stdin io.Reader, stdout, stderr io.Writer, arg ...string) error
	Name() BinaryName
}

type Binary struct {
	name BinaryName
}

func NewExecutable(bin BinaryName) *Binary {
	return &Binary{
		name: bin,
	}
}

func (b *Binary) Name() BinaryName {
	return b.name
}

func (b *Binary) Execute(stdin io.Reader, stdout, stderr io.Writer, arg ...string) error {
	slog.Debug("executing command", "binary", b.name, "args", strings.Trim(strings.Join(arg, " "), " "))

	bin := string(b.name)
	cmd := exec.Command(bin, arg...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err != nil {
		// if the command exits with a non-zero status, leave it up to the process
		// that we executed to decide what to say or do about it
		if _, ok := err.(*exec.ExitError); ok {
			return nil
		}
		slog.Error("error executing command", "binary", b.name, "cmd", strings.Join(arg, " "), "err", err)
		return &ErrCmdExec{
			Binary: b.name,
			Cmd:    strings.Join(arg, " "),
			Err:    err,
		}
	}

	return nil
}

const (
	DNF5 BinaryName = "dnf5"
	DNF  BinaryName = "dnf"
)

const (
	FZF BinaryName = "fzf"
)

const (
	Sudo   BinaryName = "sudo"
	Doas   BinaryName = "doas"
	Pkexec BinaryName = "pkexec"
)
