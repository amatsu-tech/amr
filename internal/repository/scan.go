package repository

import (
	"io/fs"
	"path/filepath"

	"github.com/amatsu-tech/amr/internal/git"
)

// Scan рекурсивно ищет git-репозитории внутри workspace. Внутрь найденного
// репозитория не заходит, скрытые каталоги (.amr, .git, ...) пропускает.
func Scan(root string) ([]Repository, error) {
	repositories := make([]Repository, 0)

	err := filepath.WalkDir(root, func(current string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !entry.IsDir() || current == root {
			return nil
		}

		name := entry.Name()

		if len(name) > 0 && name[0] == '.' {
			return filepath.SkipDir
		}

		isRepository, err := git.IsRepository(current)
		if err != nil {
			return err
		}

		if !isRepository {
			return nil
		}

		relative, err := filepath.Rel(root, current)
		if err != nil {
			return err
		}

		repositories = append(repositories, Repository{
			Path:      filepath.ToSlash(relative),
			Directory: current,
		})

		return filepath.SkipDir
	})
	if err != nil {
		return nil, err
	}

	return repositories, nil
}
