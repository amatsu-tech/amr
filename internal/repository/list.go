package repository

import (
	"path/filepath"

	"github.com/amatsu-tech/amr/internal/workspace"
)

func List(root string) ([]Repository, error) {
	config, err := workspace.ReadConfig(root)
	if err != nil {
		return nil, err
	}

	repositories := make([]Repository, len(config.Repositories))

	for i, entry := range config.Repositories {
		repositories[i] = Repository{
			Name:      entry.Name,
			Directory: filepath.Join(root, entry.Name),
		}
	}

	return repositories, nil
}
