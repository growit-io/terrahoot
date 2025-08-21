package cmd

import (
	"os"

	"github.com/growit-io/terrahoot/internal/workflow"
	"github.com/spf13/cobra"

	_ "github.com/growit-io/terrahoot/internal/platform/github"
	_ "github.com/growit-io/terrahoot/internal/platform/local"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

var debug bool

var rootCmd = &cobra.Command{
	Use:   "terrahoot",
	Short: "Opinionated wrapper for Terragrunt runs in CI/CD workflows",

	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// https://github.com/spf13/cobra/issues/340#issuecomment-374617413
		cmd.SilenceUsage = false
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		// End of command-line flags and argument handling
		cmd.SilenceUsage = true

		w := workflow.FromEnvironment(os.Environ()).
			WithDebug(debug)

		return w.Run()
	},
}

func init() {
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "enable debug output")
}
