# pkg

## What is it?

`pkg` is a tool to improve the user experience for the most common package operations by providing a performant interface for discovering, installing, and removing packages on Fedora and other DNF-based systems.

## What is it not?

`pkg` is not a full package manager. Advanced package management operations are outside of the scope of the project; you should continue to use `dnf`, `rpm`, or `yum` for those.

## How does it work?

`pkg` queries DNF's local package cache database directly rather than using `dnf search` or similar. It's significantly faster than DNF at enumerating packages, even when compared against `dnf -C` to force `dnf` to use the local cache (see [benchmarks](#benchmarks)).

The list of packages matching the query populates the interface provided by [go-fuzzyfinder](https://github.com/ktr0731/go-fuzzyfinder) where the user can search for and select any number of packages. The packages are installed or removed based on their current state using `dnf` or `dnf5` (if available).

When multiple packages are selected and those packages have mixed states (e.g. some are installed and some are not), `pkg` will always remove the installed packages before installing the available ones. This is to account for situations where a package is being replaced by another package to avoid conflicts (similar to how `dnf swap` works).

Since `pkg` relies on the local package cache being available and reasonably fresh, it will:

- run `dnf makecache` if the cache is not found
- run `dnf makecache` if the cache is older than a week
- prompt the user to enable `dnf-makecache.timer` if it's available but not enabled

If the package manager encounters an error, `pkg` will display the error message and exit with a non-zero status code.

## Features

- quickly discover packages, even if you don't know the exact name
- mixed-state multi-select to install and remove multiple packages at the same time
- easily install specific versions of a package

## Options

`pkg` is an opinionated tool with sane defaults that should work for most users. However, there are a couple of options available as flags:

| Flag | Description                         | Default                            |
| ---- | ----------------------------------- | ---------------------------------- |
| `-y` | skip the confirmation prompt        | `false`                            |
| `-c` | the path to the DNF cache database  | `/var/cache/dnf/packages.db`       |
| `-v` | toggle verbose output               | `false`                            |

## Dependencies

Nothing should need to be installed on systems that `pkg` would be useful on other than the binary itself.

`pkg` uses the [go-fuzzyfinder](https://github.com/ktr0731/go-fuzzyfinder) library, meaning the `fzf` binary is not required.

It depends on `dnf`, of course, but that should already be present on any system where `pkg` would be useful.

If [`dnf5`](https://github.com/rpm-software-management/dnf5) is available, `pkg` will prefer that over `dnf` when installing or removing packages. Installing `dnf5` is recommended as it's significantly faster than `dnf` and doesn't conflict with it.

## Benchmarks

A lot of the speed improvement to be found from `pkg` comes not directly from the speed of the application, but from the improved ergonomics: fuzzy searching, being able to install and uninstall at the same time, quickly seeing available versions of a package, etc.

That said, there are also speed improvements to be found from the application itself.

The following benchmarks were conducted using the `time` command with output directed to `/dev/null`. Each operation was run five times, and the results were averaged. All `dnf` commands were run with the `-C` flag to force `dnf` to use the local cache.

<!-- TODO: these are outdated -->
| Operation                            | `dnf` | `pkg` | Improvement |
| ------------------------------------ | ----- | ------- | ----------- |
| Search for `kernel`                  | 0.98s | 0.011s  | 89x         |
| List all packages                    | 1.38s | 0.013s  | 106x        |

## Roadmap

- [ ] get package metadata directly from local cache via libsolv
