package finder

import (
	"pkg/bin"
	"pkg/manager"
)

type Finder interface {
	Run([]manager.Package) error
}

type Fzf struct {
	fndBin    bin.Executor
	delimiter string
}
