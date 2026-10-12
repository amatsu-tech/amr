package repository

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/amatsu-tech/amr/internal/workspace"
)

func TestAdd(t *testing.T) {
	isolateGit(t)

	root := newWorkspace(t)
	remote := newRemote(t)

	cloneInto(t, remote.URL, filepath.Join(root, "top-level-repo"))
	initRepo(t, filepath.Join(root, "group", "repo-without-origin"))
	makeDir(t, root, "plain-dir")

	results, err := Add(root,
		filepath.Join(root, "top-level-repo"),
		filepath.Join(root, "group", "repo-without-origin"),
		filepath.Join(root, "plain-dir"),
		filepath.Join(root, "missing-dir"),
		filepath.Join(filepath.Dir(root), "outside-repo"),
	)
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	wantErrors := []error{nil, nil, ErrNotGitRepository, ErrRepositoryNotFound, ErrOutsideWorkspace}

	if len(results) != len(wantErrors) {
		t.Fatalf("Add() returned %d results, want %d", len(results), len(wantErrors))
	}

	for i, result := range results {
		if !errors.Is(result.Error, wantErrors[i]) {
			t.Errorf("result %d (%s): error = %v, want %v", i, result.Path, result.Error, wantErrors[i])
		}
	}

	want := []workspace.RepositoryEntry{
		{Path: "top-level-repo", URL: remote.URL},
		{Path: "group/repo-without-origin"},
	}

	if got := readEntries(t, root); !slices.Equal(got, want) {
		t.Errorf("config repositories = %+v, want %+v", got, want)
	}
}

func TestAddDuplicate(t *testing.T) {
	isolateGit(t)

	root := newWorkspace(t)
	repo := filepath.Join(root, "repo")
	initRepo(t, repo)

	if _, err := Add(root, repo); err != nil {
		t.Fatal(err)
	}

	results, err := Add(root, repo)
	if err != nil {
		t.Fatal(err)
	}

	if !errors.Is(results[0].Error, ErrRepositoryExists) {
		t.Errorf("second Add() error = %v, want %v", results[0].Error, ErrRepositoryExists)
	}

	if got := readEntries(t, root); len(got) != 1 {
		t.Errorf("config has %d repositories after duplicate add, want 1", len(got))
	}
}

func TestAddRelativeToCurrentDirectory(t *testing.T) {
	isolateGit(t)

	root := newWorkspace(t)
	initRepo(t, filepath.Join(root, "group", "nested-repo"))
	t.Chdir(filepath.Join(root, "group"))

	results, err := Add(root, "nested-repo")
	if err != nil {
		t.Fatal(err)
	}

	if results[0].Error != nil {
		t.Fatalf("Add() error = %v", results[0].Error)
	}

	if results[0].Path != "group/nested-repo" {
		t.Errorf("Add() path = %q, want %q", results[0].Path, "group/nested-repo")
	}
}

func TestRemove(t *testing.T) {
	isolateGit(t)

	root := newWorkspace(t)
	first := filepath.Join(root, "first-repo")
	second := filepath.Join(root, "group", "second-repo")
	initRepo(t, first)
	initRepo(t, second)

	if _, err := Add(root, first, second); err != nil {
		t.Fatal(err)
	}

	results, err := Remove(root,
		first,
		filepath.Join(root, "never-added"),
		filepath.Join(filepath.Dir(root), "outside-repo"),
	)
	if err != nil {
		t.Fatalf("Remove() error = %v", err)
	}

	wantErrors := []error{nil, ErrRepositoryNotTracked, ErrOutsideWorkspace}

	for i, result := range results {
		if !errors.Is(result.Error, wantErrors[i]) {
			t.Errorf("result %d (%s): error = %v, want %v", i, result.Path, result.Error, wantErrors[i])
		}
	}

	want := []workspace.RepositoryEntry{{Path: "group/second-repo"}}

	if got := readEntries(t, root); !slices.Equal(got, want) {
		t.Errorf("config repositories = %+v, want %+v", got, want)
	}

	// remove убирает репозиторий только из конфига, файлы остаются на диске.
	if _, err := os.Stat(first); err != nil {
		t.Errorf("removed repository directory is gone: %v", err)
	}
}

func TestList(t *testing.T) {
	root := newWorkspace(t)

	config := workspace.DefaultConfig()
	config.Repositories = []workspace.RepositoryEntry{
		{Path: "top-level-repo", URL: "git@example.com:org/top-level-repo.git"},
		{Path: "group/nested-repo"},
	}

	if err := workspace.WriteConfig(root, config); err != nil {
		t.Fatal(err)
	}

	got, err := List(root)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	want := []Repository{
		{Path: "top-level-repo", URL: "git@example.com:org/top-level-repo.git", Directory: filepath.Join(root, "top-level-repo")},
		{Path: "group/nested-repo", Directory: filepath.Join(root, "group", "nested-repo")},
	}

	if !slices.Equal(got, want) {
		t.Errorf("List() = %+v, want %+v", got, want)
	}
}
