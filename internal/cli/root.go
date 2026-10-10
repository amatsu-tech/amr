package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "root",
	Short: "Print workspace root",
	RunE:  rootCommand,
}

func rootCommand(cmd *cobra.Command, args []string) error {
	root, err := currentWorkspace()
	if err != nil {
		return err
	}

	cmd.Println(root)

	return nil
}

func init() {
	amrCmd.AddCommand(rootCmd)
}
