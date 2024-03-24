package manager

type Package struct {
	Name      string
	Installed bool
}

func NewPackage(pkg string, installed bool) Package {
	return Package{
		Name:      pkg,
		Installed: installed,
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
