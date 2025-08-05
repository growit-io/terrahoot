package state

// State represents the current workflow state and determines the possible
// transitions to other states.
type State interface {
	// Run implements a transition to the next workflow state. If the current
	// workflow state is final, then it is okay to return the same state or nil.
	Run() (State, error)

	// Returns a string that identifies the current workflow state to a human
	// user. It is not guaranteed that this string is unique among all possible
	// states, but it should be.
	String() string
}

// TODO: implement remaining states
type (
/*
// InDraft indicates that a draft proposal to merge the change has been
// created, but it may not be merged just yet.
InDraft struct{ NotImplemented }

// ReadyForReview indicates that the proposed change is no longer in a draft
// state. It can be reviewed and also possibly merged at any time, depending
// on the workflow configuration.
ReadyForReview struct{ NotImplemented }

// Reviewed indicates that the proposed change is no longer in a draft state
// and has received at least one review, but the current reviews do not yet
// indicate sufficient approval, according to the workflow configuration.
Reviewed struct{ NotImplemented }

// Approved indicates that the proposed change has received enough reviews
// to meet the approval criteria specified in the workflow configuration.
Approved struct{ NotImplemented }

// Merged indicates that the approved change was just merged into the base
// branch and should be applied from there immediately.
Merged struct{ NotImplemented }

// Applied indicates that the merged change was applied successfully and is
// ready to be released. The release workflow may be a separate workflow.
Applied struct{ NotImplemented }
*/
)
