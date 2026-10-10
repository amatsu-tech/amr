package workspace

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type RepositoryEntry struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url,omitempty"`
}

type Config struct {
	Version      int               `yaml:"version"`
	Repositories []RepositoryEntry `yaml:"repositories"`
}

func DefaultConfig() Config {
	return Config{
		Version:      WorkspaceConfigVersion,
		Repositories: []RepositoryEntry{},
	}
}

func WriteConfig(root string, config Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(
		configPath(root),
		data,
		0o644,
	)
}

func ReadConfig(root string) (Config, error) {
	data, err := os.ReadFile(configPath(root))
	if err != nil {
		return Config{}, err
	}

	var probe struct {
		Version int `yaml:"version"`
	}

	if err := yaml.Unmarshal(data, &probe); err != nil {
		return Config{}, err
	}

	if probe.Version > WorkspaceConfigVersion {
		return Config{}, fmt.Errorf("%w: %d (this amr CLI supports up to version %d)",
			ErrUnsupportedConfigVersion, probe.Version, WorkspaceConfigVersion)
	}

	if probe.Version < WorkspaceConfigVersion {
		config, err := migrateConfig(data, probe.Version)
		if err != nil {
			return Config{}, err
		}

		if err := WriteConfig(root, config); err != nil {
			return Config{}, err
		}

		return config, nil
	}

	var config Config

	if err := yaml.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}
