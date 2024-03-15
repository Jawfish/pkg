package main

import (
	"fmt"
	"math"
	"strings"
)

type Package struct {
	Name      string
	Package   string
	Installed bool
}

// processPkgStr converts a raw entry queried from the dnf package database to a Package
// struct.
func processPkg(pkg string, pkgType PackageType) Package {
	splitPkg := strings.Split(pkg, ".")

	name := strings.Join(splitPkg[:len(splitPkg)-2], ".")

	return Package{
		Name:      name,
		Package:   pkg,
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

// getPackageMetadata returns the metadata for a given package.
func getPackageMetadata(pkg Package, width int) (string, error) {
	name := fmt.Sprintf("Name: %s", pkg.Name)
	packageName := fmt.Sprintf("Package: %s", pkg.Package)
	installed := fmt.Sprintf("Installed: %t", pkg.Installed)

	name = wrapText(name, int(math.Round(float64(width/2)-5)))
	packageName = wrapText(packageName, int(math.Round(float64(width/2)-5)))

	return fmt.Sprintf("%s\n%s\n%s\n", name, packageName, installed), nil
}

// wrapText wraps the input text at the specified width.
func wrapText(text string, width int) string {
	if len(text) <= width {
		return text
	}

	var result string
	for len(text) > width {
		result += text[:width] + "\n"
		text = text[width:]
	}
	result += text

	return result
}
