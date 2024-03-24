# pkg

## What is it?

`pkg` is a CLI tool that improves the user experience when performing the most common package operations by providing an [fzf](https://github.com/junegunn/fzf)-backed interface for discovering, installing, and removing packages. It currently supports `dnf` and `dnf5`.

## What is it not?

`pkg` is not a full package manager. Advanced package management operations are outside of the scope of this project.

## Features

- quickly discover packages, even if you don't know the exact name
- mixed-state multi-select to install and remove multiple packages at the same time
- easily install specific versions of a package

## How does it work?

`pkg` queries the local package cache database then passes the list of packages matching the query to `fzf` where the user can search for and select any number of packages. Selected packages are installed or removed based on their current state.

When multiple packages are selected and those packages have mixed states (e.g. some are installed and some are not), `pkg` will always remove the installed packages before installing the available ones.

## Options

| Flag | Description                    | Default                            |
| ---- | -------------------------------| ---------------------------------- |
| `-c` | the path to the package cache  | `/var/cache/dnf/packages.db`       |
| `-v` | toggle verbose output          | `false`                            |
| `-y` | skip the confirmation prompt   | `false`                            |

## Dependencies

`pkg` depends on `fzf` for fuzzy searching and interactive selection. It also depends on a supported package manager being available. Currently, `pkg` supports `dnf` and `dnf5`.

If [`dnf5`](https://github.com/rpm-software-management/dnf5) is available, `pkg` will prefer that over `dnf` when installing or removing packages.

## Roadmap

- [ ] [improve responsiveness](https://github.com/Jawfish/pkg/issues/1)
- [ ] support for zypper
- [ ] support for other package managers (maybe)
