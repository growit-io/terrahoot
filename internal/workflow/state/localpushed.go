package state

import (
	"fmt"
	"reflect"
)

// LocalPushed indicates that the local change has been committed and pushed
// to a CI system in preparation for a review. A pull request may or may not
// have been opened.
type LocalPushed struct{ RemoteBranchRef string }

func (s *LocalPushed) String() string {
	return reflect.TypeOf(*s).Name()
}

func (s *LocalPushed) Run() (State, error) {
	return nil, fmt.Errorf("local Git work tree has no uncommitted changes, and is pushed to remote branch %s", s.RemoteBranchRef)
}
