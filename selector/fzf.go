package selector

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"pkg/executable"
	"pkg/manager"
	"strings"
)

type Fzf struct {
	executable string
	delimiter  string
	input      bytes.Buffer
}

func NewFzf(fzf executable.Executable) *Fzf {
	return &Fzf{
		executable: string(fzf),
		delimiter:  " ",
		input:      bytes.Buffer{},
	}
}

func (f *Fzf) SelectPackages(ctx context.Context, packages []manager.Package) ([]manager.Package, error) {
	slog.Debug("running finder", "executable", f.executable)

	pkgStr := f.prepareInput(packages)

	// TODO: the preview is pretty hacky and not very extendable once other package
	// managers are added
	args := []string{"--multi", "--with-nth", "1,2", "--delimiter", f.delimiter, "--tiebreak=length", "--ansi", "--preview", "dnf -C info {1} | tail -n +3"}

	cmd := exec.CommandContext(ctx, f.executable, args...)
	cmd.Stdin = strings.NewReader(pkgStr)
	cmd.Stdout = &f.input
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	slog.Debug("finder output", "output", f.input.String())

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 130 {
				slog.Debug("user sent SIGINT signal to deliberately exit fzf")
				return nil, &ErrProcessInterrupted{Code: 130}
			}
		}
		return nil, err
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	selectedPackages, err := f.parseOutput(packages)
	if err != nil {
		return nil, err
	}

	slog.Debug("selected packages", "packages", selectedPackages)
	return selectedPackages, nil
}

func (f *Fzf) prepareInput(packages []manager.Package) string {
	slog.Debug("preparing input to pass to fzf", "pkg_count", len(packages))

	var pkgNames []string
	for _, p := range packages {
		// append (installed) or (available) to package name
		// and also colorize that text based on its status
		name := fmt.Sprintf("%s\033[%dm (%s)\033[0m", p.Name, map[bool]int{true: 32, false: 31}[p.Installed], map[bool]string{true: "installed", false: "available"}[p.Installed])

		pkgNames = append(pkgNames, name)
	}
	return strings.Join(pkgNames, "\n")
}

func (f *Fzf) parseOutput(packages []manager.Package) ([]manager.Package, error) {
	slog.Debug("parsing fzf output", "output", f.input.String())
	selectedLines := strings.Split(strings.TrimSpace(f.input.String()), "\n")

	var selection []manager.Package
	var errs error

	pkgMap := make(map[string]manager.Package)
	for _, p := range packages {
		pkgMap[p.Name] = p
	}

	for _, line := range selectedLines {
		pkgName := strings.Split(line, " ")[0]
		if pkg, ok := pkgMap[pkgName]; ok {
			selection = append(selection, pkg)
		} else {
			errs = errors.Join(errs, fmt.Errorf("package %s not found in queried packages", pkgName))
		}
	}

	if errs != nil {
		return nil, errs
	}

	return selection, nil
}
