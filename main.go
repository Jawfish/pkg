package main

import (
	"flag"
	"log/slog"
	"os"
	"pkg/bin"
	"pkg/finder"
	"pkg/manager"
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

	binLocator := bin.OSBinaryLocator{}
	dnfBinary, err := binLocator.GetPreferredBinary(bin.DNF, bin.DNF5)
	if err != nil {
		slog.Error("error getting dnf binary", "err", err)
		os.Exit(1)
	}

	fzfBinary, err := binLocator.GetPreferredBinary(bin.FZF)
	if err != nil {
		slog.Error("error getting fzf binary", "err", err)
		os.Exit(1)
	}

	escalationBinary, err := binLocator.GetPreferredBinary(bin.Pkexec, bin.Doas, bin.Sudo)
	if err != nil {
		slog.Error("error getting escalation binary", "err", err)
		os.Exit(1)
	}

	pkgMgr := manager.NewDnf(escalationBinary, dnfBinary)
	pkgDb := initPackageDatabase(pkgMgr, cachePath)
	finder := finder.NewFinder(fzfBinary)

	packages, err := pkgDb.GetPackages(filter)
	if err != nil {
		slog.Error("error getting packages", "err", err)
		os.Exit(1)
	}

	// fmt.Println(packages)
	// fmt.Println(finder)

	// makeCacheService, err := bin.NewUnit("dnf-makecache.timer")
	// if err != nil {
	// 	slog.Error("error finding unit", "unit", "dnf-makecache.service", "err", err)
	// 	os.Exit(1)
	// }
	// makeCacheService = makeCacheService

	err = finder.Run(packages)
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
