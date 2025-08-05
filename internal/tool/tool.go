package tool

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Tool is basically a wrapper around os/exec.Command() that implements the
// builder pattern with certain defaults and output handling functions that are
// more suitable for the main program that uses this package.
type Tool struct {
	executable string
	args       []string
	dir        string
	stderr     io.Writer
}

// New prepares to invoke the given system executable command as a tool within
// the main program that uses this package.
//
// The standard error stream of the tool is connected to the operating system's
// standard error stream instead of /dev/null, by default. This behaviour is
// different from os/exec.Command.
//
// The returned Tool instance is immutable. Create modified copies of the
// initial Tool instance using the so-called "builder" pattern, by chaining
// function calls that each return a modified copy.
func New(executable string) *Tool {
	return &Tool{
		executable: executable,
		stderr:     os.Stderr,
	}
}

// WithArgs returns a new immutable Tool instance which invokes the previously
// specified executable with the given arguments.
func (t *Tool) WithArgs(arg ...string) *Tool {
	copy := *t
	copy.args = arg
	return &copy
}

// WithDir returns a new immutable Tool instance which invokes the executable
// within the given directory as the current working directory. If `dir` is the
// empty string or ".", then it is interpreted as the current working directory
// of the system process when the command is actually invoked.
func (t *Tool) WithDir(dir string) *Tool {
	copy := *t
	copy.dir = dir
	return &copy
}

// WithStderr allows redirecting stderr to /dev/null by setting it to nil or
// capturing it by setting it to an io.Writer instance.
func (t *Tool) WithStderr(stderr io.Writer) *Tool {
	copy := *t
	copy.stderr = stderr
	return &copy
}

// Command invokes the tool using the os/exec.Command function and returns the
// os/exec.Cmd object.
func (t *Tool) Command() *exec.Cmd {
	cmd := exec.Command(t.executable, t.args...)
	cmd.Stderr = t.stderr
	cmd.Dir = t.dir
	return cmd
}

func (t *Tool) Run() error {
	if err := t.Command().Run(); err != nil {
		return t.wrapError(err)
	}
	return nil
}

func (t *Tool) Output() ([]byte, error) {
	bytes, err := t.Command().Output()
	if err != nil {
		return nil, t.wrapError(err)
	}
	return bytes, nil
}

func (t *Tool) OutputLine() (string, error) {
	lines, err := t.OutputLines()
	if err != nil {
		return "", err
	}

	switch {
	case len(lines) < 1:
		return "", t.newError("expected one line of output, but got none")
	case len(lines) > 1:
		return "", t.newError("expected one line of output, but got multiple")
	}

	return lines[0], nil
}

func (t *Tool) OutputLines() ([]string, error) {
	outputBytes, err := t.Output()
	if err != nil {
		return nil, err
	}

	var outputLines []string
	scanner := bufio.NewScanner(bytes.NewReader(outputBytes))

	for scanner.Scan() {
		line := scanner.Text()
		outputLines = append(outputLines, line)
	}

	return outputLines, nil
}

func (t *Tool) wrapError(err error) error {
	return fmt.Errorf("command %s with arguments %s: %w", t.executable, t.args, err)
}

func (t *Tool) newError(msg string) error {
	return fmt.Errorf("command %s with arguments %s: %s", t.executable, t.args, msg)
}
