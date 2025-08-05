package cmd

import (
	"github.com/growit-io/terrahoot/internal/terragrunt"
	"github.com/growit-io/terrahoot/internal/workflow"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:    "run",
	Short:  "Run a command in Terragrunt units affected by changed files",
	Args:   cobra.ExactArgs(1),
	Hidden: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		command, err := terragrunt.ParsePhase(args[0])
		if err != nil {
			return err
		}

		cmd.SilenceUsage = true

		return workflow.Run(command)
	},
}
