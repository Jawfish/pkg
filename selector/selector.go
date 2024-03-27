package selector

import (
	"context"
	"pkg/manager"
)

type Finder interface {
	SelectPackages(context.Context, []manager.Package) ([]manager.Package, error)
}
