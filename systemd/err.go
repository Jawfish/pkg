package systemd

import "fmt"

type ErrUnitNotActive struct {
	Unit Unit
}

func (e *ErrUnitNotActive) Error() string {
	return fmt.Sprintf("unit %s not active", e.Unit)
}

type ErrUnitNotFound struct {
	Name UnitName
}

func (e *ErrUnitNotFound) Error() string {
	return fmt.Sprintf("unit %s not found", e.Name)
}

type ErrNoUnitsFound struct{}

func (e *ErrNoUnitsFound) Error() string {
	return "no units found"
}

type ErrInvalidUnitLoadStatus struct {
	Status string
}

func (e *ErrInvalidUnitLoadStatus) Error() string {
	return fmt.Sprintf("invalid unit load status: %s", e.Status)
}

type ErrInvalidUnitActiveStatus struct {
	Status string
}

func (e *ErrInvalidUnitActiveStatus) Error() string {
	return fmt.Sprintf("invalid unit active status: %s", e.Status)
}

type ErrEnablingUnit struct {
	Unit Unit
}

func (e *ErrEnablingUnit) Error() string {
	return fmt.Sprintf("error enabling unit %s", e.Unit.Name)
}
