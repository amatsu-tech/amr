package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/amatsu-tech/amr/internal/gittest"
)

func TestIsRepository(t *testing.T) {
	gittest.Isolate(t)

	base := t.TempDir()

	repo := filepath.Join(base, "repo")
	gittest.InitRepo(t, repo)

	plainDir := filepath.Join(base, "plain-dir")
	if err := os.Mkdir(plainDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// .git в виде файла бывает у worktree и сабмодулей; amr такие не считает репозиториями.
	gitFile := filepath.Join(base, "git-file")
	if err := os.Mkdir(gitFile, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(gitFile, ".git"), []byte("gitdir: ../somewhere\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		dir  string
		want bool
	}{
		{name: "git repository", dir: repo, want: true},
		{name: "plain directory", dir: plainDir, want: false},
		{name: "missing directory", dir: filepath.Join(base, "missing"), want: false},
		{name: ".git is a file", dir: gitFile, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsRepository(tt.dir)
			if err != nil {
				t.Fatalf("IsRepository() error = %v", err)
			}

			if got != tt.want {
				t.Errorf("IsRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTopLevel(t *testing.T) {
	gittest.Isolate(t)

	repo := filepath.Join(t.TempDir(), "repo")
	gittest.InitRepo(t, repo)

	subfolder := filepath.Join(repo, "internal", "pkg")
	if err := os.MkdirAll(subfolder, 0o755); err != nil {
		t.Fatal(err)
	}

	want, err := filepath.EvalSymlinks(repo)
	if err != nil {
		t.Fatal(err)
	}

	for _, dir := range []string{repo, subfolder} {
		got, err := TopLevel(dir)
		if err != nil {
			t.Fatalf("TopLevel(%s) error = %v", dir, err)
		}

		if got != want {
			t.Errorf("TopLevel(%s) = %q, want %q", dir, got, want)
		}
	}

	if _, err := TopLevel(t.TempDir()); err == nil {
		t.Error("TopLevel() outside a repository succeeded, want error")
	}
}
