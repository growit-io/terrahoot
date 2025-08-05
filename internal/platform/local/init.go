package local

import (
	"github.com/growit-io/terrahoot/internal/platform"
)

type Implementer struct{}

func (i *Implementer) Implement(p *platform.Platform) platform.Implementation {
	if p.Getenv("CI") != "" {
		return nil
	}

	return New()
}

func init() {
	platform.Register(&Implementer{})
}
