package cli

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/amatsu-tech/amr/internal/cli/output"
	"github.com/amatsu-tech/amr/internal/parallel"
	"github.com/amatsu-tech/amr/internal/repository"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [<path>...] -- <command>...",
	Short: "Run a command in repositories",
	Example: `  amr run -- git fetch
  amr run -g backend -- go test ./...
  amr run frontend infra -- git log -1 --oneline`,
	Args: cobra.ArbitraryArgs,
	RunE: runCommand,
}

type runResult struct {
	Name   string
	Output string
	Error  error
}

func runCommand(cmd *cobra.Command, args []string) error {
	root, err := currentWorkspace()
	if err != nil {
		return err
	}

	// Всё до -- — пути репозиториев, после — команда.
	paths, command := []string{}, args
	if dash := cmd.ArgsLenAtDash(); dash > 0 {
		paths, command = args[:dash], args[dash:]
	}

	if len(command) == 0 {
		return errors.New("no command to run, put it after --")
	}

	repositories, err := selectedRepositories(root, paths)
	if err != nil {
		return err
	}

	if len(repositories) == 0 {
		cmd.Println(noRepositoriesMessage)
		return nil
	}

	name, commandArgs := command[0], command[1:]

	var succeeded, failed []runResult

	results := parallel.Map(repositories, jobs, func(repo repository.Repository) runResult {
		var buf bytes.Buffer

		runErr := repository.Run(repo, name, commandArgs, &buf, &buf)

		return runResult{Name: repo.Path, Output: strings.TrimSpace(buf.String()), Error: runErr}
	})

	for _, result := range results {
		if result.Error != nil {
			failed = append(failed, result)
			continue
		}

		succeeded = append(succeeded, result)
	}

	if len(succeeded) > 0 {
		cmd.Println(output.SuccessStyle.Render("Succeeded:"))

		for _, result := range succeeded {
			cmd.Printf("  %s\n", result.Name)
		}
	}

	if len(failed) > 0 {
		if len(succeeded) > 0 {
			cmd.Println()
		}

		for i, result := range failed {
			if i > 0 {
				cmd.Println()
			}

			cmd.Println(output.ErrorStyle.Render(result.Name))

			if result.Output != "" {
				cmd.Println(result.Output)
			} else {
				cmd.Println(result.Error.Error())
			}
		}

		cmd.Println()
		cmd.Println(output.ErrorStyle.Render("Failed:"))

		failedTable := output.NewTable("Repository", "Reason")

		for _, result := range failed {
			failedTable.Row(result.Name, result.Error.Error())
		}

		cmd.Println(failedTable.String())
	}

	cmd.Println()
	cmd.Println(output.Summary(
		output.SummaryLine{Symbol: "✓", Style: output.SuccessStyle, Count: len(succeeded), Label: "succeeded"},
		output.SummaryLine{Symbol: "✗", Style: output.ErrorStyle, Count: len(failed), Label: "failed"},
	))

	if len(failed) > 0 {
		return output.Silent(fmt.Errorf("%d succeeded, %d failed", len(succeeded), len(failed)))
	}

	return nil
}

func init() {
	addGroupFilterFlag(runCmd)
	amrCmd.AddCommand(runCmd)
}
