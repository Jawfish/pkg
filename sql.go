package main

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

type PackageType string

const (
	Installed PackageType = "installed"
	Available PackageType = "available"
)

// checkCache returns true if the dnf package database cache exists and is a valid
// package database with the expected schema.
func ensureValidSchema(path string) error {
	slog.Debug("checking if the cache has a valid schema")

	db, err := sql.Open("sqlite3", path)
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

// getPackagesFromCache returns a list of packages from the dnf package database cache
// at the given path, either installed or available packages, depending on the pkgType
// argument.
func getPackagesFromCache(path string, pkgType PackageType, sysArch string, filter string) ([]string, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, &ErrOpeningCacheDatabase{Err: err}
	}
	defer db.Close()

	// select all packages that are installed or available for the current architecture
	// and "noarch" packages (documentation, etc.)
	query := fmt.Sprintf("SELECT pkg FROM %s WHERE (pkg LIKE '%%%s%%' OR pkg LIKE '%%.noarch') AND pkg LIKE '%%%s%%'", pkgType, sysArch, filter)

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
