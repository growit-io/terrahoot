package cmd

import (
	"github.com/spf13/cobra"
)

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

var rootCmd = &cobra.Command{
	Use:   "terrahoot",
	Short: "Opinionated wrapper for Terragrunt runs in CI/CD workflows",
	// TODO: provide a `Long` description linking to the project homepage
}
