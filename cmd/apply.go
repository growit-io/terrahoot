package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Run \"apply\" command in Terragrunt units affected by changed files",
	Args:  cobra.ExactArgs(0),

	RunE: func(cmd *cobra.Command, args []string) error {
		return runCmd.RunE(cmd, []string{"apply"})
	},
}
