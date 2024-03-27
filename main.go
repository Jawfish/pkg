package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"os"
	"pkg/executable"
	"pkg/getter"
	"pkg/manager"
	"pkg/selector"
	"strings"
	"time"
)

func setFlags() (noConfirm bool, cachePath string, verbose bool) {
	flag.BoolVar(&noConfirm, "y", false, "skip the confirmation prompt")
	flag.StringVar(&cachePath, "c", "/var/cache/dnf/packages.db", "the path to the DNF cache database")
	flag.BoolVar(&verbose, "v", false, "show verbose output")

	flag.Parse()

	return noConfirm, cachePath, verbose
}

func managePackages(ctx context.Context, mgr manager.PackageManager, packagesToRemove, packagesToInstall []manager.Package) {
	err := mgr.Remove(ctx, packagesToRemove)
	if err != nil {
		slog.Error("error removing packages", "err", err)
		os.Exit(1)
	}

	err = mgr.Install(ctx, packagesToInstall)
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

	locator := executable.SystemExecLocator{}
	dnf, err := locator.GetExecutable(executable.Dnf)
	if err != nil {
		slog.Error("error getting dnf executable", "err", err)
		os.Exit(1)
	}

	fzfExecutable, err := locator.GetExecutable(executable.Fzf)
	if err != nil {
		slog.Error("error getting fzf executable", "err", err)
		os.Exit(1)
	}

	rootExecutable, err := locator.GetExecutable(executable.Doas, executable.Sudo, executable.Pkexec)
	if err != nil {
		slog.Error("error getting escalation executable", "err", err)
		os.Exit(1)
	}

	pkgGetter, err := getter.NewDnfGetter(cachePath)
	if err != nil {
		slog.Error("error getting dnf package cache database", "err", err)
		os.Exit(1)
	}

	ctx := context.Background()

	pkgManager := manager.NewDnf(rootExecutable, dnf, noConfirm)
	pkgManagerCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	queriedPkgs, err := pkgGetter.GetPackages(pkgManagerCtx, filter, getter.All)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("getting packages timed out", "err", err)
			os.Exit(1)
		} else {
			slog.Error("error getting packages", "err", err)
			os.Exit(1)
		}
	}

	pkgSelector := selector.NewFzf(fzfExecutable)
	selectedPkgs, err := pkgSelector.SelectPackages(ctx, queriedPkgs)
	if err != nil {
		if exitError, ok := err.(*selector.ErrProcessInterrupted); ok && exitError.Code == 130 {
			slog.Debug("Exiting gracefully")
			os.Exit(0)
		} else {
			slog.Error("Error running fzf", "err", err)
			os.Exit(1)
		}
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

	managePackages(ctx, pkgManager, packagesToRemove, packagesToInstall)

	slog.Debug("Everything seems to have gone well, exiting")
}
