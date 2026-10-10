package cli

import (
	"github.com/amatsu-tech/amr/internal/repository"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List repositories in the workspace",
	RunE:  listCommand,
}

func listCommand(cmd *cobra.Command, args []string) error {
	root, err := currentWorkspace()
	if err != nil {
		return err
	}

	repositories, err := repository.List(root)
	if err != nil {
		return err
	}

	if len(repositories) == 0 {
		cmd.Println("No repositories have been added to this workspace.")
		return nil
	}

	for _, repository := range repositories {
		cmd.Println(repository.Path)
	}

	return nil
}

func init() {
	amrCmd.AddCommand(listCmd)
}
