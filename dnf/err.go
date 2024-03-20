package dnf

import (
	"fmt"

	"dnfzf/pkg"
)

type ErrMissingDep struct {
	Dep string
}

func (e *ErrMissingDep) Error() string {
	return fmt.Sprintf("Dependency %s not found in PATH", e.Dep)
}

type ErrInstallingPackage struct {
	Pkg string
	Err error
}

func (e *ErrInstallingPackage) Error() string {
	return fmt.Sprintf("error installing package %s: %s", e.Pkg, e.Err)
}

type ErrUninstallingPackage struct {
	Pkg string
	Err error
}

func (e *ErrUninstallingPackage) Error() string {
	return fmt.Sprintf("error uninstalling package %s: %s", e.Pkg, e.Err)
}

type ErrGeneratingCache struct {
	Err error
}

func (e *ErrGeneratingCache) Error() string {
	return fmt.Sprintf("error generating cache: %s", e.Err)
}

type ErrPkgMetadataNotFound struct {
	Pkg pkg.Package
	Err error
}

func (e *ErrPkgMetadataNotFound) Error() string {
	return fmt.Sprintf("error getting metadata for package %s: %s", e.Pkg.Name, e.Err)
}
