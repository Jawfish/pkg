package manager

import (
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"pkg/bin"
)

type Dnf struct {
	bin    bin.Executor
	root   bin.Executor
	cmd    string
	stdOut io.Writer
	stdErr io.Writer
}

func NewDnf(root bin.Executor, bin bin.Executor) *Dnf {
	return &Dnf{
		bin:    bin,
		root:   root,
		cmd:    string(bin.Name()),
		stdOut: os.Stdout,
		stdErr: os.Stderr,
	}
}

func (dnf *Dnf) GenerateCache() error {
	slog.Debug("generating cache")

	err := dnf.root.Execute(nil, dnf.stdOut, dnf.stdErr, dnf.cmd, "makecache")
	if err != nil {
		return &ErrGeneratingCache{Err: err}
	}

	slog.Debug("cache generated")
	return nil
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

	err := dnf.root.Execute(os.Stdin, dnf.stdOut, dnf.stdErr, dnf.cmd, "install", strings.Join(pkgNames, " "))
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

	err := dnf.root.Execute(os.Stdin, dnf.stdOut, dnf.stdErr, dnf.cmd, "remove", strings.Join(pkgNames, " "))

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
