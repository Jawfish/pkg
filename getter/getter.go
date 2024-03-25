package getter

import (
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
	GetPackages(Query, QueryType) ([]manager.Package, error)
}
