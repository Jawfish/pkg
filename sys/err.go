package sys

import "fmt"

type ErrNoValidBinaryFound struct {
	Binaries []Binary
}

func (e *ErrNoValidBinaryFound) Error() string {
	return fmt.Sprintf("none of the specified binaries (%s) could be found", e.Binaries)
}
