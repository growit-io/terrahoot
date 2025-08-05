package state

import (
	"reflect"

	"github.com/growit-io/terrahoot/internal/platform"
)

// Remote indicates that a known and supported CI environment was detected.
// It is assumed that this environment normally has the ultimate authority
// over the actual infrastructure state.
type Remote struct{ *platform.Platform }

func (s *Remote) String() string {
	return reflect.TypeOf(*s).Name()
}

func (s *Remote) Run() (State, error) {
	return nil, nil
}
