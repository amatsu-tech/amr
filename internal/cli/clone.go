package cli

import (
	"fmt"

	"github.com/amatsu-tech/amr/internal/cli/output"
	"github.com/amatsu-tech/amr/internal/parallel"
	"github.com/amatsu-tech/amr/internal/repository"
	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:   "clone [<path>...]",
	Short: "Clone repositories that are missing on disk",
	RunE:  cloneCommand,
}

type namedCloneResult struct {
	Path string
	repository.CloneResult
}

func cloneCommand(cmd *cobra.Command, args []string) error {
	root, err := currentWorkspace()
	if err != nil {
		return err
	}

	repositories, err := selectedRepositories(root, args)
	if err != nil {
		return err
	}

	if len(repositories) == 0 {
		cmd.Println("No repositories have been added to this workspace.")
		return nil
	}

	var cloned, exists, skipped, failed []namedCloneResult

	progress := output.NewProgress("Cloning", len(repositories))

	results := parallel.Map(repositories, jobs, func(repo repository.Repository) namedCloneResult {
		defer progress.Advance()

		return namedCloneResult{Path: repo.Path, CloneResult: repository.Clone(repo)}
	})

	progress.Finish()

	for _, result := range results {
		switch result.Outcome {
		case repository.CloneCloned:
			cloned = append(cloned, result)
		case repository.CloneExists:
			exists = append(exists, result)
		case repository.CloneSkipped:
			skipped = append(skipped, result)
		case repository.CloneFailed:
			failed = append(failed, result)
		}
	}

	printed := false

	if len(cloned) > 0 {
		cmd.Println(output.SuccessStyle.Render("Cloned:"))

		for _, result := range cloned {
			cmd.Printf("  %s\n", result.Path)
		}

		printed = true
	}

	if len(exists) > 0 {
		if printed {
			cmd.Println()
		}

		cmd.Println(output.NoticeStyle.Render("Already present:"))

		for _, result := range exists {
			cmd.Printf("  %s\n", result.Path)
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
			skippedTable.Row(result.Path, result.Reason)
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
			failedTable.Row(result.Path, result.Reason)
		}

		cmd.Println(failedTable.String())
	}

	cmd.Println()
	cmd.Println(output.Summary(
		output.SummaryLine{Symbol: "✓", Style: output.SuccessStyle, Count: len(cloned), Label: "cloned"},
		output.SummaryLine{Symbol: "•", Style: output.NoticeStyle, Count: len(exists), Label: "already present"},
		output.SummaryLine{Symbol: "•", Style: output.NoticeStyle, Count: len(skipped), Label: "skipped"},
		output.SummaryLine{Symbol: "✗", Style: output.ErrorStyle, Count: len(failed), Label: "failed"},
	))

	if len(failed) > 0 {
		return output.Silent(fmt.Errorf("%d cloned, %d already present, %d skipped, %d failed",
			len(cloned), len(exists), len(skipped), len(failed)))
	}

	return nil
}

func init() {
	addGroupFlag(cloneCmd, "only repositories in these groups (comma-separated `names` or repeated)")
	amrCmd.AddCommand(cloneCmd)
}
