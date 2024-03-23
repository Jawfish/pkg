package manager

import (
	"log/slog"
	"os/exec"

	"pkg/bin"
)

type Dnf struct {
	dnfBin bin.Executor
	escBin bin.Executor
	dnfCmd string
}

func NewDnf(escalationBinary bin.Executor, dnfBinary bin.Executor) *Dnf {
	return &Dnf{
		escBin: escalationBinary,
		dnfBin: dnfBinary,
		dnfCmd: string(dnfBinary.Name()),
	}
}

func (dnf *Dnf) GenerateCache() error {
	slog.Debug("generating cache")

	_, err := dnf.escBin.Execute(dnf.dnfCmd, "makecache")
	if err != nil {
		return &ErrGeneratingCache{Err: err}
	}

	slog.Debug("cache generated")
	return nil
}

func (dnf *Dnf) Install(pkg Package) error {
	slog.Debug("installing package", "pkg", pkg.Name)

	_, err := dnf.escBin.Execute(dnf.dnfCmd, "install", pkg.Name)
	if err != nil {
		return &ErrInstallingPackage{Pkg: pkg.Name, Err: err}
	}

	return nil
}

func (dnf *Dnf) Remove(pkg Package) (err error) {
	slog.Debug("uninstalling package", "pkg", pkg.Name)

	_, err = dnf.escBin.Execute(dnf.dnfCmd, "remove", pkg.Name)
	if err != nil {
		return &ErrRemovingPackage{Pkg: pkg.Name, Err: err}
	}

	return nil
}

func (dnf *Dnf) GetMetadata(pack Package) (Metadata, error) {
	slog.Debug("getting metadata for package", "package", pack.Name)

	out, err := exec.Command(dnf.dnfCmd, "info", pack.Name).Output()
	if err != nil {
		return Metadata{}, &ErrPkgMetadataNotFound{Pkg: pack, Err: err}
	}

	pkgMetadata := NewMetadata(pack.Name, "0.0.0", string(out))

	return pkgMetadata, nil
}
