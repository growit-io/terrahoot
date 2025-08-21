package cmd

import (
	"fmt"
	"log"

	"github.com/growit-io/terrahoot/internal/git"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(changedFilesCmd)
}

var changedFilesCmd = &cobra.Command{
	Use:    "changed-files [flags] [revision]",
	Short:  "List files changed since base revision",
	Args:   cobra.MaximumNArgs(1),
	Hidden: true,

	Long: `List files change since the given base revision.
	
The optional "revision" argument can be specified in any of the <rev> formats
accepted by git-rev-parse(1) and defaults to "` + git.DefaultGitRemote + `/HEAD".`,

	Run: func(cmd *cobra.Command, args []string) {
		baseRef := ""

		if len(args) > 0 {
			baseRef = args[0]
		}

		dir, err := git.New().Dir()
		if err != nil {
			log.Fatal(err)
		}

		files, err := dir.ChangedFiles(baseRef)
		if err != nil {
			log.Fatal(err)
		}

		// XXX: output would be easier to read if items were sorted by path
		for path, change := range files {
			fmt.Println(change, path)
		}
	},
}
