package getter

import (
	"context"

	_ "github.com/mattn/go-sqlite3"

	"pkg/manager"
)

type Query string
type QueryType int

const (
	Installed QueryType = iota
	Available
	All
)

type PackageGetter interface {
	GetPackages(context.Context, Query, QueryType) ([]manager.Package, error)
}
