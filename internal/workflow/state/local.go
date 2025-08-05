package state

import (
	"fmt"
	"reflect"

	"github.com/growit-io/terrahoot/internal/platform"
)

// Local indicates that the workflow is running outside of any known and
// supported CI environment. It is assumed that it is being executed from
// the local workstation of a human administrator.
type Local struct{ *platform.Platform }

func (s *Local) String() string {
	return reflect.TypeOf(*s).Name()
}

func (s *Local) Run() (State, error) {
	gd, err := s.Platform.GitDir()
	if err != nil {
		return nil, fmt.Errorf("failed to find Git work tree: %w", err)
	}

	isClean, err := gd.IsClean()
	if err != nil {
		return nil, fmt.Errorf("failed to determine whether Git work tree is clean: %w", err)
	}
	isDirty := !isClean

	remoteBranchRef, err := gd.RemoteBranchRef()
	if err != nil {
		return nil, fmt.Errorf("failed to determine the remote Git branch reference: %w", err)
	}

	isUpToDate, err := gd.IsUpToDate()
	if err != nil {
		return nil, fmt.Errorf("failed to determine whether remote Git branch is up-to-date: %w", err)
	}

	switch {
	case isDirty:
		return &LocalUncommitted{}, nil
	case isUpToDate:
		return &LocalPushed{
			RemoteBranchRef: remoteBranchRef,
		}, nil
	default:
		return &LocalCommitted{
			RemoteBranchRef: remoteBranchRef,
		}, nil
	}
}
