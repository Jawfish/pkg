package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"dnfzf/pkg"
)

type PackageDatabase struct {
	Path string
}

func NewPackageDatabase(path string) (*PackageDatabase, error) {
	pdb := &PackageDatabase{Path: path}
	err := pdb.validate()

	return pdb, err
}

func (pdb *PackageDatabase) validate() error {
	if err := pdb.validatePath(); err != nil {
		return err
	}

	if err := pdb.validateSchema(); err != nil {
		return err
	}

	return nil
}

// validatePath returns nil if the dnf package database cache exists and is a valid
// package database with the expected schema, otherwise an error is returned.
func (pdb *PackageDatabase) validatePath() error {
	slog.Debug("checking if cache exists")
	_, err := os.Stat(pdb.Path)

	if err != nil {
		if os.IsNotExist(err) {
			slog.Warn("cache does not exist")
			return &ErrCacheNotFound{Path: pdb.Path}

		} else {
			slog.Error("error checking if cache exists", "err", err)
			return err
		}
	}

	return nil
}

// validateSchema returns true if the dnf package database cache exists and is a valid
// package database with the expected schema.
func (pdb *PackageDatabase) validateSchema() error {
	slog.Debug("checking if the cache has a valid schema")

	db, err := sql.Open("sqlite3", pdb.Path)
	if err != nil {
		return &ErrOpeningCacheDatabase{Err: err}
	}
	defer db.Close()

	tables := []string{"installed", "available"}
	for _, table := range tables {
		query := "PRAGMA table_info(" + table + ")"
		rows, err := db.Query(query)
		if err != nil {
			return &ErrQueryFailed{Query: query, Err: err}
		}

		var cid int
		var name string
		var dtype string
		var notnull int
		var dflt_value sql.NullString
		var pk int

		for rows.Next() {
			err = rows.Scan(&cid, &name, &dtype, &notnull, &dflt_value, &pk)
			if err != nil {
				return &ErrScanFailed{Table: table, Err: err}
			}

			if cid != 0 || name != "pkg" || dtype != "TEXT" || notnull != 0 || pk != 0 {
				return &ErrInvalidCacheSchema{Table: table}
			}
		}
	}

	slog.Debug("cache has valid schema")
	return nil
}

// GetPackages returns a list of packages from the package database cache. The filter
// argument is used to filter the list of packages by name. Both installed and available
// packages are returned.
func (pdb *PackageDatabase) GetPackages(filter string) ([]pkg.Package, error) {
	rawAvailable, err := pdb.getRawPackages(pkg.Installed, filter)
	if err != nil {
		return nil, err
	}

	rawInstalled, err := pdb.getRawPackages(pkg.Installed, filter)
	if err != nil {
		return nil, err
	}

	available := pdb.processRawPackages(rawAvailable, pkg.Available)
	installed := pdb.processRawPackages(rawInstalled, pkg.Installed)

	return append(available, installed...), nil
}

// getRawPackages returns a list of packages from the dnf package database cache
// at the given path, either installed or available packages, depending on the pkgType
// argument.
func (pdb *PackageDatabase) getRawPackages(pkgType pkg.PackageType, filter string) ([]string, error) {
	db, err := sql.Open("sqlite3", pdb.Path)
	if err != nil {
		return nil, &ErrOpeningCacheDatabase{Err: err}
	}
	defer db.Close()

	// select all packages that are installed or available and match the filter
	query := fmt.Sprintf("SELECT pkg FROM %s WHERE pkg LIKE '%%%s%%'", pkgType, filter)

	rows, err := db.Query(query)
	if err != nil {
		return nil, &ErrQueryFailed{Query: query, Err: err}
	}
	defer rows.Close()

	var packages []string
	for rows.Next() {
		var pkg string
		if err := rows.Scan(&pkg); err != nil {
			return nil, &ErrScanFailed{Table: string(pkgType), Err: err}
		}
		packages = append(packages, pkg)
	}

	if err := rows.Err(); err != nil {
		return nil, &ErrQueryFailed{Query: query, Err: err}
	}

	return packages, nil
}

// processRawPackages processes the raw list of packages queried from the dnf package
// database and returns a list of Package structs.
func (pdb *PackageDatabase) processRawPackages(raw []string, pkgType pkg.PackageType) []pkg.Package {
	var packages []pkg.Package

	for _, pkgName := range raw {
		packages = append(packages, pkg.NewPackage(pkgName, pkgType))
	}
	return packages
}
