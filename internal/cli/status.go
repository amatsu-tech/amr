package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/amatsu-tech/amr/internal/cli/output"
	"github.com/amatsu-tech/amr/internal/git"
	"github.com/amatsu-tech/amr/internal/parallel"
	"github.com/amatsu-tech/amr/internal/repository"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show workspace status",
	RunE:  statusCommand,
}

type statusEntry struct {
	Name    string
	Summary git.StatusSummary
}

type repositoryStatus struct {
	Path    string
	URL     string
	Summary git.StatusSummary
	Err     error
}

type unavailableEntry struct {
	Name   string
	Reason string
}

func statusCommand(cmd *cobra.Command, args []string) error {
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

	var clean, attention []statusEntry
	var unavailable []unavailableEntry

	statuses := parallel.Map(repositories, jobs, func(repo repository.Repository) repositoryStatus {
		summary, err := git.Status(repo.Directory)
		return repositoryStatus{Path: repo.Path, URL: repo.URL, Summary: summary, Err: err}
	})

	// Сколько отсутствующих на диске репозиториев можно скачать через amr clone.
	cloneable := 0

	for _, status := range statuses {
		if errors.Is(status.Err, git.ErrDirectoryNotFound) && status.URL != "" {
			cloneable++
		}

		if status.Err != nil {
			unavailable = append(unavailable, unavailableEntry{Name: status.Path, Reason: status.Err.Error()})
			continue
		}

		summary := status.Summary
		entry := statusEntry{Name: status.Path, Summary: summary}

		if summary.Clean() && summary.Ahead == 0 && summary.Behind == 0 {
			clean = append(clean, entry)
			continue
		}

		attention = append(attention, entry)
	}

	printed := false

	if len(clean) > 0 {
		cmd.Println(output.SuccessStyle.Render("Clean:"))

		cleanTable := output.NewTable("Repository", "Branch")

		for _, item := range clean {
			cleanTable.Row(item.Name, item.Summary.Branch)
		}

		cmd.Println(cleanTable.String())

		printed = true
	}

	if len(attention) > 0 {
		if printed {
			cmd.Println()
		}

		cmd.Println("Needs attention:")

		attentionTable := output.NewTable("Repository", "Branch", "Status")

		for _, item := range attention {
			attentionTable.Row(item.Name, item.Summary.Branch, formatStatus(item.Summary))
		}

		cmd.Println(attentionTable.String())

		printed = true
	}

	if len(unavailable) > 0 {
		if printed {
			cmd.Println()
		}

		cmd.Println(output.ErrorStyle.Render("Unavailable:"))

		unavailableTable := output.NewTable("Repository", "Reason")

		for _, item := range unavailable {
			unavailableTable.Row(item.Name, item.Reason)
		}

		cmd.Println(unavailableTable.String())

		printed = true
	}

	if printed {
		cmd.Println()
	}

	cmd.Println(output.Summary(
		output.SummaryLine{Symbol: "✓", Style: output.SuccessStyle, Count: len(clean), Label: "clean"},
		output.SummaryLine{Symbol: "•", Style: output.NoticeStyle, Count: len(attention), Label: "need attention"},
		output.SummaryLine{Symbol: "✗", Style: output.ErrorStyle, Count: len(unavailable), Label: "unavailable"},
	))

	if cloneable > 0 {
		cmd.Println()
		cmd.Println("Clone missing repositories with:\n  amr clone")
	}

	if len(unavailable) > 0 {
		return output.Silent(fmt.Errorf("%d unavailable", len(unavailable)))
	}

	return nil
}

func formatStatus(summary git.StatusSummary) string {
	if summary.Clean() && summary.Ahead == 0 && summary.Behind == 0 {
		return output.SuccessStyle.Render("clean")
	}

	var parts []string

	if summary.Conflicted > 0 {
		parts = append(parts, fmt.Sprintf("%d conflicted", summary.Conflicted))
	}

	if summary.Staged > 0 {
		parts = append(parts, fmt.Sprintf("%d staged", summary.Staged))
	}

	if summary.Modified > 0 {
		parts = append(parts, fmt.Sprintf("%d modified", summary.Modified))
	}

	if summary.Untracked > 0 {
		parts = append(parts, fmt.Sprintf("%d untracked", summary.Untracked))
	}

	if summary.Ahead > 0 {
		parts = append(parts, fmt.Sprintf("%d ahead", summary.Ahead))
	}

	if summary.Behind > 0 {
		parts = append(parts, fmt.Sprintf("%d behind", summary.Behind))
	}

	text := strings.Join(parts, ", ")

	if summary.Conflicted > 0 {
		return output.ErrorStyle.Render(text)
	}

	return output.NoticeStyle.Render(text)
}

func init() {
	amrCmd.AddCommand(statusCmd)
}
