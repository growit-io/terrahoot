package platform

import (
	"fmt"
	"strings"

	"github.com/growit-io/terrahoot/internal/git"
)

// Implementation describes a specific Platform implementation. Every Platform
// instance requires a concrete underlying implementation.
type Implementation interface {
	// CI indicates whether the implementation is for a CI/CD platform such as
	// GitHub Actions or GitLab CI/CD. The CI function should return false when
	// the implementation is for a local development environment.
	CI() bool

	// GitDir provides platform-specific access to a local Git work tree for the
	// infrastructure configuration repository.
	GitDir() (*git.Dir, error)
}

// Platform describes the project hosting platform or local development
// environment that underlies the workflow which is facilitated and largely
// implemented by the main application to which this package belogns.
type Platform struct {
	environ []string
	impl    Implementation
}

// Implementer describes an interface for constructing Platform implementations.
type Implementer interface {
	// Implement is allowed to return nil to indicate that the platform's
	// environment doesn't match the implementation.
	Implement(p *Platform) Implementation
}

var implementers []Implementer

// Register adds a new Platform implementation constructor to the list of known
// constructors for automatic discovery by the FromEnvironment function.
func Register(c Implementer) {
	implementers = append(implementers, c)
}

// FromEnvironment creates a new Platform instance based on the given set of
// environment variables in the form of a list of `key=value` pairs.
//
// Under normal conditions, the `environ` argument should be the return value of
// a call to `os.Environ()`. For testing purposes, it is also possible to pass a
// synthetic set of environment variables here in the same `key=value` form.
func FromEnvironment(environ []string) (*Platform, error) {
	p := &Platform{
		environ: environ,
		impl:    nil,
	}

	for _, c := range implementers {
		impl := c.Implement(p)
		if impl != nil {
			p.impl = impl
			return p, nil
		}
	}

	return nil, fmt.Errorf("no supported platform detected from environment")
}

// Getenv should behave similar to os.Getenv, but it operates on the environment
// variables of the Platform. It may be useful for Platform Implementations and
// their Implementers.
func (p *Platform) Getenv(key string) string {
	for _, kv := range p.environ {
		parts := strings.SplitN(kv, "=", 2)
		if parts[0] == key {
			return parts[1]
		}
	}
	return ""
}

// CI indicates whether the platform is a CI/CD platform such as GitHub Actions
// or GitLab CI/CD. It should return false when the program executes on a local
// machine during development.
func (p *Platform) CI() bool {
	return p.impl.CI()
}

// GitDir provides access to the platform's work tree of the infrastructure
// configuration Git repository.
func (p *Platform) GitDir() (*git.Dir, error) {
	return p.impl.GitDir()
}
