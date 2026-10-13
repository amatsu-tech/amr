package repository

import (
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
			Path:      entry.Path,
			URL:       entry.URL,
			Groups:    entry.Groups,
			Directory: directory(root, entry.Path),
		}
	}

	return repositories, nil
}
