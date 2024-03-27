package selector

import (
	"context"
	"pkg/manager"
)

type Selector interface {
	SelectPackages(context.Context, []manager.Package) ([]manager.Package, error)
}
