package executable

import (
	"fmt"
)

type ErrExecutableNotFound struct {
	Executable Executable
}

func (e *ErrExecutableNotFound) Error() string {
	return fmt.Sprintf("executable %s not found in $PATH", e.Executable)
}

type ErrNoValidExecutableFound struct {
	Binaries []Executable
}

func (e *ErrNoValidExecutableFound) Error() string {
	return fmt.Sprintf("none of the specified binaries (%s) could be found", e.Binaries)
}

type ErrCmdExec struct {
	Executable Executable
	Cmd        string
	Err        error
}

func (e *ErrCmdExec) Error() string {
	return fmt.Sprintf("error executing %s: %s", e.Executable, e.Err)
}
