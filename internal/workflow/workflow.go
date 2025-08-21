package workflow

import (
	"fmt"

	"github.com/growit-io/terrahoot/internal/workflow/state"
)

// Workflow represents the opinionated Terragrunt CI/CD workflow that is
// facilitated and largely implemented by the main application to which this
// package belongs.  Each instance of Workflow represents a partial execution
// of the workflow and the current workflow state is discovered entirely from
// the surrounding platform.
type Workflow struct {
	state state.State
	debug bool
}

// FromEnvironment creates a new Workflow instance from the given environment
// variables.
func FromEnvironment(environ []string) *Workflow {
	state := &state.Unknown{Environ: environ}
	debug := false

	return &Workflow{
		state,
		debug,
	}
}

// WithDebug enables or disables additional debug information during workflow
// execution.
func (w *Workflow) WithDebug(value bool) *Workflow {
	copy := *w
	copy.debug = value
	return &copy
}

// Run proceeds with the workflow execution until the next point at which the
// underlying platform needs to take control again and is possibly waiting for
// an interaction with the change author(s).
func (w *Workflow) Run() error {
	for {
		s, err := w.state.Run()
		if err != nil {
			return err
		}

		if s == nil || w.state == s {
			return nil
		}

		if w.debug {
			fmt.Printf("State: %s -> %s\n", w.state.String(), s.String())
		}

		w.state = s
	}
}
