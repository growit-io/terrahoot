package workflow

import (
	"fmt"
	"os"

	"github.com/growit-io/terrahoot/internal/git"
	"github.com/growit-io/terrahoot/internal/terragrunt"
	"github.com/growit-io/terrahoot/internal/ui"
)

func Run(phase terragrunt.Phase) error {
	dir, err := git.New().Dir()
	if err != nil {
		return err
	}

	clean, err := dir.IsClean()
	if err != nil {
		return fmt.Errorf("cannot determine whether git working copy is clean: %w", err)
	}
	if !clean {
		return fmt.Errorf("git working copy appears to be dirty; make sure it is clean")
	}

	baseRef := os.Getenv("GIT_BASE_REF")
	changedFiles, err := dir.ChangedFiles(baseRef)
	if err != nil {
		return err
	}

	var deletedUnits []string
	var updatedFiles []string

	for path, change := range changedFiles {
		switch change {
		case git.Deleted:
			if terragrunt.IsUnitFile(path) {
				deletedUnits = append(deletedUnits, terragrunt.UnitDir(path))
			} else {
				updatedFiles = append(updatedFiles, path)
			}
		case git.Added, git.Modified:
			updatedFiles = append(updatedFiles, path)
		default:
			panic(fmt.Sprintf("unhandled file change in git repository: %s %s", change, path))
		}
	}

	options := terragrunt.Options{
		DeletedUnits: deletedUnits,
		UpdatedFiles: updatedFiles,
	}

	return terragrunt.Run(ui.New(), phase, options)
}
