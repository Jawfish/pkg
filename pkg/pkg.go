package pkg

type PackageType string

const (
	Installed PackageType = "installed"
	Available PackageType = "available"
)

type Package struct {
	Name      string
	Installed bool
}

func NewPackage(pkg string, pkgType PackageType) Package {

	return Package{
		Name:      pkg,
		Installed: pkgType == Installed,
	}
}
