package main

import "fmt"

type ErrUnsupportedArch struct {
	Arch string
}

func (e *ErrUnsupportedArch) Error() string {
	return fmt.Sprintf("architecture %s not supported", e.Arch)
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
	return "no package manager found in PATH"
}

type ErrInvalidCacheSchema struct {
	Table string
}

func (e *ErrInvalidCacheSchema) Error() string {
	return fmt.Sprintf("invalid schema for table %s", e.Table)
}

type ErrQueryFailed struct {
	Query string
	Err   error
}

func (e *ErrQueryFailed) Error() string {
	return fmt.Sprintf("query %s failed: %s", e.Query, e.Err)
}

type ErrScanFailed struct {
	Table string
	Err   error
}

func (e *ErrScanFailed) Error() string {
	return fmt.Sprintf("scanning table %s failed: %s", e.Table, e.Err)
}

type ErrOpeningCacheDatabase struct {
	Err error
}

func (e *ErrOpeningCacheDatabase) Error() string {
	return fmt.Sprintf("error opening cache database: %s", e.Err)
}

type ErrCacheNotFound struct {
	Path string
}

func (e *ErrCacheNotFound) Error() string {
	return fmt.Sprintf("cache not found at path %s", e.Path)
}
