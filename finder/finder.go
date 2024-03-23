package finder

import (
	"fmt"
	"log/slog"
	"pkg/bin"
	"pkg/manager"
	"strconv"
	"strings"
)

type Finder struct {
	fndBin    bin.Executor
	delimiter string
}

func NewFinder(fzfBinary bin.Executor) *Finder {
	slog.Debug("creating new finder", "binary", fzfBinary.Name())
	return &Finder{
		fndBin: fzfBinary,
		// non-whitespace characters show up for some reason, and space is too useful
		// to use as a delimiter, so use ‎ (invisible space character)
		delimiter: "‎",
	}
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

func (f *Finder) getPackagePreview(p *manager.Package, width int) string {
	fn := f.wrapText("Package: "+p.Name, int(float64(width/2)-5))

	installed := f.wrapText("Installed: "+fmt.Sprintf("%t", p.Installed), int(float64(width/2)-5))

	return fmt.Sprintf("%s\n%s\n", fn, installed)
}

func (f *Finder) Run(packages []manager.Package) error {
	slog.Debug("running finder", "binary", f.fndBin.Name())

	pkgStr := f.prepareInput(packages)

	args := []string{"--multi", "--with-nth", "1", "--delimiter", f.delimiter}

	out, err := f.fndBin.ExecuteWithStdin(pkgStr, args...)
	if err != nil {
		return err
	}

	selectedPackages, err := f.parseOutput(out)
	if err != nil {
		return err
	}

	slog.Debug("selected packages", "packages", selectedPackages)

	return nil
}

func (f *Finder) prepareInput(packages []manager.Package) string {
	slog.Debug("preparing input to give to finder", "pkg_count", len(packages))

	var pkgNames []string
	for _, p := range packages {
		name := p.Name

		if p.Installed {
			name += " (installed)"
		}

		pkgNames = append(pkgNames, fmt.Sprintf("%s%s%t", name, f.delimiter, p.Installed))
	}
	return strings.Join(pkgNames, "\n")
}

func (f *Finder) parseOutput(out []byte) ([]map[string]bool, error) {
	slog.Debug("parsing output from finder")
	selectedLines := strings.Split(strings.TrimSpace(string(out)), "\n")

	var selectedPackages []map[string]bool

	for _, line := range selectedLines {
		parts := strings.SplitN(line, f.delimiter, 2)
		if len(parts) < 2 {
			slog.Debug("too few parts in line", "line", line)
			continue
		}

		installed, err := strconv.ParseBool(parts[1])
		if err != nil {
			return nil, err
		}

		pkgMap := map[string]bool{parts[0]: installed}

		selectedPackages = append(selectedPackages, pkgMap)
	}

	return selectedPackages, nil
}
