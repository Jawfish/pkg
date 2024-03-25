package executor

import (
	"fmt"
)

type ErrNoValidBinaryFound struct {
	Binaries []BinaryName
}

func (e *ErrNoValidBinaryFound) Error() string {
	return fmt.Sprintf("none of the specified binaries (%s) could be found", e.Binaries)
}

type ErrCmdExec struct {
	Binary BinaryName
	Cmd    string
	Err    error
}

func (e *ErrCmdExec) Error() string {
	return fmt.Sprintf("error executing %s: %s", e.Binary, e.Err)
}
