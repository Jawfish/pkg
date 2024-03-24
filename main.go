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
	dnfBinary, err := binLocator.GetPreferredBinary(bin.DNF5, bin.DNF)
	if err != nil {
		slog.Error("error getting dnf binary", "err", err)
		os.Exit(1)
	}

	fzfBinary, err := binLocator.GetPreferredBinary(bin.FZF)
	if err != nil {
		slog.Error("error getting fzf binary", "err", err)
		os.Exit(1)
	}

	rootBinary, err := binLocator.GetPreferredBinary(bin.Doas, bin.Sudo, bin.Pkexec)
	if err != nil {
		slog.Error("error getting escalation binary", "err", err)
		os.Exit(1)
	}

	pkgMgr := manager.NewDnf(rootBinary, dnfBinary)
	pkgDb := initPackageDatabase(pkgMgr, cachePath)
	finder := finder.NewFzf(fzfBinary)

	queriedPackages, err := pkgDb.GetPackages(filter)
	if err != nil {
		slog.Error("error getting packages", "err", err)
		os.Exit(1)
	}

	selectedPackages, err := finder.SelectPackages(queriedPackages)
	if err != nil {
		slog.Error("error running fzf", "err", err)
		os.Exit(1)
	}

	packagesToRemove := []manager.Package{}
	for _, pkg := range selectedPackages {
		if !pkg.Installed {
			continue
		}
		packagesToRemove = append(packagesToRemove, pkg)
	}

	packagesToInstall := []manager.Package{}
	for _, pkg := range selectedPackages {
		if pkg.Installed {
			continue
		}
		packagesToInstall = append(packagesToInstall, pkg)
	}

	err = pkgMgr.Remove(packagesToRemove)
	if err != nil {
		slog.Error("error removing packages", "err", err)
		os.Exit(1)
	}

	err = pkgMgr.Install(packagesToInstall)
	if err != nil {
		slog.Error("error installing packages", "err", err)
		os.Exit(1)
	}

	pkgMgr.GenerateCache()

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
