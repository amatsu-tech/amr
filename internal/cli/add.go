package cli

import (
	"fmt"

	"github.com/amatsu-tech/amr/internal/cli/output"
	"github.com/amatsu-tech/amr/internal/repository"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <path>...",
	Short: "Add one or more repositories to the workspace",
	Args:  cobra.MinimumNArgs(1),
	RunE:  addCommand,
}

func addCommand(cmd *cobra.Command, args []string) error {
	root, err := currentWorkspace()
	if err != nil {
		return err
	}

	results, err := repository.Add(root, args...)
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

func init() {
	amrCmd.AddCommand(addCmd)
}
