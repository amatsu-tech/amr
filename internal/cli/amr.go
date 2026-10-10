package cli

import (
	"fmt"

	"github.com/amatsu-tech/amr/internal/cli/output"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
)

var amrCmd = &cobra.Command{
	Use:           "amr",
	Short:         "Amatsu multi-repo: manage many Git repositories as one workspace",
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
