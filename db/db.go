package db

import (
	_ "github.com/mattn/go-sqlite3"

	"pkg/manager"
)

type PackageDatabase interface {
	GetPackages(filter string) ([]manager.Package, error)
}
