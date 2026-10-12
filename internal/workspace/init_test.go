package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestInit(t *testing.T) {
	root := t.TempDir()

	if err := Init(root); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	info, err := os.Stat(filepath.Join(root, DirectoryName))
	if err != nil {
		t.Fatalf("workspace directory not created: %v", err)
	}

	if !info.IsDir() {
		t.Fatalf("%s is not a directory", DirectoryName)
	}

	config, err := ReadConfig(root)
	if err != nil {
		t.Fatalf("ReadConfig() after Init error = %v", err)
	}

	if config.Version != WorkspaceConfigVersion {
		t.Errorf("Version = %d, want %d", config.Version, WorkspaceConfigVersion)
	}

	if len(config.Repositories) != 0 {
		t.Errorf("Repositories = %v, want empty", config.Repositories)
	}
}

func TestInitTwice(t *testing.T) {
	root := t.TempDir()

	if err := Init(root); err != nil {
		t.Fatalf("first Init() error = %v", err)
	}

	err := Init(root)
	if !errors.Is(err, ErrAlreadyInitialized) {
		t.Fatalf("second Init() error = %v, want %v", err, ErrAlreadyInitialized)
	}
}

func TestInitMissingDirectory(t *testing.T) {
	root := filepath.Join(t.TempDir(), "does-not-exist")

	if err := Init(root); err == nil {
		t.Fatal("Init() in a missing directory succeeded, want error")
	}
}
