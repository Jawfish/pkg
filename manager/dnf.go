package manager

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"pkg/executable"
)

type Dnf struct {
	executable     string
	rootExecutable string
	nonInteractive bool
}

func NewDnf(rootExecutable executable.Executable, executable executable.Executable, nonInteractive bool) *Dnf {
	return &Dnf{
		executable:     string(executable),
		rootExecutable: string(rootExecutable),
		nonInteractive: nonInteractive,
	}
}

func (dnf *Dnf) Install(ctx context.Context, packages []Package) error {
	slog.Debug("installing packages", "packages", packages)

	if len(packages) == 0 {
		return nil
	}

	var pkgNames []string
	for _, pkg := range packages {
		pkgNames = append(pkgNames, pkg.Name)
	}

	args := []string{dnf.executable, "install"}
	args = append(args, pkgNames...)
	if dnf.nonInteractive {
		args = append(args, "-y")
	}

	cmd := exec.CommandContext(ctx, dnf.rootExecutable, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	return err
}

func (dnf *Dnf) Remove(ctx context.Context, packages []Package) error {
	slog.Debug("removing packages", "packages", packages)

	if len(packages) == 0 {
		return nil
	}

	var pkgNames []string
	for _, pkg := range packages {
		pkgNames = append(pkgNames, dnf.getCleanName(pkg))
	}

	args := []string{dnf.executable, "remove"}
	args = append(args, pkgNames...)
	if dnf.nonInteractive {
		args = append(args, "-y")
	}

	cmd := exec.CommandContext(ctx, dnf.rootExecutable, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	return err
}

func (dnf *Dnf) GetMetadata(ctx context.Context, pack Package) (Metadata, error) {
	slog.Debug("getting metadata for package", "package", pack.Name)

	out, err := exec.CommandContext(ctx, "info", pack.Name).Output()
	if err != nil {
		return Metadata{}, &ErrPkgMetadataNotFound{Pkg: pack, Err: err}
	}

	pkgMetadata := NewMetadata(pack.Name, "0.0.0", string(out))

	return pkgMetadata, nil
}

func (dnf *Dnf) getCleanName(pkg Package) string {
	cleanName := strings.TrimSpace(strings.Fields(pkg.Name)[0])
	slog.Debug("cleaning package name", "original", pkg.Name, "cleaned", cleanName)
	return cleanName
}
