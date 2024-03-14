package main

// import (
// 	"database/sql"
// )

// const availablePackagesQuery = "SELECT * FROM available"
// const installedPackagesQuery = "SELECT * FROM installed"

// func nothing() {
// 	db, err := sql.Open("sqlite3", "file:packages.db?cache=shared")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()
// }
import (
	"database/sql"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
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
