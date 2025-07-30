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

	RunE: func(cmd *cobra.Command, args []string) error {
		return runCmd.RunE(cmd, append([]string{"--all", "plan"}, args...))
	},
}
