package cli

import (
	"errors"
	"fmt"

	"github.com/amatsu-tech/amr/internal/cli/output"
	"github.com/amatsu-tech/amr/internal/repository"
	"github.com/spf13/cobra"
)

var addAll bool

var addCmd = &cobra.Command{
	Use:   "add [<path>...]",
	Short: "Add one or more repositories to the workspace",
	Long: `Add repositories to the workspace.

Use --all to add every repository found in the workspace that is not added yet.`,
	Example: `  amr add services/api web
  amr add --all`,
	Args: cobra.ArbitraryArgs,
	RunE: addCommand,
}

func addCommand(cmd *cobra.Command, args []string) error {
	root, err := currentWorkspace()
	if err != nil {
		return err
	}

	paths, err := pathsToAdd(cmd, root, args)
	if err != nil || len(paths) == 0 {
		return err
	}

	results, err := repository.Add(root, paths...)
	if err != nil {
		return err
	}

	var added, failed []repository.AddResult

	for _, result := range results {
		if result.Error != nil {
			failed = append(failed, result)
			continue
		}

		added = append(added, result)
	}

	if len(added) > 0 {
		cmd.Println(output.SuccessStyle.Render("Added:"))

		for _, result := range added {
			cmd.Printf("  %s\n", result.Path)
		}
	}

	if len(failed) > 0 {
		if len(added) > 0 {
			cmd.Println()
		}

		cmd.Println(output.ErrorStyle.Render("Failed:"))

		failedTable := output.NewTable("Repository", "Reason")

		for _, result := range failed {
			failedTable.Row(result.Path, result.Error.Error())
		}

		cmd.Println(failedTable.String())
	}

	cmd.Println()
	cmd.Println(output.Summary(
		output.SummaryLine{Symbol: "✓", Style: output.SuccessStyle, Count: len(added), Label: "added"},
		output.SummaryLine{Symbol: "✗", Style: output.ErrorStyle, Count: len(failed), Label: "failed"},
	))

	if len(failed) > 0 {
		return output.Silent(fmt.Errorf("%d added, %d failed", len(added), len(failed)))
	}

	return nil
}

// pathsToAdd решает, что добавлять: пути из аргументов или все найденные (--all).
// Пустой результат без ошибки — добавлять нечего.
func pathsToAdd(cmd *cobra.Command, root string, args []string) ([]string, error) {
	if addAll && len(args) > 0 {
		return nil, errors.New("specify paths or --all, not both")
	}

	if len(args) > 0 {
		return args, nil
	}

	if !addAll {
		return nil, output.WithHint(
			errors.New("no repositories specified"),
			"Add specific repositories:\n  amr add <path>...\n\nOr add every repository found in the workspace:\n  amr add --all",
		)
	}

	unmanaged, err := repository.Unmanaged(root)
	if err != nil {
		return nil, err
	}

	if len(unmanaged) == 0 {
		cmd.Println("No new repositories found.")
		return nil, nil
	}

	paths := make([]string, len(unmanaged))
	for i, repo := range unmanaged {
		paths[i] = repo.Directory
	}

	return paths, nil
}

func init() {
	addCmd.Flags().BoolVar(&addAll, "all", false, "add every repository found in the workspace")
	amrCmd.AddCommand(addCmd)
}
