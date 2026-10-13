package repository

import (
	"slices"

	"github.com/amatsu-tech/amr/internal/workspace"
)

// Unmanaged возвращает найденные в workspace репозитории, которых ещё нет в конфиге.
func Unmanaged(root string) ([]Repository, error) {
	config, err := workspace.ReadConfig(root)
	if err != nil {
		return nil, err
	}

	discovered, err := Scan(root)
	if err != nil {
		return nil, err
	}

	return slices.DeleteFunc(discovered, func(repo Repository) bool {
		return slices.ContainsFunc(config.Repositories, func(entry workspace.RepositoryEntry) bool {
			return entry.Path == repo.Path
		})
	}), nil
}
