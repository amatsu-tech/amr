package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateGitignore(t *testing.T) {
	root := t.TempDir()

	created, err := CreateGitignore(root)
	if err != nil {
		t.Fatalf("CreateGitignore() error = %v", err)
	}

	if !created {
		t.Fatal("CreateGitignore() = false, want true")
	}

	data, err := os.ReadFile(filepath.Join(root, ".gitignore"))
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != Gitignore {
		t.Errorf(".gitignore = %q, want %q", data, Gitignore)
	}
}

func TestCreateGitignoreKeepsExisting(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".gitignore")

	if err := os.WriteFile(path, []byte("node_modules/\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	created, err := CreateGitignore(root)
	if err != nil {
		t.Fatal(err)
	}

	if created {
		t.Error("CreateGitignore() = true for an existing .gitignore, want false")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "node_modules/\n" {
		t.Errorf("existing .gitignore was changed to %q", data)
	}
}

func TestCreateGitignoreSkipsGitRepository(t *testing.T) {
	root := t.TempDir()
	mkdirAll(t, filepath.Join(root, ".git"))

	created, err := CreateGitignore(root)
	if err != nil {
		t.Fatal(err)
	}

	if created {
		t.Error("CreateGitignore() = true inside a git repository, want false")
	}

	if _, err := os.Stat(filepath.Join(root, ".gitignore")); err == nil {
		t.Error(".gitignore was created inside a git repository")
	}
}
