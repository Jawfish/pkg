package selector

import "fmt"

type ErrProcessInterrupted struct {
	Code int
}

func (e *ErrProcessInterrupted) Error() string {
	return fmt.Sprintf("process exited with code: %d", e.Code)
}
