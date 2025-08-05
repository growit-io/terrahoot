package state

import (
	"fmt"
	"reflect"
)

// LocalCommitted indicates that the change is locally committed to version
// control entirely, but the current revision has not yet been submitted to
// a CI system for review.
type LocalCommitted struct {
	RemoteBranchRef string
}

func (s *LocalCommitted) String() string {
	return reflect.TypeOf(*s).Name()
}

func (s *LocalCommitted) Run() (State, error) {
	return nil, fmt.Errorf("local Git branch is out of sync with remote tracking branch \"%s\"", s.RemoteBranchRef)
}
