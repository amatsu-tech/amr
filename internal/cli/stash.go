package cli

import (
	"errors"
	"fmt"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/amatsu-tech/amr/internal/cli/output"
	"github.com/amatsu-tech/amr/internal/parallel"
	"github.com/amatsu-tech/amr/internal/repository"
	"github.com/spf13/cobra"
)

// Значение флага --message у amr stash.
var stashName string

var stashCmd = &cobra.Command{
	Use:   "stash [<path>...]",
	Short: "Stash changes in repositories, including new files",
	Long: `Stash changes in repositories under one name, including new (untracked) files.
Ignored files are left alone.

Without -m, a name is generated from the current time.`,
	Example: `  amr stash -m login-fix
  amr stash . -m login-fix
  amr stash pop login-fix`,
	RunE: stashCommand,
}

var stashListCmd = &cobra.Command{
	Use:   "list [<path>...]",
	Short: "List amr stashes by name",
	RunE:  stashListCommand,
}

var stashPopCmd = &cobra.Command{
	Use:   "pop [<name> [<path>...]]",
	Short: "Restore a stash and remove it",
	Long: `Restore changes from the stash with this name and remove it.
Without a name, the latest amr stash in each repository is restored.`,
	Example: `  amr stash pop
  amr stash pop login-fix
  amr stash pop login-fix .`,
	RunE: stashPopCommand,
}

var stashDropCmd = &cobra.Command{
	Use:   "drop <name> [<path>...]",
	Short: "Remove a stash without restoring it",
	Args:  cobra.MinimumNArgs(1),
	RunE:  stashDropCommand,
}

type namedStashResult struct {
	Path string
	repository.StashResult
}

type namedStashActionResult struct {
	Path string
	repository.StashActionResult
}

func stashCommand(cmd *cobra.Command, args []string) error {
	root, err := currentWorkspace()
	if err != nil {
		return err
	}

	name := strings.TrimSpace(stashName)
	if name == "" {
		name = repository.DefaultStashName(time.Now())
	}

	if strings.ContainsAny(name, "\n\r") {
		return errors.New("stash name must be a single line")
	}

	repositories, err := selectedRepositories(root, args)
	if err != nil {
		return err
	}

	if len(repositories) == 0 {
		cmd.Println(noRepositoriesMessage)
		return nil
	}

	results := parallel.Map(repositories, jobs, func(repo repository.Repository) namedStashResult {
		return namedStashResult{Path: repo.Path, StashResult: repository.StashSave(repo, name)}
	})

	var stashed, unchanged, failed []namedStashResult

	for _, result := range results {
		switch result.Outcome {
		case repository.StashSaved:
			stashed = append(stashed, result)
		case repository.StashNoChanges:
			unchanged = append(unchanged, result)
		case repository.StashFailed:
			failed = append(failed, result)
		}
	}

	if len(stashed) == 0 && len(failed) == 0 {
		cmd.Println("Nothing to stash.")
		return nil
	}

	if len(stashed) > 0 {
		cmd.Println(output.SuccessStyle.Render("Stashed as " + name + ":"))

		for _, result := range stashed {
			cmd.Printf("  %s\n", result.Path)
		}
	}

	if len(failed) > 0 {
		if len(stashed) > 0 {
			cmd.Println()
		}

		cmd.Println(output.ErrorStyle.Render("Failed:"))

		table := output.NewTable("Repository", "Reason")
		for _, result := range failed {
			table.Row(result.Path, result.Reason)
		}

		cmd.Println(table.String())
	}

	cmd.Println()
	cmd.Println(output.Summary(
		output.SummaryLine{Symbol: "✓", Style: output.SuccessStyle, Count: len(stashed), Label: "stashed"},
		output.SummaryLine{Symbol: "•", Style: output.NoticeStyle, Count: len(unchanged), Label: "without changes"},
		output.SummaryLine{Symbol: "✗", Style: output.ErrorStyle, Count: len(failed), Label: "failed"},
	))

	if len(stashed) > 0 {
		cmd.Println()
		cmd.Printf("Restore with:\n  amr stash pop %s\n", name)
	}

	if len(failed) > 0 {
		return output.Silent(fmt.Errorf("%d repositories failed to stash", len(failed)))
	}

	return nil
}

