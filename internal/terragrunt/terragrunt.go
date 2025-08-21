package terragrunt

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/growit-io/terrahoot/internal/ui"
)

type Options struct {
	// DeletedUnits is similar to UpdatedFiles, but it lists only the parent
	// directories of deleted configuration unit files. A "destroy" action may
	// be performed for units that appear in this list before performing the
	// "apply" action for any other units that may be affected by the change.
	DeletedUnits []string

	// UpdatedFiles limits the "apply" action to only those configuration units
	// marked as reading one or more of the given files.  When UpdatedFiles is
	// empty, "run --all" does not limit units this way.
	//
	// Each path in UpdatedFiles must be either an absolute path or relative to
	// the current working directory. In other words, for any given file's path,
	// os.Stat() is supposed to return no error.
	UpdatedFiles []string
}

type Phase int

const (
	Plan Phase = iota
	Apply
)

func (c Phase) String() string {
	switch c {
	case Plan:
		return "plan"
	case Apply:
		return "apply"
	}
	panic(fmt.Sprintf("invalid %T value: %d", c, c))
}

// ParsePhase attempts to translate an arbitrary string into a known and
// supported Terragrunt command. It raises an error if the given named command
// isn't supported by the Run() function.
func ParsePhase(name string) (Phase, error) {
	switch name {
	case Plan.String():
		return Plan, nil
	case Apply.String():
		return Apply, nil
	}
	return -1, fmt.Errorf("unsupported Terragrunt command: %s", name)
}

// Run executes the command "terragrunt run" with additional flags and
// environment variables to support common high-level use cases across multiple
// configuration units. The most important use cases are to plan and to apply
// infrastructure changes, of course.
func Run(ui ui.UserInterface, phase Phase, options Options) error {
	args := []string{"run", "--all"}

	if !ui.IsInteractive() {
		args = append(args, "--non-interactive")
	}

	args = append(args, runQueueArgs(options.UpdatedFiles)...)

	args = append(args, "--json-out-dir", filepath.Join(".terrahoot/output"))

	args = append(args, phase.String(), "--")

	if !ui.IsInteractive() {
		args = append(args, "-input=false")

		if phase == Apply {
			args = append(args, "-auto-approve")
		}
	}

	cmd := exec.Command("terragrunt", args...)

	if ui.IsInteractive() {
		cmd.Stdin = os.Stdin
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func runQueueArgs(changedFiles []string) []string {
	args := []string{"--queue-strict-include"}

	for _, f := range changedFiles {
		if filepath.Base(f) == "terragrunt.hcl" {
			args = append(args, "--queue-include-dir", filepath.Dir(f))
		} else {
			args = append(args, "--queue-include-units-reading", f)
		}
	}

	return args
}

// IsUnitFile returns true if the given path points to a Terragrunt unit
// configuration file named "terragrunt.hcl", as opposed to any other HCL file
// that ends with the ".hcl" suffix and may or may not be read by Terragrunt.
func IsUnitFile(path string) bool {
	return filepath.Base(path) == "terragrunt.hcl"
}

// UnitDir returns the directory of the given Terragrunt unit configuration
// file. Use UnitDir in conjunction with IsUnitFile.
func UnitDir(path string) string {
	return filepath.Dir(path)
}
