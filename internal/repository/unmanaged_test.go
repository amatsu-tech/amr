package repository

import (
	"path/filepath"
	"slices"
	"testing"

	"github.com/amatsu-tech/amr/internal/gittest"
)

func TestUnmanaged(t *testing.T) {
	gittest.Isolate(t)

	root := newWorkspace(t)
	gittest.InitRepo(t, filepath.Join(root, "managed-repo"))
	gittest.InitRepo(t, filepath.Join(root, "group", "new-repo"))
	gittest.InitRepo(t, filepath.Join(root, "another-new-repo"))

	if _, err := Add(root, filepath.Join(root, "managed-repo")); err != nil {
		t.Fatal(err)
	}

	repositories, err := Unmanaged(root)
	if err != nil {
		t.Fatalf("Unmanaged() error = %v", err)
	}

	got := make([]string, len(repositories))
	for i, repo := range repositories {
		got[i] = repo.Path
	}

	if want := []string{"another-new-repo", "group/new-repo"}; !slices.Equal(got, want) {
		t.Errorf("Unmanaged() = %v, want %v", got, want)
	}
}

func TestUnmanagedAllAdded(t *testing.T) {
	gittest.Isolate(t)

	root := newWorkspace(t)
	repo := filepath.Join(root, "repo")
	gittest.InitRepo(t, repo)

	if _, err := Add(root, repo); err != nil {
		t.Fatal(err)
	}

	repositories, err := Unmanaged(root)
	if err != nil {
		t.Fatal(err)
	}

	if len(repositories) != 0 {
		t.Errorf("Unmanaged() = %v, want none", repositories)
	}
}
