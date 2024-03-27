package executable

type Executable string

// Executables for package managers
const (
	Dnf5 Executable = "dnf5"
	Dnf  Executable = "dnf"
)

// Executables for fuzzy finders
const (
	Fzf Executable = "fzf"
)

// Executables for privilege escalation
const (
	Sudo   Executable = "sudo"
	Doas   Executable = "doas"
	Pkexec Executable = "pkexec"
)
