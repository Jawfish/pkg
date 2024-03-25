package executor

type BinaryName string

const (
	DNF5 BinaryName = "dnf5"
	DNF  BinaryName = "dnf"
)

const (
	FZF BinaryName = "fzf"
)

const (
	Sudo   BinaryName = "sudo"
	Doas   BinaryName = "doas"
	Pkexec BinaryName = "pkexec"
)
