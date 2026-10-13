package repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/amatsu-tech/amr/internal/gittest"
)

func TestSync(t *testing.T) {
	gittest.Isolate(t)

	remote := gittest.NewRemote(t)
	base := t.TempDir()

	present := filepath.Join(base, "present")
	gittest.Clone(t, remote.URL, present)

	behind := filepath.Join(base, "behind")
	gittest.Clone(t, remote.URL, behind)

	remote.Push(t, "new-file.txt")

	// present подтягиваем заранее, чтобы он был уже актуальным.
	gittest.Run(t, present, "pull", "-q")

	tests := []struct {
		name string
		repo Repository
		want SyncOutcome
	}{
		{name: "missing repository is cloned", repo: Repository{URL: remote.URL, Directory: filepath.Join(base, "group", "missing")}, want: SyncCloned},
		{name: "repository behind is updated", repo: Repository{URL: remote.URL, Directory: behind}, want: SyncUpdated},
		{name: "current repository is up to date", repo: Repository{URL: remote.URL, Directory: present}, want: SyncUpToDate},
		{name: "missing repository without url is skipped", repo: Repository{Directory: filepath.Join(base, "no-url")}, want: SyncSkipped},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sync(tt.repo); got.Outcome != tt.want {
				t.Errorf("Sync() = %+v, want outcome %v", got, tt.want)
			}
		})
	}
}

func TestSyncFailure(t *testing.T) {
	gittest.Isolate(t)

	repo := Repository{URL: filepath.Join(t.TempDir(), "missing.git"), Directory: filepath.Join(t.TempDir(), "repo")}

	if got := Sync(repo); got.Outcome != SyncFailed || got.Reason == "" {
		t.Errorf("Sync() = %+v, want failure with a reason", got)
	}
}

func TestUpdateWorkspace(t *testing.T) {
	gittest.Isolate(t)

	remote := gittest.NewRemote(t)

	behind := filepath.Join(t.TempDir(), "workspace")
	gittest.Clone(t, remote.URL, behind)

	dirty := filepath.Join(t.TempDir(), "workspace")
	gittest.Clone(t, remote.URL, dirty)

	remote.Push(t, "amr.yaml")

	if err := os.WriteFile(filepath.Join(dirty, "local.txt"), []byte("wip\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	withoutUpstream := filepath.Join(t.TempDir(), "workspace")
	gittest.InitRepo(t, withoutUpstream)

	tests := []struct {
		name string
		root string
		want WorkspaceOutcome
	}{
		{name: "plain folder", root: t.TempDir(), want: WorkspaceNotRepository},
		{name: "behind upstream", root: behind, want: WorkspaceUpdated},
		{name: "already updated", root: behind, want: WorkspaceUpToDate},
		{name: "local changes", root: dirty, want: WorkspaceSkipped},
		{name: "no upstream", root: withoutUpstream, want: WorkspaceSkipped},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UpdateWorkspace(tt.root); got.Outcome != tt.want {
				t.Errorf("UpdateWorkspace() = %+v, want outcome %v", got, tt.want)
			}
		})
	}

	if _, err := os.Stat(filepath.Join(behind, "amr.yaml")); err != nil {
		t.Errorf("workspace update did not bring the new file: %v", err)
	}
}
