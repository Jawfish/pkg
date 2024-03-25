package manager

type PackageInstaller interface {
	Install(packages []Package) error
}

type PackageRemover interface {
	Remove(packages []Package) error
}

type PackageManager interface {
	PackageInstaller
	PackageRemover
}

type MetadataGetter interface {
	GetMetadata(Package) (Metadata, error)
}
