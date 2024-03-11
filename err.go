package main

import "fmt"

type ErrUnsupportedArch struct {
	Arch string
}

func (e *ErrUnsupportedArch) Error() string {
	return fmt.Sprintf("Architecture %s not supported", e.Arch)
}

type ErrMissingDep struct {
	Dep string
}

func (e *ErrMissingDep) Error() string {
	return fmt.Sprintf("Dependency %s not found in PATH", e.Dep)
}

type ErrPackageManagerNotFound struct {
}

func (e *ErrPackageManagerNotFound) Error() string {
	return "No package manager found in PATH. This program requires dnf or dnf5 to be installed and in the PATH."
}
