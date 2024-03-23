package main

import (
	"log/slog"
	"os"
	"pkg/db"
	"pkg/manager"
)

func initPackageDatabase(pkgMgr manager.PackageManager, cachePath string) *db.PackageDatabase {
	pkgDb, err := db.NewPackageDatabase(cachePath)

	// if the cache couldn't be validated, try to generate it
	if _, ok := err.(*db.ErrCacheNotFound); ok {
		slog.Warn("package database cache not found, generating...")
		err := pkgMgr.GenerateCache()
		if err != nil {
			slog.Error("error generating package database cache", "err", err)
			os.Exit(1)
		}

		// retry after generating the cache
		pkgDb, err = db.NewPackageDatabase(cachePath)
		if err != nil {
			slog.Error("error validating package database", "err", err)
			os.Exit(1)
		}
	}

	return pkgDb
}
