package repository

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/amatsu-tech/amr/internal/gittest"
)

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()

	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, dir, name string) string {
	t.Helper()

	data, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		t.Fatal(err)
	}

	return string(data)
}

func TestDefaultStashName(t *testing.T) {
	now := time.Date(2026, 10, 15, 14, 30, 5, 0, time.UTC)

	if got := DefaultStashName(now); got != "amr-2026-10-15-143005" {
		t.Errorf("DefaultStashName() = %q", got)
	}
}

func TestStashSave(t *testing.T) {
	gittest.Isolate(t)

	clean := Repository{Directory: filepath.Join(t.TempDir(), "clean")}
	gittest.InitRepo(t, clean.Directory)

	dirty := Repository{Directory: filepath.Join(t.TempDir(), "dirty")}
	gittest.InitRepo(t, dirty.Directory)
	writeFile(t, dirty.Directory, "new-file.txt", "new\n")

	missing := Repository{Directory: filepath.Join(t.TempDir(), "missing")}

	tests := []struct {
		name string
		repo Repository
		want StashOutcome
	}{
		{name: "no changes", repo: clean, want: StashNoChanges},
		{name: "only a new file", repo: dirty, want: StashSaved},
		{name: "missing directory", repo: missing, want: StashFailed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StashSave(tt.repo, "task"); got.Outcome != tt.want {
				t.Errorf("StashSave() = %+v, want outcome %v", got, tt.want)
			}
		})
	}
}

func TestStashByName(t *testing.T) {
	gittest.Isolate(t)

	repo := Repository{Directory: filepath.Join(t.TempDir(), "repo")}
	gittest.InitRepo(t, repo.Directory)

	writeFile(t, repo.Directory, "README.md", "login work\n")
	StashSave(repo, "login-fix")

	// Обычный git stash между нашими: amr не должен его трогать.
	writeFile(t, repo.Directory, "README.md", "manual work\n")
	gittest.Run(t, repo.Directory, "stash", "push", "-m", "manual")

	writeFile(t, repo.Directory, "README.md", "search work\n")
	StashSave(repo, "search")

	names, err := StashNames(repo)
	if err != nil {
		t.Fatal(err)
	}

	if want := []string{"search", "login-fix"}; !slices.Equal(names, want) {
		t.Fatalf("StashNames() = %v, want %v", names, want)
	}

	// Возвращаем более старый стеш по имени, а не последний.
	if got := StashPop(repo, "login-fix"); got.Outcome != StashDone {
		t.Fatalf("StashPop(login-fix) = %+v", got)
	}

	if got := readFile(t, repo.Directory, "README.md"); got != "login work\n" {
		t.Errorf("README.md after pop = %q, want the login-fix version", got)
	}

	if got := StashPop(repo, "login-fix"); got.Outcome != StashNotFound {
		t.Errorf("second StashPop(login-fix) = %+v, want not found", got)
	}

	if got := StashDrop(repo, "search"); got.Outcome != StashDone {
		t.Errorf("StashDrop(search) = %+v", got)
	}

	if names, _ := StashNames(repo); len(names) != 0 {
		t.Errorf("StashNames() after pop and drop = %v, want none", names)
	}

	// Обычный стеш на месте.
	if got := gittest.Run(t, repo.Directory, "stash", "list"); got == "" {
		t.Error("manual git stash was removed")
	}
}

func TestStashPopLatest(t *testing.T) {
	gittest.Isolate(t)

	repo := Repository{Directory: filepath.Join(t.TempDir(), "repo")}
	gittest.InitRepo(t, repo.Directory)

	writeFile(t, repo.Directory, "README.md", "older\n")
	StashSave(repo, "older")
	writeFile(t, repo.Directory, "README.md", "newer\n")
	StashSave(repo, "newer")

	if got := StashPop(repo, ""); got.Outcome != StashDone {
		t.Fatalf("StashPop(\"\") = %+v", got)
	}

	if got := readFile(t, repo.Directory, "README.md"); got != "newer\n" {
		t.Errorf("README.md = %q, want the latest stash", got)
	}
}
