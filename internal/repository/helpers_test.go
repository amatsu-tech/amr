package repository

import (
	"testing"

	"github.com/amatsu-tech/amr/internal/workspace"
)

func newWorkspace(t *testing.T) string {
	t.Helper()

	root := t.TempDir()

	if err := workspace.Init(root); err != nil {
		t.Fatal(err)
	}

	return root
}

func readEntries(t *testing.T, root string) []workspace.RepositoryEntry {
	t.Helper()

	config, err := workspace.ReadConfig(root)
	if err != nil {
		t.Fatal(err)
	}

	return config.Repositories
}
