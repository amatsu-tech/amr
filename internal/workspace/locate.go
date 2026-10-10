package workspace

import (
	"os"
	"path/filepath"
)

func Locate(start string) (string, error) {
	current, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(workspaceDirectory(current)); err == nil {
			return current, nil
		}

		parent := filepath.Dir(current)

		if parent == current {
			return "", ErrNotFound
		}

		current = parent
	}
}
