package cli

import (
	"fmt"

	"github.com/amatsu-tech/amr/internal/cli/output"
	"github.com/spf13/cobra"
)

var amrCmd = &cobra.Command{
	Use:           "amr",
	Short:         "Amatsu multi-repo: manage many Git repositories as one workspace",
	Long:          "amr (Amatsu multi-repo) is a command-line tool for managing multiple Git repositories from a single interface.",
	Version:       fmt.Sprintf("%s (commit: %s)", version, commit),
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	amrCmd.SetVersionTemplate("{{.Name}} {{.Version}}\n")
}

func Execute() error {
	if err := amrCmd.Execute(); err != nil {
		output.PrintError(err)
		return err
	}

	return nil
}
