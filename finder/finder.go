package finder

import (
	"bytes"
	"dnfzf/pkg"
	"dnfzf/sys"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Finder struct {
	fzfBinary sys.FinderBinary
}

func NewFinder(fzfBinary sys.FinderBinary) *Finder {
	slog.Debug("creating new finder", "binary", fzfBinary)
	return &Finder{fzfBinary: fzfBinary}
}

// wrapText wraps the input text at the specified width.
func (f *Finder) wrapText(text string, width int) string {
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

func (f *Finder) getPackagePreview(p *pkg.Package, width int) string {
	fn := f.wrapText("Package: "+p.Name, int(float64(width/2)-5))

	installed := f.wrapText("Installed: "+fmt.Sprintf("%t", p.Installed), int(float64(width/2)-5))

	return fmt.Sprintf("%s\n%s\n", fn, installed)
}

// TODO: split this up
func (f *Finder) RunFinder(packages []pkg.Package) error {
	slog.Debug("running fzf", "binary", f.fzfBinary)
	// execute fzf with the list of packages using pkg.Name as the display name
	b := string(f.fzfBinary)

	// convert the packages into a newline-separated string for fzf
	var pkgNames []string
	for _, p := range packages {
		name := p.Name

		if p.Installed {
			name += "(installed)"
		}

		pkgNames = append(pkgNames, fmt.Sprintf("%s %s %t", p.Name, name, p.Installed))
	}
	pkgStr := strings.Join(pkgNames, "\n")

	cmd := exec.Command(b, "--multi", "--with-nth", "1", "--delimiter", " ")

	cmd.Stdin = strings.NewReader(pkgStr)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	selectedLines := strings.Split(strings.TrimSpace(out.String()), "\n")

	var selectedPackages []map[string]bool

	for _, line := range selectedLines {
		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}

		installed, err := strconv.ParseBool(parts[1])
		if err != nil {
			return err
		}

		pkgMap := map[string]bool{parts[0]: installed}

		selectedPackages = append(selectedPackages, pkgMap)
	}

	slog.Debug("selected packages", "packages", selectedPackages)

	return nil
}
