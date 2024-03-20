package main

import (
	"dnfzf/db"
	"dnfzf/dnf"
	"dnfzf/finder"
	"dnfzf/sys"
	"flag"
	"log/slog"
	"os"
	"strings"
)

func setFlags(noConfirm *bool, cachePath *string, verbose *bool) {
	flag.BoolVar(noConfirm, "y", false, "skip the confirmation prompt")
	flag.StringVar(cachePath, "c", "/var/cache/dnf/packages.db", "the path to the DNF cache database")
	flag.BoolVar(verbose, "v", false, "show verbose output")

	flag.Parse()
}

func main() {
	var (
		noConfirm bool
		cachePath string
		verbose   bool
	)
	setFlags(&noConfirm, &cachePath, &verbose)

	filters := flag.Args()
	filter := strings.Join(filters, " ")

	initLogger(verbose)

	dnfBinary, fzfBinary, escalationBinary, err := sys.GetBinaries()
	if err != nil {
		slog.Error("error getting binary", "err", err)
		os.Exit(1)
	}

	pkgMgr := dnf.NewPackageManager(escalationBinary, dnfBinary)
	pkgDb, err := db.NewPackageDatabase(cachePath)
	finder := finder.NewFinder(fzfBinary)

	if err != nil {
		// if the cache couldn't be validated, try to generate it
		if _, ok := err.(*db.ErrCacheNotFound); ok {
			err := pkgMgr.GenerateCache()
			if err != nil {
				slog.Error("error generating package database cache", "err", err)
				os.Exit(1)
			}

			// retry after generating the cache
			pkgDb, err = db.NewPackageDatabase(cachePath)
		}

		if err != nil {
			slog.Error("error validating package database", "err", err)
			os.Exit(1)
		}
	}

	packages, err := pkgDb.GetPackages(filter)
	if err != nil {
		slog.Error("error getting packages", "err", err)
		os.Exit(1)
	}

	// fmt.Println(packages)
	// fmt.Println(finder)

	err = finder.RunFinder(packages)
	if err != nil {
		slog.Error("error running fzf", "err", err)
		os.Exit(1)
	}

	// type PackageError struct {
	// 	Package Package
	// 	Err     error
	// }

	// var errors []PackageError

	// idx, err := fuzzyfinder.FindMulti(
	// 	processedPackages,
	// 	func(i int) string {
	// 		name := processedPackages[i].Name

	// 		if processedPackages[i].Installed {
	// 			name += " (installed)"
	// 		}
	// 		return name
	// 	},
	// 	fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
	// 		// if there aren't any packages, don't try to show a preview
	// 		if i == -1 {
	// 			return ""
	// 		}

	// 		md, err := getPackageMetadata(processedPackages[i], w)
	// 		if err != nil {
	// 			errors = append(errors, PackageError{Package: processedPackages[i], Err: err})
	// 			return ""
	// 		}

	// 		return md
	// 	}))

	// for _, pkgErr := range errors {
	// 	slog.Error("error getting package metadata", "package", pkgErr.Package.Name, "err", pkgErr.Err)
	// }

	// if err != nil {
	// 	slog.Error("error finding package", "err", err)
	// 	os.Exit(1)
	// }
	// fmt.Printf("selected: %v\n", idx)
}
