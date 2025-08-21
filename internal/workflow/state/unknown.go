package state

import (
	"fmt"
	"reflect"

	"github.com/growit-io/terrahoot/internal/platform"
)

// Unknown indicates that the discovery of the actual state of the workflow
// has not yet been attempted or that the attempt has failed.
//
// This is the most generic event that can trigger the infrastructure change
// management workflow. It fires when the workflow command-line program is
// called without a more specific subcommand.
//
// The handler for this event must inspect the execution environment of the
// current process to determine the "platform" it is running in, that is,
// whether it is executing locally or in a supported CI environment.
//
// The handler must then inspect the platform's local work tree of the
// infrastructure configuration Git repository and the process environment
// to determine the next action and resulting next workflow state.
type Unknown struct{ Environ []string }

func (s *Unknown) String() string {
	return reflect.TypeOf(*s).Name()
}

func (s *Unknown) Run() (State, error) {
	p, err := platform.FromEnvironment(s.Environ)
	if err != nil {
		return nil, fmt.Errorf("unable to determine whether the workflow is running in CI or locally: %w", err)
	}

	if p.CI() {
		return &Remote{p}, nil
	} else {
		return &Local{p}, nil
	}
}
