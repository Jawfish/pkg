package manager

type Status string

const (
	Installed Status = "installed"
	Available Status = "available"
)

type Package struct {
	Name      string
	Installed bool
}

func NewPackage(pkg string, status Status) Package {
	return Package{
		Name:      pkg,
		Installed: status == Installed,
	}
}

type Metadata struct {
	Name        string
	Version     string
	Description string
}

func NewMetadata(name, version, description string) Metadata {
	return Metadata{
		Name:        name,
		Version:     version,
		Description: description,
	}
}
