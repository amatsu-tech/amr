package repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/amatsu-tech/amr/internal/gittest"
)

// clonedRepo клонирует remote в workspace и возвращает Repository для него.
func clonedRepo(t *testing.T, remote gittest.Remote) Repository {
	t.Helper()

	dir := filepath.Join(t.TempDir(), "repo")
	gittest.Clone(t, remote.URL, dir)

	return Repository{Path: "repo", URL: remote.URL, Directory: dir}
}

func TestPullUpdated(t *testing.T) {
	gittest.Isolate(t)

	remote := gittest.NewRemote(t)
	repo := clonedRepo(t, remote)
	remote.Push(t, "new-file.txt")

	result := Pull(repo)
	if result.Outcome != PullUpdated {
		t.Fatalf("Pull() outcome = %v (%s), want %v", result.Outcome, result.Reason, PullUpdated)
	}

	if _, err := os.Stat(filepath.Join(repo.Directory, "new-file.txt")); err != nil {
		t.Errorf("pulled file is missing: %v", err)
	}
}

func TestPullUpToDate(t *testing.T) {
	gittest.Isolate(t)

	repo := clonedRepo(t, gittest.NewRemote(t))

	if result := Pull(repo); result.Outcome != PullUpToDate {
		t.Errorf("Pull() outcome = %v (%s), want %v", result.Outcome, result.Reason, PullUpToDate)
	}
}

func TestPullSkipsModifiedRepository(t *testing.T) {
	gittest.Isolate(t)

	remote := gittest.NewRemote(t)
	repo := clonedRepo(t, remote)
	remote.Push(t, "new-file.txt")

	if err := os.WriteFile(filepath.Join(repo.Directory, "local-change.txt"), []byte("wip\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	result := Pull(repo)
	if result.Outcome != PullSkipped || result.Reason != "modified files" {
		t.Fatalf("Pull() = %+v, want skipped with reason %q", result, "modified files")
	}

	// Изменения с сервера не должны были подтянуться.
	if _, err := os.Stat(filepath.Join(repo.Directory, "new-file.txt")); err == nil {
		t.Error("Pull() updated a repository with local changes")
	}
}

func TestPullFailures(t *testing.T) {
	gittest.Isolate(t)

	remote := gittest.NewRemote(t)

	diverged := clonedRepo(t, remote)
	gittest.CommitFile(t, diverged.Directory, "local-commit.txt", "local\n")
	remote.Push(t, "remote-commit.txt")

	withoutUpstream := filepath.Join(t.TempDir(), "local-only")
	gittest.InitRepo(t, withoutUpstream)

	tests := []struct {
		name       string
		repo       Repository
		wantReason string
	}{
		{
			name:       "no upstream",
			repo:       Repository{Path: "local-only", Directory: withoutUpstream},
			wantReason: "no upstream branch configured",
		},
		{
			name: "diverged history cannot be fast-forwarded",
			repo: diverged,
		},
		{
			name:       "missing directory",
			repo:       Repository{Path: "missing", Directory: filepath.Join(t.TempDir(), "missing")},
			wantReason: "repository directory not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Pull(tt.repo)

			if result.Outcome != PullFailed {
				t.Fatalf("Pull() outcome = %v (%s), want %v", result.Outcome, result.Reason, PullFailed)
			}

			if result.Reason == "" {
				t.Error("Pull() returned no reason")
			}

			if tt.wantReason != "" && result.Reason != tt.wantReason {
				t.Errorf("Pull() reason = %q, want %q", result.Reason, tt.wantReason)
			}
		})
	}
}
