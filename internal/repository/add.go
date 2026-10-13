package repository

import (
	"errors"
	"os"
	"slices"

	"github.com/amatsu-tech/amr/internal/git"
	"github.com/amatsu-tech/amr/internal/workspace"
)

type AddResult struct {
	Path  string
	Error error
}

// Add добавляет репозитории в конфиг; groups записываются каждому добавленному.
func Add(root string, groups []string, paths ...string) ([]AddResult, error) {
	config, err := workspace.ReadConfig(root)
	if err != nil {
		return nil, err
	}

	results := make([]AddResult, 0, len(paths))

	for _, rawPath := range paths {
		path, err := resolvePath(root, rawPath)
		if err != nil {
			results = append(results, AddResult{Path: rawPath, Error: err})
			continue
		}

		repoDirectory := directory(root, path)

		if _, err := os.Stat(repoDirectory); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				results = append(results, AddResult{Path: path, Error: ErrRepositoryNotFound})
				continue
			}

			return nil, err
		}

		isRepository, err := git.IsRepository(repoDirectory)
		if err != nil {
			return nil, err
		}

		if !isRepository {
			results = append(results, AddResult{Path: path, Error: ErrNotGitRepository})
			continue
		}

		if slices.ContainsFunc(config.Repositories, func(entry workspace.RepositoryEntry) bool {
			return entry.Path == path
		}) {
			results = append(results, AddResult{Path: path, Error: ErrRepositoryExists})
			continue
		}

		url, err := git.RemoteURL(repoDirectory)
		if err != nil {
			results = append(results, AddResult{Path: path, Error: err})
			continue
		}

		config.Repositories = append(config.Repositories, workspace.RepositoryEntry{
			Path:   path,
			URL:    url,
			Groups: groups,
		})
		results = append(results, AddResult{Path: path})
	}

	if err := workspace.WriteConfig(root, config); err != nil {
		return nil, err
	}

	return results, nil
}
