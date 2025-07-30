package cmd

import (
	"fmt"
	"log"

	"github.com/growit-io/terrahoot/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(moduleFilesCmd)
}

var moduleFilesCmd = &cobra.Command{
	Use:   "module-files",
	Short: "List all source files of a Terraform module",
	Long:  "Takes the directory of a Terragrunt module as the first and only argument.",

	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		p := internal.ModuleParser{
			Dir: path,
		}

		m, err := p.ParseModule()
		if err != nil {
			log.Fatal(err)
		}

		err = printTFFiles(m)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func printTFFiles(m *internal.ModuleInfo) error {
	for _, path := range m.TFFiles {
		fmt.Println(path)
	}

	deps, err := m.Deps()
	if err != nil {
		return err
	}

	for _, dep := range deps {
		printTFFiles(dep)
	}

	return nil
}
