package manager

import (
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"pkg/executor"
)

type Dnf struct {
	dnfExecutable  *executor.Executable
	rootExecutable *executor.Executable
	nonInteractive bool
}

func NewDnf(rootExecutable *executor.Executable, executable *executor.Executable, nonInteractive bool) *Dnf {
	return &Dnf{
		dnfExecutable:  executable,
		rootExecutable: rootExecutable,
		nonInteractive: nonInteractive,
	}
}

func (dnf *Dnf) Install(packages []Package) error {
	slog.Debug("installing multiple packages", "packages", packages)

	if len(packages) == 0 {
		return nil
	}

	var pkgNames []string
	for _, pkg := range packages {
		pkgNames = append(pkgNames, pkg.Name)
	}

	args := []string{string(dnf.dnfExecutable.Name), "install"}
	args = append(args, pkgNames...)
	if dnf.nonInteractive {
		args = append(args, "-y")
	}

	execIo := executor.ExecutorIo{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	err := dnf.rootExecutable.Execute(execIo, args...)
	if err != nil {
		return err
	}

	return nil
}

func (dnf *Dnf) Remove(packages []Package) error {
	slog.Debug("removing packages", "packages", packages)

	if len(packages) == 0 {
		return nil
	}

	var pkgNames []string
	for _, pkg := range packages {
		pkgNames = append(pkgNames, dnf.getCleanName(pkg))
	}

	args := []string{string(dnf.dnfExecutable.Name), "remove"}
	args = append(args, pkgNames...)
	if dnf.nonInteractive {
		args = append(args, "-y")
	}

	execIo := executor.ExecutorIo{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	err := dnf.rootExecutable.Execute(execIo, args...)

	if err != nil {
		return err
	}

	return nil
}

func (dnf *Dnf) GetMetadata(pack Package) (Metadata, error) {
	slog.Debug("getting metadata for package", "package", pack.Name)

	out, err := exec.Command("info", pack.Name).Output()
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
