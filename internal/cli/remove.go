package cli

import (
	"fmt"

	"github.com/amatsu-tech/amr/internal/cli/output"
	"github.com/amatsu-tech/amr/internal/repository"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <repository>...",
	Short: "Remove one or more repositories from the workspace",
	Args:  cobra.MinimumNArgs(1),
	RunE:  removeCommand,
}

func removeCommand(cmd *cobra.Command, args []string) error {
	root, err := currentWorkspace()
	if err != nil {
		return err
	}

	results, err := repository.Remove(root, args...)
	if err != nil {
		return err
	}

	var removed, failed []repository.RemoveResult

	for _, result := range results {
		if result.Error != nil {
			failed = append(failed, result)
			continue
		}

		removed = append(removed, result)
	}

	if len(removed) > 0 {
		cmd.Println(output.SuccessStyle.Render("Removed:"))

		for _, result := range removed {
			cmd.Printf("  %s\n", result.Name)
		}
	}

	if len(failed) > 0 {
		if len(removed) > 0 {
			cmd.Println()
		}

		cmd.Println(output.ErrorStyle.Render("Failed:"))

		failedTable := output.NewTable("Repository", "Reason")

		for _, result := range failed {
			failedTable.Row(result.Name, result.Error.Error())
		}

		cmd.Println(failedTable.String())
	}

	cmd.Println()
	cmd.Println(output.Summary(
		output.SummaryLine{Symbol: "✓", Style: output.SuccessStyle, Count: len(removed), Label: "removed"},
		output.SummaryLine{Symbol: "✗", Style: output.ErrorStyle, Count: len(failed), Label: "failed"},
	))

	if len(failed) > 0 {
		return output.Silent(fmt.Errorf("%d removed, %d failed", len(removed), len(failed)))
	}

	return nil
}

func init() {
	amrCmd.AddCommand(removeCmd)
}
