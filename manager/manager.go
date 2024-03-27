package manager

import "context"

type PackageInstaller interface {
	Install(context.Context, []Package) error
}

type PackageRemover interface {
	Remove(context.Context, []Package) error
}

type PackageManager interface {
	PackageInstaller
	PackageRemover
}

type MetadataGetter interface {
	GetMetadata(context.Context, Package) (Metadata, error)
}
