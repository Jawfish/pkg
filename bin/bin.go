package bin

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

type BinaryName string

type Executor interface {
	Execute(args ...string) (output []byte, err error)
	ExecuteWithStdin(stdin string, args ...string) (output []byte, err error)
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

func (b *Binary) Execute(arg ...string) (output []byte, err error) {
	slog.Debug("executing command", "binary", b.name, "args", strings.Trim(strings.Join(arg, " "), " "))
	return b.executeCommand(nil, arg...)
}

func (b *Binary) ExecuteWithStdin(stdin string, arg ...string) (output []byte, err error) {
	slog.Debug("executing command with stdin", "binary", b.name, "args", strings.Trim(strings.Join(arg, " "), " "))
	return b.executeCommand(strings.NewReader(stdin), arg...)
}
func (b *Binary) executeCommand(stdin io.Reader, arg ...string) (output []byte, err error) {
	bin := string(b.name)
	cmd := exec.Command(bin, arg...)
	cmd.Stdin = stdin

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

const (
	DNF5 BinaryName = "dnf5"
	DNF  BinaryName = "dnf"
)

const (
	FZF BinaryName = "fzf"
)

const (
	Pkexec BinaryName = "pkexec"
	Sudo   BinaryName = "sudo"
	Doas   BinaryName = "doas"
)
