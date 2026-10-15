package cli

import (
	"fmt"

	"github.com/amatsu-tech/amr/internal/cli/output"
	"github.com/spf13/cobra"
)

// Значение флага --jobs.
var jobs int

var amrCmd = &cobra.Command{
	Use:           "amr",
	Short:         "Amatsu multi-repo: manage many Git repositories as one workspace",
	Long:          "amr (Amatsu multi-repo) is a command-line tool for managing multiple Git repositories from a single interface.",
	Version:       fmt.Sprintf("%s (commit: %s)", buildVersion(), buildCommit()),
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if jobs < 1 {
			return fmt.Errorf("--jobs must be at least 1, got %d", jobs)
		}

		return nil
	},
}

func init() {
	amrCmd.SetVersionTemplate("{{.Name}} {{.Version}}\n")
	amrCmd.SetOut(output.Stdout)
	amrCmd.SetErr(output.Stderr)
	amrCmd.PersistentFlags().IntVarP(&jobs, "jobs", "j", 8, "number of repositories to process at the same time")
}

func Execute() error {
	if err := amrCmd.Execute(); err != nil {
		output.PrintError(err)
		return err
	}

	return nil
}
