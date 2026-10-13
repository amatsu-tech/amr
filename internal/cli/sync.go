package cli

import (
	"fmt"
	"slices"

	"charm.land/lipgloss/v2"
	"github.com/amatsu-tech/amr/internal/cli/output"
	"github.com/amatsu-tech/amr/internal/parallel"
	"github.com/amatsu-tech/amr/internal/repository"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync [<path>...]",
	Short: "Update the workspace config, clone missing repositories and pull the rest",
	Long: `Bring the workspace up to date in one command.

If the workspace folder is itself a Git repository (a shared config repository),
amr first pulls it to get the latest config. Then it clones repositories that are
missing on disk and pulls the rest.`,
	RunE: syncCommand,
}

type namedSyncResult struct {
	Path string
	repository.SyncResult
}

func syncCommand(cmd *cobra.Command, args []string) error {
	root, err := currentWorkspace()
	if err != nil {
		return err
	}

	// Пути из конфига до обновления — чтобы показать, сколько репозиториев добавилось.
	before, err := repository.List(root)
	if err != nil {
		return err
	}

	workspaceFailed := printWorkspaceUpdate(cmd, repository.UpdateWorkspace(root), root, before)

	repositories, err := selectedRepositories(root, args)
	if err != nil {
		return err
	}

	if len(repositories) == 0 {
		cmd.Println(noRepositoriesMessage)
		return nil
	}

	progress := output.NewProgress("Syncing", len(repositories))

	results := parallel.Map(repositories, jobs, func(repo repository.Repository) namedSyncResult {
		defer progress.Advance()

		return namedSyncResult{Path: repo.Path, SyncResult: repository.Sync(repo)}
	})

	progress.Finish()

	var cloned, updated, upToDate, skipped, failed []namedSyncResult

	for _, result := range results {
		switch result.Outcome {
		case repository.SyncCloned:
			cloned = append(cloned, result)
		case repository.SyncUpdated:
			updated = append(updated, result)
		case repository.SyncUpToDate:
			upToDate = append(upToDate, result)
		case repository.SyncSkipped:
			skipped = append(skipped, result)
		case repository.SyncFailed:
			failed = append(failed, result)
		}
	}

	printed := false

	printList := func(title string, style lipgloss.Style, items []namedSyncResult) {
		if len(items) == 0 {
			return
		}

		if printed {
			cmd.Println()
		}

		cmd.Println(style.Render(title))

		for _, item := range items {
			cmd.Printf("  %s\n", item.Path)
		}

		printed = true
	}

	printTable := func(title string, style lipgloss.Style, items []namedSyncResult) {
		if len(items) == 0 {
			return
		}

		if printed {
			cmd.Println()
		}

		cmd.Println(style.Render(title))

		table := output.NewTable("Repository", "Reason")

		for _, item := range items {
			table.Row(item.Path, item.Reason)
		}

		cmd.Println(table.String())

		printed = true
	}

	printList("Cloned:", output.SuccessStyle, cloned)
	printList("Pulled:", output.SuccessStyle, updated)
	printList("Already up to date:", output.NoticeStyle, upToDate)
	printTable("Skipped:", output.NoticeStyle, skipped)
	printTable("Failed:", output.ErrorStyle, failed)

	cmd.Println()
	cmd.Println(output.Summary(
		output.SummaryLine{Symbol: "✓", Style: output.SuccessStyle, Count: len(cloned), Label: "cloned"},
		output.SummaryLine{Symbol: "✓", Style: output.SuccessStyle, Count: len(updated), Label: "pulled"},
		output.SummaryLine{Symbol: "•", Style: output.NoticeStyle, Count: len(upToDate), Label: "up to date"},
		output.SummaryLine{Symbol: "•", Style: output.NoticeStyle, Count: len(skipped), Label: "skipped"},
		output.SummaryLine{Symbol: "✗", Style: output.ErrorStyle, Count: len(failed), Label: "failed"},
	))

	if len(failed) > 0 || workspaceFailed {
		return output.Silent(fmt.Errorf("%d repositories failed to sync", len(failed)))
	}

	return nil
}

// printWorkspaceUpdate печатает, что стало с самим workspace, и возвращает true при ошибке.
// Если workspace не git-репозиторий, ничего не печатает.
func printWorkspaceUpdate(cmd *cobra.Command, result repository.WorkspaceResult, root string, before []repository.Repository) bool {
	switch result.Outcome {
	case repository.WorkspaceNotRepository:
		return false
	case repository.WorkspaceUpdated:
		message := "Workspace updated"

		if added := newRepositories(root, before); added > 0 {
			message += fmt.Sprintf(" (%d new in config)", added)
		}

		cmd.Println(output.SuccessStyle.Render(message))
	case repository.WorkspaceUpToDate:
		cmd.Println(output.NoticeStyle.Render("Workspace already up to date"))
	case repository.WorkspaceSkipped:
		cmd.Println(output.NoticeStyle.Render("Workspace not updated: " + result.Reason))
	case repository.WorkspaceFailed:
		cmd.Println(output.ErrorStyle.Render("Workspace update failed: " + result.Reason))
	}

	cmd.Println()

	return result.Outcome == repository.WorkspaceFailed
}

// newRepositories считает, сколько репозиториев появилось в конфиге после обновления.
func newRepositories(root string, before []repository.Repository) int {
	after, err := repository.List(root)
	if err != nil {
		return 0
	}

	added := 0

	for _, repo := range after {
		if !slices.ContainsFunc(before, func(old repository.Repository) bool { return old.Path == repo.Path }) {
			added++
		}
	}

	return added
}

func init() {
	addGroupFilterFlag(syncCmd)
	amrCmd.AddCommand(syncCmd)
}
