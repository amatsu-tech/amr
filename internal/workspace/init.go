package workspace

import (
	"errors"
	"fmt"
	"os"
)

func Init(root string) error {
	workspaceDir := workspaceDirectory(root)

	_, err := os.Stat(workspaceDir)

	if err == nil {
		return fmt.Errorf("%w in %s", ErrAlreadyInitialized, root)
	}

	if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := os.Mkdir(workspaceDir, 0o755); err != nil {
		return err
	}

	if err := WriteConfig(root, DefaultConfig()); err != nil {
		return err
	}

	return nil
}
