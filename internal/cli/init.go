package cli

import (
	"github.com/amatsu-tech/amr/internal/workspace"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Initialize a workspace",
	Long: `Initialize a workspace in the specified directory.

If no directory is provided, the current working directory is used.`,
	Example: `  amr init
  amr init ~/projects/my-workspace`,
	RunE: initCommand,
	Args: cobra.MaximumNArgs(1),
}

func initCommand(cmd *cobra.Command, args []string) error {
	root, err := initDirectory(args)
	if err != nil {
		return err
	}

	if err := workspace.Init(root); err != nil {
		return err
	}

	cmd.Printf("Initialized workspace in %s\n", root)

	return nil
}

func init() {
	amrCmd.AddCommand(initCmd)
}
