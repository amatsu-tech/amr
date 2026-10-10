package cli

import (
	"os"
	"path/filepath"
)

func initDirectory(args []string) (string, error) {
	if len(args) == 0 {
		return os.Getwd()
	}

	return filepath.Abs(args[0])
}
