package getter

import "fmt"

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