func stashListCommand(cmd *cobra.Command, args []string) error {
	root, err := currentWorkspace()
	if err != nil {
		return err
	}

	repositories, err := selectedRepositories(root, args)
	if err != nil {
		return err
	}

	type repositoryStashes struct {
		Path  string
		Names []string
		Err   error
	}

	results := parallel.Map(repositories, jobs, func(repo repository.Repository) repositoryStashes {
		names, err := repository.StashNames(repo)
		return repositoryStashes{Path: repo.Path, Names: names, Err: err}
	})

	// Имя стеша → репозитории, где он есть. Порядок — как встретились.
	var order []string
	byName := map[string][]string{}

	for _, result := range results {
		for _, name := range result.Names {
			if _, seen := byName[name]; !seen {
				order = append(order, name)
			}

			byName[name] = append(byName[name], result.Path)
		}
	}

	if len(order) == 0 {
		cmd.Println("No amr stashes.")
	} else {
		table := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
		for _, name := range order {
			fmt.Fprintf(table, "%s\t%s\n", name, strings.Join(byName[name], ", "))
		}

		if err := table.Flush(); err != nil {
			return err
		}
	}

	var failed int

	for _, result := range results {
		if result.Err != nil {
			failed++
			cmd.PrintErrln(output.ErrorStyle.Render(result.Path + ": " + result.Err.Error()))
		}
	}

	if failed > 0 {
		return output.Silent(fmt.Errorf("%d repositories could not be read", failed))
	}

	return nil
}

func stashPopCommand(cmd *cobra.Command, args []string) error {
	name, paths := "", args
	if len(args) > 0 {
		name, paths = args[0], args[1:]
	}

	return stashAction(cmd, name, paths, "Restored", repository.StashPop)
}

func stashDropCommand(cmd *cobra.Command, args []string) error {
	return stashAction(cmd, args[0], args[1:], "Dropped", repository.StashDrop)
}

// stashAction выполняет pop или drop стеша name в выбранных репозиториях.
func stashAction(cmd *cobra.Command, name string, paths []string, done string, action func(repository.Repository, string) repository.StashActionResult) error {
	root, err := currentWorkspace()
	if err != nil {
		return err
	}

	repositories, err := selectedRepositories(root, paths)
	if err != nil {
		return err
	}

	results := parallel.Map(repositories, jobs, func(repo repository.Repository) namedStashActionResult {
		return namedStashActionResult{Path: repo.Path, StashActionResult: action(repo, name)}
	})

	var succeeded, failed []namedStashActionResult

	for _, result := range results {
		switch result.Outcome {
		case repository.StashDone:
			succeeded = append(succeeded, result)
		case repository.StashActionFailed:
			failed = append(failed, result)
		}
	}

	if len(succeeded) == 0 && len(failed) == 0 {
		if name == "" {
			return output.WithHint(errors.New("no amr stashes found"), "See stashes with:\n  amr stash list")
		}

		return output.WithHint(fmt.Errorf("no amr stash named %s", name), "See stashes with:\n  amr stash list")
	}

	if len(succeeded) > 0 {
		cmd.Println(output.SuccessStyle.Render(done + ":"))

		for _, result := range succeeded {
			cmd.Printf("  %s\n", result.Path)
		}
	}

	if len(failed) > 0 {
		if len(succeeded) > 0 {
			cmd.Println()
		}

		cmd.Println(output.ErrorStyle.Render("Failed:"))

		table := output.NewTable("Repository", "Reason")
		for _, result := range failed {
			table.Row(result.Path, result.Reason)
		}

		cmd.Println(table.String())
	}

	if len(failed) > 0 {
		return output.Silent(fmt.Errorf("%d repositories failed", len(failed)))
	}

	return nil
}

func init() {
	stashCmd.Flags().StringVarP(&stashName, "message", "m", "", "stash name (default: generated from the current time)")
	addGroupFilterFlag(stashCmd)

	for _, sub := range []*cobra.Command{stashListCmd, stashPopCmd, stashDropCmd} {
		addGroupFilterFlag(sub)
		stashCmd.AddCommand(sub)
	}

	amrCmd.AddCommand(stashCmd)
}
