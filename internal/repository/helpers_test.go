package repository

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/amatsu-tech/amr/internal/workspace"
)

// isolateGit запускает git в тесте без глобальной и системной конфигурации
// пользователя (подписи коммитов, insteadOf, хуки) и с временным автором.
func isolateGit(t *testing.T) {
	t.Helper()

	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_AUTHOR_NAME", "amr test")
	t.Setenv("GIT_AUTHOR_EMAIL", "test@example.com")
	t.Setenv("GIT_COMMITTER_NAME", "amr test")
	t.Setenv("GIT_COMMITTER_EMAIL", "test@example.com")
}

func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()

	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
	}

	return strings.TrimSpace(string(output))
}

// commitFile создаёт или перезаписывает файл и коммитит его.
func commitFile(t *testing.T, dir, name, content string) {
	t.Helper()

	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	runGit(t, dir, "add", name)
	runGit(t, dir, "commit", "-q", "-m", "update "+name)
}

// testRemote — bare-репозиторий, который играет роль сервера (GitHub),
// и рабочая копия, через которую в него можно пушить новые коммиты.
type testRemote struct {
	URL  string
	seed string
}

func newRemote(t *testing.T) testRemote {
	t.Helper()

	base := t.TempDir()
	url := filepath.Join(base, "remote.git")
	seed := filepath.Join(base, "seed")

	runGit(t, base, "init", "-q", "--bare", "-b", "main", url)
	runGit(t, base, "clone", "-q", url, seed)
	commitFile(t, seed, "README.md", "initial\n")
	runGit(t, seed, "push", "-q", "origin", "main")

	return testRemote{URL: url, seed: seed}
}

// push добавляет в удалённый репозиторий новый коммит.
func (r testRemote) push(t *testing.T, name string) {
	t.Helper()

	commitFile(t, r.seed, name, name+"\n")
	runGit(t, r.seed, "push", "-q", "origin", "main")
}

func cloneInto(t *testing.T, url, dir string) {
	t.Helper()

	runGit(t, filepath.Dir(dir), "clone", "-q", url, dir)
}

// initRepo создаёт локальный репозиторий без origin с одним коммитом.
func initRepo(t *testing.T, dir string) {
	t.Helper()

	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	runGit(t, dir, "init", "-q", "-b", "main")
	commitFile(t, dir, "README.md", "local\n")
}

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
