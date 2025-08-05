package state

import (
	"fmt"
	"reflect"
)

// LocalUncommitted indicates that the change is not yet entirely committed
// into a version control system.
type LocalUncommitted struct{}

func (s *LocalUncommitted) String() string {
	return reflect.TypeOf(*s).Name()
}

func (s *LocalUncommitted) Run() (State, error) {
	return nil, fmt.Errorf("local Git work tree has uncommitted changes")
}
