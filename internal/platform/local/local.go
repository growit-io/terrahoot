package local

import "github.com/growit-io/terrahoot/internal/git"

type Local struct{}

func New() *Local {
	return &Local{}
}

func (gh *Local) CI() bool {
	return false
}

func (gh *Local) GitDir() (*git.Dir, error) {
	return git.New().Dir()
}
