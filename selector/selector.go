package selector

import (
	"pkg/manager"
)

type Finder interface {
	SelectPackages([]manager.Package) ([]manager.Package, error)
}
