package repository

import (
	"slices"

	"github.com/amatsu-tech/amr/internal/workspace"
)

type RemoveResult struct {
	Path  string
	Error error
}

func Remove(root string, paths ...string) ([]RemoveResult, error) {
	config, err := workspace.ReadConfig(root)
	if err != nil {
		return nil, err
	}

	results := make([]RemoveResult, 0, len(paths))

	for _, rawPath := range paths {
		path, err := resolvePath(root, rawPath)
		if err != nil {
			results = append(results, RemoveResult{Path: rawPath, Error: err})
			continue
		}

		if !slices.ContainsFunc(config.Repositories, func(entry workspace.RepositoryEntry) bool {
			return entry.Path == path
		}) {
			results = append(results, RemoveResult{Path: path, Error: ErrRepositoryNotTracked})
			continue
		}

		config.Repositories = slices.DeleteFunc(config.Repositories, func(entry workspace.RepositoryEntry) bool {
			return entry.Path == path
		})

		results = append(results, RemoveResult{Path: path})
	}

	if err := workspace.WriteConfig(root, config); err != nil {
		return nil, err
	}

	return results, nil
}
