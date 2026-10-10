package repository

import (
	"errors"
	"os"
	"path/filepath"
	"slices"

	"github.com/amatsu-tech/amr/internal/git"
	"github.com/amatsu-tech/amr/internal/workspace"
)

type AddResult struct {
	Name  string
	Error error
}

func Add(root string, names ...string) ([]AddResult, error) {
	config, err := workspace.ReadConfig(root)
	if err != nil {
		return nil, err
	}

	results := make([]AddResult, 0, len(names))

	for _, rawName := range names {
		name := filepath.Base(filepath.Clean(rawName))
		directory := filepath.Join(root, name)

		if _, err := os.Stat(directory); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				results = append(results, AddResult{Name: name, Error: ErrRepositoryNotFound})
				continue
			}

			return nil, err
		}

		isRepository, err := git.IsRepository(directory)
		if err != nil {
			return nil, err
		}

		if !isRepository {
			results = append(results, AddResult{Name: name, Error: ErrNotGitRepository})
			continue
		}

		if slices.ContainsFunc(config.Repositories, func(entry workspace.RepositoryEntry) bool {
			return entry.Name == name
		}) {
			results = append(results, AddResult{Name: name, Error: ErrRepositoryExists})
			continue
		}

		config.Repositories = append(config.Repositories, workspace.RepositoryEntry{
			Name:   name,
			Source: workspace.SourceLocal,
		})
		results = append(results, AddResult{Name: name})
	}

	if err := workspace.WriteConfig(root, config); err != nil {
		return nil, err
	}

	return results, nil
}
