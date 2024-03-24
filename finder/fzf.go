package finder

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"pkg/bin"
	"pkg/manager"
	"strings"
)

type Fzf struct {
	bin       bin.Executor
	delimiter string
	stdIn     bytes.Buffer
}

func NewFzf(fzfBinary bin.Executor) *Fzf {
	slog.Debug("creating new finder", "binary", fzfBinary.Name())

	return &Fzf{
		bin: fzfBinary,
		// non-whitespace characters show up for some reason, and space is too useful
		// to use as a delimiter, so use ‎ (invisible space character)
		delimiter: "‎",
		stdIn:     bytes.Buffer{},
	}
}

// wrapText wraps the input text at the specified width.
func (f *Fzf) wrapText(text string, width int) string {
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

func (f *Fzf) getPackagePreview(p *manager.Package, width int) string {
	fn := f.wrapText("Package: "+p.Name, int(float64(width/2)-5))

	installed := f.wrapText("Installed: "+fmt.Sprintf("%t", p.Installed), int(float64(width/2)-5))

	return fmt.Sprintf("%s\n%s\n", fn, installed)
}

func (f *Fzf) SelectPackages(packages []manager.Package) ([]manager.Package, error) {
	slog.Debug("running finder", "binary", f.bin.Name())

	pkgStr := f.prepareInput(packages)

	args := []string{"--multi", "--with-nth", "1", "--delimiter", f.delimiter, "--tiebreak=length", "--ansi"}

	err := f.bin.Execute(strings.NewReader(pkgStr), &f.stdIn, os.Stderr, args...)
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 130 {
				slog.Debug("user sent SIGINT signal to deliberately exit fzf")
				return nil, nil
			}
		}
		return nil, err
	}

	slog.Debug("finder output", "output", f.stdIn.String())

	selectedPackages, err := f.parseOutput(packages)
	if err != nil {
		return nil, err
	}

	slog.Debug("selected packages", "packages", selectedPackages)
	return selectedPackages, nil
}

func (f *Fzf) prepareInput(packages []manager.Package) string {
	slog.Debug("preparing input to give to finder", "pkg_count", len(packages))

	var pkgNames []string
	for _, p := range packages {
		name := p.Name

		if p.Installed {
			name = fmt.Sprintf("%s\033[32m (installed)\033[0m", name)
		} else {
			name = fmt.Sprintf("%s\033[34m (available)\033[0m", name)
		}

		pkgNames = append(pkgNames, name)
	}
	return strings.Join(pkgNames, "\n")
}

func (f *Fzf) parseOutput(packages []manager.Package) ([]manager.Package, error) {
	slog.Debug("parsing output from finder", "output", f.stdIn.String())
	selectedLines := strings.Split(strings.TrimSpace(f.stdIn.String()), "\n")

	var selection []manager.Package

	for _, line := range selectedLines {
		pkgName := strings.Split(line, " ")[0]
		slog.Debug("selected package", "package", pkgName)

		for _, p := range packages {
			if p.Name == pkgName {
				selection = append(selection, p)
				break
			}
		}
	}

	return selection, nil
}
