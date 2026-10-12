package repository

import (
	"path/filepath"
	"testing"
)

func TestClone(t *testing.T) {
	isolateGit(t)

	root := newWorkspace(t)
	remote := newRemote(t)

	existing := filepath.Join(root, "existing-repo")
	initRepo(t, existing)

	tests := []struct {
		name        string
		repo        Repository
		wantOutcome CloneOutcome
	}{
		{
			name:        "missing repository is cloned into nested path",
			repo:        Repository{Path: "group/cloned-repo", URL: remote.URL, Directory: filepath.Join(root, "group", "cloned-repo")},
			wantOutcome: CloneCloned,
		},
		{
			name:        "existing directory is left alone",
			repo:        Repository{Path: "existing-repo", URL: remote.URL, Directory: existing},
			wantOutcome: CloneExists,
		},
		{
			name:        "repository without url is skipped",
			repo:        Repository{Path: "no-url-repo", Directory: filepath.Join(root, "no-url-repo")},
			wantOutcome: CloneSkipped,
		},
		{
			name:        "unreachable url fails",
			repo:        Repository{Path: "broken-repo", URL: filepath.Join(t.TempDir(), "missing.git"), Directory: filepath.Join(root, "broken-repo")},
			wantOutcome: CloneFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Clone(tt.repo)

			if result.Outcome != tt.wantOutcome {
				t.Fatalf("Clone() outcome = %v (%s), want %v", result.Outcome, result.Reason, tt.wantOutcome)
			}

			if (result.Outcome == CloneSkipped || result.Outcome == CloneFailed) && result.Reason == "" {
				t.Error("Clone() returned no reason")
			}
		})
	}

	cloned := filepath.Join(root, "group", "cloned-repo")

	if got := runGit(t, cloned, "config", "--get", "remote.origin.url"); got != remote.URL {
		t.Errorf("cloned repository origin = %q, want %q", got, remote.URL)
	}
}
