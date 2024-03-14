package main

type Package struct {
	Name      string
	Version   string
	Installed bool
}

// processPkgStr converts a raw entry queried from the dnf package database to a Package
// struct.
func processPkg(pkg string, pkgType PackageType) Package {
	name := pkg

	if pkgType == Installed {
		name += " (installed)"
	}

	return Package{
		Name:      name,
		Version:   "0.0.0",
		Installed: pkgType == Installed,
	}
}

// processPkgQuery processes the raw list of packages queried from the dnf package
// database and returns a list of Package structs.
func processPkgQuery(raw []string, pkgType PackageType) []Package {
	var packages []Package
	for _, pkg := range raw {
		packages = append(packages, processPkg(pkg, pkgType))
	}
	return packages
}
