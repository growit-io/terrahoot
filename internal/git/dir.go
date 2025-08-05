package git

import (
	"fmt"
	"os"
	"strings"
)

// DefaultGitRemote is the default name of the remote that is used when cloning
// a remote repository.
const DefaultGitRemote = "origin"

// Dir represents a Git work tree, which usually also contains a local copy
// of the repository in the ".git" subdirectory in its top-level directory.
type Dir struct {
	git      Git
	toplevel string
}

// NewDir constructs a new git.Dir instance which represents a valid local work
// tree.
func NewDir(git *Git) (*Dir, error) {
	toplevel, err := git.command("rev-parse", "--show-toplevel").OutputLine()
	if err != nil {
		return nil, err
	}

	return &Dir{
		git:      *git,
		toplevel: toplevel,
	}, nil
}

// Toplevel returns the top-level directory of the current Git work tree. Note
// that the current working directory within the work tree may be different.
func (d *Dir) Toplevel() string {
	return d.toplevel
}

// IsClean determines whether the work tree is clean. A clean work tree does not
// contain uncommitted changes or unmerged conflicts.
func (d *Dir) IsClean() (bool, error) {
	lines, err := d.git.command("status", "--short").OutputLines()
	if err != nil {
		return false, err
	}

	return len(lines) == 0, nil
}

// RemoteBranchRef determines whether the work tree is connected to a remote
// branch and returns the remote branch reference.
func (d *Dir) RemoteBranchRef() (string, error) {
	remoteBranchRef, err := d.git.command("rev-parse",
		"--abbrev-ref", "--symbolic-full-name", "@{u}").
		WithStderr(nil).
		OutputLine()
	if err != nil {
		return "", fmt.Errorf("unable to determine remote tracking branch: %w", err)
	}

	return remoteBranchRef, nil
}

// IsUpToDate determines whether the work tree is connected to a remote branch
// and currently up-to-date with that branch.
func (d *Dir) IsUpToDate() (bool, error) {
	if err := d.git.command("fetch", "--quiet").Run(); err != nil {
		return false, fmt.Errorf("unable to update remote branch references: %w", err)
	}

	remoteBranchRef, err := d.RemoteBranchRef()
	if err != nil {
		return false, nil
	}

	aheadCount, err := d.git.command("rev-list", "--count", "@{u}..HEAD").OutputLine()
	if err != nil {
		return false, fmt.Errorf("unable to count commits ahead of remote branch %s: %w", remoteBranchRef, err)
	}

	behindCount, err := d.git.command("rev-list", "--count", "HEAD..@{u}").OutputLine()
	if err != nil {
		return false, fmt.Errorf("unable to count commits behind remote branch %s: %w", remoteBranchRef, err)
	}

	return aheadCount == "0" && behindCount == "0", nil
}

type FileChange string

const (
	Added    FileChange = "A"
	Deleted  FileChange = "D"
	Modified FileChange = "M"
)

// ChangedFiles returns the paths of files relative to the working copy's
// toplevel directory that have changed compared to the provided baseRef.
// Compares against "origin/HEAD" if baseRef is the empty string. Ingores
// all other changes in the local working copy, steged and unstaged.
func (d *Dir) ChangedFiles(baseRef string) (map[string]FileChange, error) {
	files := make(map[string]FileChange)

	// The default base reference should point to the default branch of the
	// remote repository. Of course, this assumes that the local repository
	// *does* have a remote repository that it was cloned from.
	if baseRef == "" {
		ref, err := d.defaultRemoteBranchRef()
		if err != nil {
			return nil, err
		}
		baseRef = ref
	}

	// List files that differ between the two committed trees.
	diffLines, err := d.git.command("diff", "--no-renames", "--name-status", baseRef, "--", ".").OutputLines()
	if err != nil {
		return nil, err
	}

	for _, line := range diffLines {
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		// See --diff-filter in git-diff(1) for the meaning of each constant.
		var change FileChange
		switch parts[0] {
		case "A":
			change = Added
		case "D":
			change = Deleted
		case "M":
			change = Modified
		default:
			// FIXME: handle all relevant status latters!
			return nil, fmt.Errorf("unhandled status letter in \"git diff --name-status\" output: %s", parts[0])
		}
		files[parts[1]] = change
	}

	return files, nil
}

// defaultRemoteBranchRef returns a symbolic reference that points to the
// default branch of the primary remote that the repository has been cloned
// from. Typicalled the remote is named "origin" and the name of the default
// branch reference is "HEAD.""
func (d *Dir) defaultRemoteBranchRef() (string, error) {
	ref := DefaultGitRemote + "/HEAD"

	// Resolve the path to the file where the remote head reference would
	// normally be stored, respecting all environment and config variables.
	gitPath, err := d.revParseGitPath("refs/remotes/" + ref)
	if err != nil {
		return "", err
	}

	// The remote's default branch information is not present in this local
	// repsoitory. Try to fetch it now from the actual remote.
	_, err = os.Stat(gitPath)
	if err != nil {
		// The "HEAD" reference exists when a repository gets cloned
		// initially from the remote, but needs to be discovered explicitly
		// if the remote gets added to an existing repository.
		err := d.git.command("remote", "set-head", DefaultGitRemote, "-a").Run()
		if err != nil {
			return "", err
		}
	}

	// The reference should now be valid and the gitPath should exist, but we
	// forego the os.Stat() call to verify it.
	return ref, nil
}

// revParseGitPath returns the first line of the output of the "git rev-parse
// --git-path" command. It returns a relative path from the current working
// directory to the given path relative to the ".git" repository directory.
func (d *Dir) revParseGitPath(gitPath string) (string, error) {
	relPath, err := d.git.command("rev-parse", "--git-path", gitPath).OutputLine()
	if err != nil {
		return "", err
	}

	return relPath, nil
}
