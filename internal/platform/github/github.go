package github

import (
	"github.com/growit-io/terrahoot/internal/git"
	"github.com/growit-io/terrahoot/internal/platform/local"
)

type GitHub struct {
	local *local.Local
}

func New() *GitHub {
	return &GitHub{local.New()}
}

func (gh *GitHub) CI() bool {
	return true
}

func (gh *GitHub) GitDir() (*git.Dir, error) {
	return gh.local.GitDir()
}
