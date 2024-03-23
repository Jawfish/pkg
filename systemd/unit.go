package systemd

type LoadStatus string

const (
	Loaded   LoadStatus = "loaded"
	NotFound LoadStatus = "not-found"
)

type ActiveStatus string

const (
	Active   ActiveStatus = "active"
	Inactive ActiveStatus = "inactive"
)

type UnitName string

type Unit struct {
	Name   UnitName
	Loaded LoadStatus
	Active ActiveStatus
}

func NewUnit(name UnitName, unitLoaded LoadStatus, unitActive ActiveStatus) Unit {
	return Unit{
		Name:   name,
		Loaded: unitLoaded,
		Active: unitActive,
	}
}
