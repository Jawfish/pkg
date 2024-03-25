package main

import (
	"flag"
	"log/slog"
	"os"
	"pkg/executor"
	"pkg/getter"
	"pkg/manager"
	"pkg/selector"
	"strings"
)

func setFlags() (noConfirm bool, cachePath string, verbose bool) {
	flag.BoolVar(&noConfirm, "y", false, "skip the confirmation prompt")
	flag.StringVar(&cachePath, "c", "/var/cache/dnf/packages.db", "the path to the DNF cache database")
	flag.BoolVar(&verbose, "v", false, "show verbose output")

	flag.Parse()

	return noConfirm, cachePath, verbose
}

func managePackages(mgr manager.PackageManager, packagesToRemove, packagesToInstall []manager.Package) {
	err := mgr.Remove(packagesToRemove)
	if err != nil {
		slog.Error("error removing packages", "err", err)
		os.Exit(1)
	}

	err = mgr.Install(packagesToInstall)
	if err != nil {
		slog.Error("error installing packages", "err", err)
		os.Exit(1)
	}
}

func main() {
	noConfirm, cachePath, verbose := setFlags()

	filters := flag.Args()
	filter := getter.Query(strings.Join(filters, " "))

	initLogger(verbose)

	binLocator := executor.OSBinaryLocator{}
	dnfBinary, err := binLocator.GetExecutable(executor.DNF)
	if err != nil {
		slog.Error("error getting dnf binary", "err", err)
		os.Exit(1)
	}

	fzfBinary, err := binLocator.GetExecutable(executor.FZF)
	if err != nil {
		slog.Error("error getting fzf binary", "err", err)
		os.Exit(1)
	}

	rootBinary, err := binLocator.GetExecutable(executor.Doas, executor.Sudo, executor.Pkexec)
	if err != nil {
		slog.Error("error getting escalation binary", "err", err)
		os.Exit(1)
	}

	pkgGetter, err := getter.NewDnfGetter(cachePath)
	if err != nil {
		slog.Error("error getting dnf package cache database", "err", err)
		os.Exit(1)
	}

	pkgManager := manager.NewDnf(rootBinary, dnfBinary, noConfirm)

	queriedPkgs, err := pkgGetter.GetPackages(filter, getter.All)
	if err != nil {
		slog.Error("error getting packages", "err", err)
		os.Exit(1)
	}

	pkgSelector := selector.NewFzf(fzfBinary)
	selectedPkgs, err := pkgSelector.SelectPackages(queriedPkgs)
	if err != nil {
		slog.Error("error running fzf", "err", err)
		os.Exit(1)
	}

	packagesToRemove := []manager.Package{}
	for _, pkg := range selectedPkgs {
		if !pkg.Installed {
			continue
		}
		packagesToRemove = append(packagesToRemove, pkg)
	}

	packagesToInstall := []manager.Package{}
	for _, pkg := range selectedPkgs {
		if pkg.Installed {
			continue
		}
		packagesToInstall = append(packagesToInstall, pkg)
	}

	managePackages(pkgManager, packagesToRemove, packagesToInstall)

}
