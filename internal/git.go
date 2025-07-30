package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

// executable is the name of the "git" executable to call.
const executable = "git"

// defaultRemote is the name of the remote that is used when resolving the
// default base reference.
const defaultRemote = "origin"

// Git provides methods to inspect a git repository. It is basically a wrapper
// around the "git" command for very specific use cases.
type Git struct {
	// Optionally add fields here, e.g. RepoPath string
}

// ChangedFiles returns a sorted, deduplicated list of files that have changed
// compared to the provided baseRef (or "origin/HEAD" if empty).
func (g *Git) ChangedFiles(baseRef string) ([]string, error) {
	if baseRef == "" {
		// The default base reference should point to the default branch of the
		// remote repository. Of course, this assumes that the local repository
		// *does* have a remote repository that it was cloned from.
		baseRef = defaultRemote + "/HEAD"

		// The "HEAD" reference exists when a repository gets cloned initially
		// from the remote, but needs to be discovered explicitly if the remote
		// gets added to an existing repository.
		_, err := g.revParse(baseRef)
		if err != nil {
			_, err := g.commandOutputLines("remote", "set-head", defaultRemote, "-a")
			if err != nil {
				return nil, err
			}
		}
	}

	paths := make(map[string]struct{})

	// List files modified in the local working copy, except ignored. Including
	// ignored files can quickly cause the list to get too long, since Terraform
	// and Terragrunt generate a lot of temporary files.
	statusLines, err := g.commandOutputLines("status", "--short", "--untracked-files", ".")
	if err != nil {
		return nil, err
	}

	for _, line := range statusLines {
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		paths[parts[1]] = struct{}{}
	}

	// List files that differ between the two committed trees. Note that we're
	// not giving any special attention to renames, currently.
	diffLines, err := g.commandOutputLines("diff", "--name-status", baseRef, "--", ".")
	if err != nil {
		return nil, err
	}

	for _, line := range diffLines {
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		paths[parts[1]] = struct{}{}
	}

	// Collect, sort, and return unique paths.
	var uniquePaths []string
	for path := range paths {
		uniquePaths = append(uniquePaths, path)
	}
	sort.Strings(uniquePaths)
	return uniquePaths, nil
}

func (g *Git) revParse(ref string) ([]string, error) {
	return g.commandOutputLines("rev-parse", ref)
}

func (g *Git) commandOutputLines(args ...string) ([]string, error) {
	cmd := exec.Command(executable, args...)
	outputBytes, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running \"%s\" command with arguments %s: %w", executable, args, err)
	}

	var outputLines []string
	scanner := bufio.NewScanner(bytes.NewReader(outputBytes))

	for scanner.Scan() {
		line := scanner.Text()
		outputLines = append(outputLines, line)
	}

	return outputLines, nil
}
