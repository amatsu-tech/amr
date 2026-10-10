package workspace

import (
	"fmt"
	"os"
	"path/filepath"
)

func Locate(start string) (string, error) {
	origin, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}

	current := origin

	for {
		if _, err := os.Stat(workspaceDirectory(current)); err == nil {
			return current, nil
		}

		parent := filepath.Dir(current)

		if parent == current {
			return "", fmt.Errorf("%w: no %s directory in %s or any parent directory", ErrNotFound, DirectoryName, origin)
		}

		current = parent
	}
}
