package manager

type PackageManager interface {
	GenerateCache() error
	Install(Package) error
	Remove(Package) error
	GetMetadata(Package) (Metadata, error)
}
