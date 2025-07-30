package cmd

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/growit-io/terrahoot/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a command in Terragrunt units affected by changed files",

	RunE: func(cmd *cobra.Command, args []string) error {
		git := internal.Git{}
		baseRef := "" // TODO: look up base ref in environment variable

		changedFiles, err := git.ChangedFiles(baseRef)
		if err != nil {
			return err
		}

		terragruntArgs := []string{"run", "--queue-strict-include"}
		for _, f := range changedFiles {
			if filepath.Base(f) == "terragrunt.hcl" {
				terragruntArgs = append(terragruntArgs, "--queue-include-dir", filepath.Dir(f))
			} else {
				terragruntArgs = append(terragruntArgs, "--queue-include-units-reading", f)
			}
		}

		terragruntArgs = append(terragruntArgs, args...)
		terragruntCmd := exec.Command("terragrunt", terragruntArgs...)
		terragruntCmd.Stdin = os.Stdin
		terragruntCmd.Stdout = os.Stdout
		terragruntCmd.Stderr = os.Stderr
		return terragruntCmd.Run()
	},
}
