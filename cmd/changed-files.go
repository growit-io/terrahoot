package cmd

import (
	"fmt"
	"log"

	"github.com/growit-io/terrahoot/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(changedFilesCmd)
}

var changedFilesCmd = &cobra.Command{
	Use:   "changed-files",
	Short: "List files changed since base revision",

	Run: func(cmd *cobra.Command, args []string) {
		git := &internal.Git{}

		files, err := git.ChangedFiles("")
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {
			fmt.Println(f)
		}
	},
}
