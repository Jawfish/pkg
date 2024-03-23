package systemd

import (
	"log/slog"
	"os/exec"
	"pkg/bin"
	"strings"
)

type Systemd struct {
	Binary bin.BinaryName
}

func NewSystemd(bin bin.BinaryName) Systemd {
	return Systemd{
		Binary: bin,
	}
}

type RawUnitList []string

func (s *Systemd) getUnitList() (RawUnitList, error) {
	slog.Debug("getting unit list")

	cmd := exec.Command(string(s.Binary), "list-units", "--all", "--no-legend", "--full")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, &ErrNoUnitsFound{}
	}

	if len(output) == 0 {
		return nil, &ErrNoUnitsFound{}
	}

	list := strings.Split(string(output), "\n")

	return list, nil
}

func (s *Systemd) EnableUnit(unit Unit) error {
	slog.Debug("enabling unit", "unit", unit.Name)

	cmd := exec.Command(string(s.Binary), "enable", string(unit.Name), "--now")
	err := cmd.Run()
	if err != nil {
		return &ErrEnablingUnit{Unit: unit}
	}

	return nil
}

func (s *Systemd) findUnit(name UnitName, units RawUnitList) (Unit, error) {
	slog.Debug("looking for unit", "name", name)

	for _, u := range units {
		fields := strings.Fields(u)
		if len(fields) < 3 {
			continue
		}

		if fields[0] == string(name) {

			loadStatus := LoadStatus(fields[1])
			if loadStatus != Loaded && loadStatus != NotFound {
				return Unit{}, &ErrInvalidUnitLoadStatus{fields[1]}
			}

			activeStatus := ActiveStatus(fields[2])
			if activeStatus != Active && activeStatus != Inactive {
				return Unit{}, &ErrInvalidUnitActiveStatus{fields[2]}
			}

			slog.Debug("unit found", "name", name, "loaded", loadStatus, "active", activeStatus)
			unit := NewUnit(name, loadStatus, activeStatus)

			return unit, nil
		}
	}

	return Unit{}, &ErrUnitNotFound{name}
}
