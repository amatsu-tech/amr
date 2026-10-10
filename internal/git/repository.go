package git

import (
	"errors"
	"os"
	"path/filepath"
)

func IsRepository(directory string) (bool, error) {
	info, err := os.Stat(filepath.Join(directory, ".git"))
	if err == nil {
		return info.IsDir(), nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err

}
