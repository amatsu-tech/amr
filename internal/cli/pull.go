package cli

import (
	"fmt"

	"github.com/amatsu-tech/amr/internal/cli/output"
	"github.com/amatsu-tech/amr/internal/repository"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Fast-forward every managed repository from its upstream",
	RunE:  pullCommand,
}

type namedPullResult struct {
	Name string
	repository.PullResult
}

func pullCommand(cmd *cobra.Command, args []string) error {
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

	var pulled, upToDate, skipped, failed []namedPullResult

	for _, repo := range repositories {
		result := namedPullResult{Name: repo.Path, PullResult: repository.Pull(repo)}

		switch result.Outcome {
		case repository.PullUpdated:
			pulled = append(pulled, result)
		case repository.PullUpToDate:
			upToDate = append(upToDate, result)
		case repository.PullSkipped:
			skipped = append(skipped, result)
		case repository.PullFailed:
			failed = append(failed, result)
		}
	}

	printed := false

	if len(pulled) > 0 {
		cmd.Println(output.SuccessStyle.Render("Pulled:"))

		for _, result := range pulled {
			cmd.Printf("  %s\n", result.Name)
		}

		printed = true
	}

	if len(upToDate) > 0 {
		if printed {
			cmd.Println()
		}

		cmd.Println(output.NoticeStyle.Render("Already up to date:"))

		for _, result := range upToDate {
			cmd.Printf("  %s\n", result.Name)
		}

		printed = true
	}

	if len(skipped) > 0 {
		if printed {
			cmd.Println()
		}

		cmd.Println(output.NoticeStyle.Render("Skipped:"))

		skippedTable := output.NewTable("Repository", "Reason")

		for _, result := range skipped {
			skippedTable.Row(result.Name, result.Reason)
		}

		cmd.Println(skippedTable.String())

		printed = true
	}

	if len(failed) > 0 {
		if printed {
			cmd.Println()
		}

		cmd.Println(output.ErrorStyle.Render("Failed:"))

		failedTable := output.NewTable("Repository", "Reason")

		for _, result := range failed {
			failedTable.Row(result.Name, result.Reason)
		}

		cmd.Println(failedTable.String())
	}

	cmd.Println()
	cmd.Println(output.Summary(
		output.SummaryLine{Symbol: "✓", Style: output.SuccessStyle, Count: len(pulled), Label: "pulled"},
		output.SummaryLine{Symbol: "•", Style: output.NoticeStyle, Count: len(upToDate), Label: "up to date"},
		output.SummaryLine{Symbol: "•", Style: output.NoticeStyle, Count: len(skipped), Label: "skipped"},
		output.SummaryLine{Symbol: "✗", Style: output.ErrorStyle, Count: len(failed), Label: "failed"},
	))

	if len(failed) > 0 {
		return output.Silent(fmt.Errorf("%d pulled, %d up to date, %d skipped, %d failed",
			len(pulled), len(upToDate), len(skipped), len(failed)))
	}

	return nil
}

func init() {
	amrCmd.AddCommand(pullCmd)
}
