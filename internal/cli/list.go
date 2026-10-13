package cli

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [<path>...]",
	Short: "List repositories in the workspace",
	RunE:  listCommand,
}

func listCommand(cmd *cobra.Command, args []string) error {
	root, err := currentWorkspace()
	if err != nil {
		return err
	}

	repositories, err := selectedRepositories(root, args)
	if err != nil {
		return err
	}

	if len(repositories) == 0 {
		cmd.Println(noRepositoriesMessage)
		return nil
	}

	table := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)

	for _, repo := range repositories {
		fmt.Fprintf(table, "%s\t%s\n", repo.Path, strings.Join(repo.Groups, ", "))
	}

	return table.Flush()
}

func init() {
	addGroupFilterFlag(listCmd)
	amrCmd.AddCommand(listCmd)
}
