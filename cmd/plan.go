package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(planCmd)
}

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Run \"plan\" command in Terragrunt units affected by changed files",
	Args:  cobra.ExactArgs(0),

	RunE: func(cmd *cobra.Command, args []string) error {
		return runCmd.RunE(cmd, []string{"plan"})
	},
}
