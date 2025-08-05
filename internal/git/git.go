package git

import (
	"github.com/growit-io/terrahoot/internal/tool"
)

// DefaultGitExecutable is the default name of the Git executable to call.
const DefaultGitExecutable = "git"

// Git represents the executable system command that is typically named "git".
type Git struct {
	// executable is the name of the Git system command and defaults to the
	// value of DefaultGitExecutable.
	executable string

	// path is the name of a local working directory that Git should be executed
	// in. It is assumed that this directory is within a valid work tree that
	// also allows the repository to be discovered, but these assumptions may be
	// spoiled if the environment variables GIT_DIR or GIT_WORK_TREE are set,
	// for example.
	path string
}

// New constructs an immutable Git instance.
func New() *Git {
	return &Git{}
}

// WithExecutable returns a new immutable Git instance with the name of the Git
// executable command set to the given executable. It may be useful for testing
// or on systems where the Git executable name differs from the default name.
func (g *Git) WithExecutable(executable string) *Git {
	copy := *g
	copy.executable = executable
	return &copy
}

// WithPath returns a new immutable Git instance which executes all subcommands
// with the given path as the current working directory.
func (g *Git) WithPath(path string) *Git {
	copy := *g
	copy.path = path
	return &copy
}

// Dir returns an object that provides convenience methods to access a local Git
// repository's work tree as if "-C <path>" was specified for every Git command,
// where `path` is the current working directory, by default, or whatever other
// value was passed to the WithPath method.
func (g *Git) Dir() (*Dir, error) {
	return NewDir(g)
}

// command prepares the given "git" command-line for execution, much like the
// os/exec.Command() function, but using a slightly more convenient interface
// and default behaviour that is provided by the "tool" package.
func (g *Git) command(arg ...string) *tool.Tool {
	executable := g.executable
	if executable == "" {
		executable = DefaultGitExecutable
	}

	return tool.New(executable).WithArgs(arg...).WithDir(g.path)
}
