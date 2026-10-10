package repository

import (
	"os"
	"path/filepath"

	"github.com/amatsu-tech/amr/internal/git"
)

func Scan(root string) ([]Repository, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	repositories := make([]Repository, 0)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		// Скрытые каталоги (.workspace, .git, ...)
		if len(name) > 0 && name[0] == '.' {
			continue
		}

		directory := filepath.Join(root, name)

		isRepository, err := git.IsRepository(directory)
		if err != nil {
			return nil, err
		}

		if !isRepository {
			continue
		}

		repositories = append(repositories, Repository{
			Name:      name,
			Directory: directory,
		})
	}

	return repositories, nil
}
