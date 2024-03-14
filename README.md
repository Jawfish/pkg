# dnfzf

## What is it?

`dnfzf` is a tool to improve the user experience for the most common package operations by providing a performant interface for discovering, installing, and removing packages on Fedora and other DNF-based systems.

## What is it not?

`dnfzf` is not a full package manager. Advanced package management operations are outside of the scope of the project; you should continue to use `dnf`, `rpm`, or `yum` for those.

## How does it work?

`dnfzf` queries DNF's local package cache database directly rather than using `dnf search` or similar. It's significantly faster than DNF at enumerating packages, even when compared against `dnf -C` to force `dnf` to use the local cache (see [benchmarks](#benchmarks)).

The list of packages matching the query populates the interface provided by [go-fuzzyfinder](https://github.com/ktr0731/go-fuzzyfinder) where the user can search for and select any number of packages. The packages are installed or removed based on their current state using `dnf` or `dnf5` (if available).

When multiple packages are selected and those packages have mixed states (e.g. some are installed and some are not), `dnfzf` will always remove the installed packages before installing the available ones. This is to account for situations where a package is being replaced by another package to avoid conflicts (similar to how `dnf swap` works).

Since `dnfzf` relies on the local package cache being available and reasonably fresh, it will:

- prompt the user to run `dnf makecache` if the cache is not found
- prompt the user to run `dnf makecache` if the cache is older than a week
- prompt the user to enable `dnf-makecache.timer` if it's available but not enabled

If the package manager encounters an error, `dnfzf` will display the error message and exit with a non-zero status code.

## Features

- quickly discover packages, even if you don't know the exact name
- mixed-state multi-select support for installing and removing multiple packages at once
- preview package metadata without having to run a separate command
- easily install specific versions of a package

## Options

`dnfzf` is an opinionated tool with sane defaults that should work for most users. However, there are a couple of options available as flags:

| Flag | Description                         | Default                            |
| ---- | ----------------------------------- | ---------------------------------- |
| `-y` | skip the confirmation prompt        | `false`                            |
| `-c` | the path to the DNF cache database  | `/var/cache/dnf/packages.db`       |

## Dependencies

Nothing should need to be installed on systems that `dnfzf` would be useful on other than the binary itself.

`dnfzf` uses the [go-fuzzyfinder](https://github.com/ktr0731/go-fuzzyfinder) library, meaning the `fzf` binary is not required.

It depends on `dnf`, of course, but that should already be present on any system where `dnfzf` would be useful.

If [`dnf5`](https://github.com/rpm-software-management/dnf5) is available, `dnfzf` will prefer that over `dnf` when installing or removing packages. Installing `dnf5` is recommended as it's significantly faster than `dnf` and doesn't conflict with it.

## Benchmarks

A lot of the speed improvement to be found from `dnfzf` comes not directly from the speed of the application, but from the improved ergonomics: fuzzy searching, being able to install and uninstall at the same time, quickly seeing available versions of a package, etc.

That said, there are also speed improvements to be found from the application itself.

The following benchmarks were conducted using the `time` command with output directed to `/dev/null`. Each operation was run five times, and the results were averaged. All `dnf` commands were run with the `-C` flag to force `dnf` to use the local cache.

| Operation                            | `dnf` | `dnfzf` | Improvement |
| ------------------------------------ | ----- | ------- | ----------- |
| Search for `kernel`                  | 0.98s | 0.011s  | 89x         |
| List all packages                    | 1.38s | 0.013s  | 106x        |
| Query `kernel-6.5.6-300.fc39.x86_64` | 0.74s | 0.002s  | 372x*       |

*`dnfzf` will find the package, but displaying the metadata takes roughly the same amount of time as `dnf` because it's a separate operation that uses `dnf` under the hood. However, you don't have to wait for the metadata query to complete before interacting with the package list.

## Roadmap

- [ ] handle conflicting packages
- [ ] get package metadata directly from local cache instead of using `dnf -C info`
