package getter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"pkg/manager"
)

type DnfDatabase struct {
	Path string
}

func NewDnfGetter(path string) (*DnfDatabase, error) {
	pdb := &DnfDatabase{Path: path}
	err := pdb.validate()

	return pdb, err
}

func (pdb *DnfDatabase) validate() error {
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
func (pdb *DnfDatabase) validatePath() error {
	slog.Debug("looking for cache", "path", pdb.Path)
	_, err := os.Stat(pdb.Path)

	if err != nil {
		if os.IsNotExist(err) {
			return &ErrCacheNotFound{Path: pdb.Path}

		} else {
			return err
		}
	}

	slog.Debug("cache found", "path", pdb.Path)
	return nil
}

// validateSchema returns true if the dnf package database cache exists and is a valid
// package database with the expected schema.
func (pdb *DnfDatabase) validateSchema() error {
	slog.Debug("checking cache schema")

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

	slog.Debug("cache schema valid")
	return nil
}

// GetPackages returns a list of packages from the package database cache. The query
// argument is used to filter the list of packages by name. Both installed and available
// packages are returned.
func (pdb *DnfDatabase) GetPackages(ctx context.Context, query Query, queryType QueryType) ([]manager.Package, error) {
	// use a slice and a corresponding map of package names to their index in the
	// slice so that we can return a slice without having to flatten a map
	pkgMap := make(map[string]int)
	pkgs := make([]manager.Package, 0)

	packageTypes := []string{}
	if queryType == Available || queryType == All {
		packageTypes = append(packageTypes, "available")
	}
	if queryType == Installed || queryType == All {
		packageTypes = append(packageTypes, "installed")
	}

	for _, packageType := range packageTypes {
		rawPackages, err := pdb.getRawPackages(packageType, query)
		if err != nil {
			return nil, fmt.Errorf("error getting %s packages: %w", packageType, err)
		}

		for _, pkg := range rawPackages {
			// if the package is installed and already in the slice,
			// update the installed field
			isInstalled := packageType == "installed"
			if i, ok := pkgMap[pkg]; ok {
				pkgs[i].Installed = isInstalled
			} else {
				// otherwise, we haven't seen this package yet, so add it to the slice
				p := manager.NewPackage(pkg, isInstalled)
				pkgMap[p.Name] = len(pkgs)
				pkgs = append(pkgs, p)
			}
		}
	}

	return pkgs, nil
}

// getRawPackages returns a list of packages from the dnf package database cache for
// the given table (installed or available).
func (pdb *DnfDatabase) getRawPackages(table string, filter Query) ([]string, error) {

	db, err := sql.Open("sqlite3", pdb.Path)
	if err != nil {
		return nil, &ErrOpeningCacheDatabase{Err: err}
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT pkg FROM %s WHERE pkg LIKE '%%%s%%'", table, filter)

	rows, err := db.Query(query)
	if err != nil {
		return nil, &ErrQueryFailed{Query: query, Err: err}
	}
	defer rows.Close()

	var packages []string
	for rows.Next() {
		var pkg string
		if err := rows.Scan(&pkg); err != nil {
			return nil, &ErrScanFailed{Table: string(table), Err: err}
		}
		packages = append(packages, pkg)
	}

	if err := rows.Err(); err != nil {
		return nil, &ErrQueryFailed{Query: query, Err: err}
	}

	return packages, nil
}
