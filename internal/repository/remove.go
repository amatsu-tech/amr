package repository

import (
	"path/filepath"
	"slices"

	"github.com/amatsu-tech/amr/internal/workspace"
)

type RemoveResult struct {
	Name  string
	Error error
}

func Remove(root string, names ...string) ([]RemoveResult, error) {
	config, err := workspace.ReadConfig(root)
	if err != nil {
		return nil, err
	}

	results := make([]RemoveResult, 0, len(names))

	for _, rawName := range names {
		name := filepath.Base(filepath.Clean(rawName))

		if !slices.ContainsFunc(config.Repositories, func(entry workspace.RepositoryEntry) bool {
			return entry.Name == name
		}) {
			results = append(results, RemoveResult{Name: name, Error: ErrRepositoryNotTracked})
			continue
		}

		config.Repositories = slices.DeleteFunc(config.Repositories, func(entry workspace.RepositoryEntry) bool {
			return entry.Name == name
		})

		results = append(results, RemoveResult{Name: name})
	}

	if err := workspace.WriteConfig(root, config); err != nil {
		return nil, err
	}

	return results, nil
}
