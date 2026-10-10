package cli

import (
	"slices"

	"github.com/amatsu-tech/amr/internal/cli/output"
	"github.com/amatsu-tech/amr/internal/repository"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan workspace for Git repositories",
	RunE:  scanCommand,
}

func scanCommand(cmd *cobra.Command, args []string) error {
	root, err := currentWorkspace()
	if err != nil {
		return err
	}

	managed, err := repository.List(root)
	if err != nil {
		return err
	}

	discovered, err := repository.Scan(root)
	if err != nil {
		return err
	}

	if len(discovered) == 0 {
		cmd.Println("No Git repositories found.")
		return nil
	}

	managedPaths := make([]string, len(managed))
	for i, repo := range managed {
		managedPaths[i] = repo.Path
	}

	var managedRepos, unmanagedRepos []repository.Repository

	for _, repo := range discovered {
		if slices.Contains(managedPaths, repo.Path) {
			managedRepos = append(managedRepos, repo)
			continue
		}

		unmanagedRepos = append(unmanagedRepos, repo)
	}

	printed := false

	if len(managedRepos) > 0 {
		cmd.Println(output.SuccessStyle.Render("Managed:"))

		for _, repo := range managedRepos {
			cmd.Printf("  %s\n", repo.Path)
		}

		printed = true
	}

	if len(unmanagedRepos) > 0 {
		if printed {
			cmd.Println()
		}

		cmd.Println(output.NoticeStyle.Render("Unmanaged:"))

		for _, repo := range unmanagedRepos {
			cmd.Printf("  %s\n", repo.Path)
		}
	}

	cmd.Println()
	cmd.Println(output.Summary(
		output.SummaryLine{Symbol: "✓", Style: output.SuccessStyle, Count: len(managedRepos), Label: "managed"},
		output.SummaryLine{Symbol: "•", Style: output.NoticeStyle, Count: len(unmanagedRepos), Label: "unmanaged"},
	))

	return nil
}

func init() {
	amrCmd.AddCommand(scanCmd)
}
