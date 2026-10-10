package cli

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/amatsu-tech/amr/internal/cli/output"
	"github.com/amatsu-tech/amr/internal/repository"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run -- <command>...",
	Short: "Run a command in every managed repository",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runCommand,
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

	repositories, err := repository.List(root)
	if err != nil {
		return err
	}

	if len(repositories) == 0 {
		cmd.Println("No repositories have been added to this workspace.")
		return nil
	}

	name, commandArgs := args[0], args[1:]

	var succeeded, failed []runResult

	for _, repo := range repositories {
		var buf bytes.Buffer

		runErr := repository.Run(repo, name, commandArgs, &buf, &buf)
		result := runResult{Name: repo.Name, Output: strings.TrimSpace(buf.String()), Error: runErr}

		if runErr != nil {
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
	amrCmd.AddCommand(runCmd)
}
