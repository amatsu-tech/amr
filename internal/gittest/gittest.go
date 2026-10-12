package gittest

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func Isolate(t testing.TB) {
	t.Helper()

	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_AUTHOR_NAME", "amr test")
	t.Setenv("GIT_AUTHOR_EMAIL", "test@example.com")
	t.Setenv("GIT_COMMITTER_NAME", "amr test")
	t.Setenv("GIT_COMMITTER_EMAIL", "test@example.com")
}

// Run выполняет git в каталоге dir и возвращает вывод без пробелов по краям.
func Run(t testing.TB, dir string, args ...string) string {
	t.Helper()

	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
	}

	return strings.TrimSpace(string(output))
}

// CommitFile создаёт или перезаписывает файл и коммитит его.
func CommitFile(t testing.TB, dir, name, content string) {
	t.Helper()

	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	Run(t, dir, "add", name)
	Run(t, dir, "commit", "-q", "-m", "update "+name)
}

// Remote — bare-репозиторий во временном каталоге
type Remote struct {
	URL  string
	seed string
}

// NewRemote создаёт "сервер" с одним коммитом в ветке main.
func NewRemote(t testing.TB) Remote {
	t.Helper()

	base := t.TempDir()
	url := filepath.Join(base, "remote.git")
	seed := filepath.Join(base, "seed")

	Run(t, base, "init", "-q", "--bare", "-b", "main", url)
	Run(t, base, "clone", "-q", url, seed)
	CommitFile(t, seed, "README.md", "initial\n")
	Run(t, seed, "push", "-q", "origin", "main")

	return Remote{URL: url, seed: seed}
}

// Push добавляет на "сервер" новый коммит с файлом name.
func (r Remote) Push(t testing.TB, name string) {
	t.Helper()

	CommitFile(t, r.seed, name, name+"\n")
	Run(t, r.seed, "push", "-q", "origin", "main")
}

// Clone клонирует url в dir.
func Clone(t testing.TB, url, dir string) {
	t.Helper()

	Run(t, filepath.Dir(dir), "clone", "-q", url, dir)
}

// InitRepo создаёт в dir локальный репозиторий без origin с одним коммитом.
func InitRepo(t testing.TB, dir string) {
	t.Helper()

	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	Run(t, dir, "init", "-q", "-b", "main")
	CommitFile(t, dir, "README.md", "local\n")
}
