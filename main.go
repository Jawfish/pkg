package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"pkg/getter"
	"pkg/manager"
	"pkg/runner"
	"pkg/selector"
)

type Flags struct {
	noConfirm bool
	cachePath string
	verbose   bool
}

type Config struct {
	filter  getter.Query
	runners runner.Runners
	flags   Flags
}

func initialize() (Config, error) {
	noConfirm, cachePath, verbose := setFlags()
	flags := Flags{noConfirm, cachePath, verbose}

	args := flag.Args()
	filter := getter.Query(strings.Join(args, " "))

	initLogger(verbose)

	runners, err := runner.GetRunners()
	if err != nil {
		return Config{}, fmt.Errorf("error getting executable(s): %w", err)
	}

	return Config{filter, runners, flags}, nil
}

func main() {
	ctx := context.Background()

	i, err := initialize()
	if err != nil {
		slog.Error("Initialization error: ", "err", err)
		os.Exit(1)
	}

	mgr := manager.NewDnf(i.runners.Escalator, i.runners.Manager, i.flags.noConfirm)

	packageGetter, err := getter.NewDnfGetter(i.flags.cachePath)
	if err != nil {
		slog.Error("Error validating the package cache database: ", "err", err)
		os.Exit(1)
	}

	queriedPkgs, err := queryPkgs(ctx, packageGetter, i.filter)
	if err != nil {
		slog.Error("Error querying packages: ", "err", err)
		os.Exit(1)
	}

	pkgSelector := selector.NewFzf(i.runners.Selector)
	selectedPkgs, err := pkgSelector.SelectPackages(ctx, queriedPkgs)
	if err != nil {
		if _, ok := err.(*selector.ErrProcessInterrupted); ok {
			slog.Debug("Exiting gracefully")
			os.Exit(0)
		} else {
			slog.Error("Error running fzf", "err", err)
			os.Exit(1)
		}
	}

	err = managePackages(ctx, mgr, selectedPkgs)
	if err != nil {
		slog.Error("Package management error: ", "err", err)
		os.Exit(1)
	}

	slog.Debug("Everything seems to have gone well, exiting")
}

func setFlags() (noConfirm bool, cachePath string, verbose bool) {
	flag.BoolVar(&noConfirm, "y", false, "skip the confirmation prompt")
	flag.StringVar(&cachePath, "c", "/var/cache/dnf/packages.db", "the path to the DNF cache database")
	flag.BoolVar(&verbose, "v", false, "show verbose output")

	flag.Parse()

	return noConfirm, cachePath, verbose
}

func queryPkgs(ctx context.Context, packageGetter getter.PackageGetter, filter getter.Query) ([]manager.Package, error) {
	queryCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	queriedPkgs, err := packageGetter.GetPackages(queryCtx, filter, getter.All)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("getting packages timed out: %w", err)
		} else {
			return nil, fmt.Errorf("error getting packages: %w", err)
		}
	}
	return queriedPkgs, nil
}

func managePackages(ctx context.Context, mgr manager.PackageManager, selection []manager.Package) error {

	packagesToRemove := []manager.Package{}
	for _, pkg := range selection {
		if !pkg.Installed {
			continue
		}
		packagesToRemove = append(packagesToRemove, pkg)
	}

	packagesToInstall := []manager.Package{}
	for _, pkg := range selection {
		if pkg.Installed {
			continue
		}
		packagesToInstall = append(packagesToInstall, pkg)
	}

	err := mgr.Remove(ctx, packagesToRemove)
	if err != nil {
		return fmt.Errorf("error removing packages: %w", err)
	}

	err = mgr.Install(ctx, packagesToInstall)
	if err != nil {
		return fmt.Errorf("error installing packages: %w", err)
	}

	return nil
}
