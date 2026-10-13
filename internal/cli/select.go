package cli

import (
	"github.com/amatsu-tech/amr/internal/repository"
	"github.com/spf13/cobra"
)

// Значение флага --group у команд, которые работают с выбранными репозиториями.
var groups []string

func addGroupFlag(cmd *cobra.Command, usage string) {
	cmd.Flags().StringSliceVarP(&groups, "group", "g", nil, usage)
}

// addGroupFilterFlag — флаг --group для команд, которые выбирают репозитории по группам.
func addGroupFilterFlag(cmd *cobra.Command) {
	addGroupFlag(cmd, "only repositories in these groups (comma-separated `names` or repeated)")
}

// Сообщение, когда в workspace ещё нет ни одного репозитория.
const noRepositoriesMessage = "No repositories have been added to this workspace."

// selectedRepositories — репозитории workspace, выбранные путями из args и флагом --group.
// Без путей и групп — все.
func selectedRepositories(root string, paths []string) ([]repository.Repository, error) {
	repositories, err := repository.List(root)
	if err != nil {
		return nil, err
	}

	return repository.Select(root, repositories, paths, groups)
}
