package dnf

import (
	"log/slog"
	"os/exec"

	"dnfzf/pkg"
	"dnfzf/sys"
)

type PackageManager struct {
	dnfBinary        sys.DnfBinary
	escalationBinary sys.EscalationBinary
}

func NewPackageManager(escalationBinary sys.EscalationBinary, dnfBinary sys.DnfBinary) *PackageManager {
	return &PackageManager{
		escalationBinary: escalationBinary,
		dnfBinary:        dnfBinary,
	}
}

// GenerateCache creates the DNF package database cache.
func (pm *PackageManager) GenerateCache() error {
	slog.Debug("generating cache")

	esc := string(pm.escalationBinary)
	dnf := string(pm.dnfBinary)

	err := exec.Command(esc, dnf, "makecache").Run()

	if err != nil {
		return &ErrGeneratingCache{Err: err}
	}

	slog.Debug("cache generated")
	return nil
}

func (pm *PackageManager) InstallPackage(pkg pkg.Package) (err error) {
	slog.Debug("installing package", "pkg", pkg.Name)

	esc := string(pm.escalationBinary)
	dnf := string(pm.dnfBinary)

	err = exec.Command(esc, dnf, "install", pkg.Name).Run()
	if err != nil {
		return &ErrInstallingPackage{Pkg: pkg.Name, Err: err}
	}

	return nil
}

func (pm *PackageManager) UninstallPackage(pkg pkg.Package) (err error) {
	slog.Debug("uninstalling package", "pkg", pkg.Name)

	esc := string(pm.escalationBinary)
	dnf := string(pm.dnfBinary)

	err = exec.Command(esc, dnf, "remove", pkg.Name).Run()
	if err != nil {
		return &ErrUninstallingPackage{Pkg: pkg.Name, Err: err}
	}

	return nil
}

func (pm *PackageManager) GetPackageMetadata(pkg pkg.Package) (string, error) {
	slog.Debug("getting metadata for package", "pkg", pkg.Name)

	dnf := string(pm.dnfBinary)

	out, err := exec.Command(dnf, "info", pkg.Name).Output()
	if err != nil {
		return "", &ErrPkgMetadataNotFound{Pkg: pkg, Err: err}
	}

	return string(out), nil
}

// TODO: this is probably beyond the responsibility of the PackageManager. A method to
// create the cache is fine, but this is doing too much. Probably the caller should be
// the one checking the cache (via PackageDatabase) and then calling PackageManager if
// necessary.

// EnsureCache ensures that the dnf package database cache exists and is a valid package
// database with the expected schema. If the cache does not exist, it is generated.
// func (pm *PackageManager) EnsureCache(dnfBinary string, db PackageDatabase) {
// 	err := db.Validate()
// 	if err == nil {
// 		return
// 	}

// 	var e *ErrCacheNotFound
// 	if !errors.As(err, &e) {
// 		slog.Error("error checking cache", "err", err)
// 		os.Exit(1)
// 	}

// 	slog.Warn("cache not found, generating")
// 	if err = pm.generateCache(dnfBinary); err != nil {
// 		slog.Error("error generating cache", "err", err)
// 		os.Exit(1)
// 	}

// 	if err = db.Validate(); err != nil {
// 		if errors.As(err, &e) {
// 			slog.Error("cache still can't be found after generating. If you are using the -c flag, ensure the path matches the configuration for DNF.")
// 		} else {
// 			slog.Error("error checking cache", "err", err)
// 		}
// 		os.Exit(1)
// 	}
// }
