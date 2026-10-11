package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

// makeRepo создаёт каталог, который git.IsRepository считает репозиторием.
func makeRepo(tb testing.TB, root, path string) {
	tb.Helper()

	if err := os.MkdirAll(filepath.Join(root, filepath.FromSlash(path), ".git"), 0o755); err != nil {
		tb.Fatal(err)
	}
}

func makeDir(tb testing.TB, root, path string) {
	tb.Helper()

	if err := os.MkdirAll(filepath.Join(root, filepath.FromSlash(path)), 0o755); err != nil {
		tb.Fatal(err)
	}
}

func scannedPaths(tb testing.TB, root string) []string {
	tb.Helper()

	repositories, err := Scan(root)
	if err != nil {
		tb.Fatal(err)
	}

	paths := make([]string, len(repositories))
	for i, repo := range repositories {
		paths[i] = repo.Path
	}

	return paths
}

func TestScan(t *testing.T) {
	root := t.TempDir()

	makeRepo(t, root, "top-level-repo")
	makeRepo(t, root, "group/nested-repo")
	makeRepo(t, root, "group/another-nested-repo")
	makeRepo(t, root, "top-level-repo/inner-repo") // внутри репозитория: не ищем
	makeRepo(t, root, ".hidden/repo")              // скрытый каталог: пропускаем
	makeDir(t, root, "plain-dir/subdir")           // не репозиторий

	got := scannedPaths(t, root)
	want := []string{"group/another-nested-repo", "group/nested-repo", "top-level-repo"}

	if !slices.Equal(got, want) {
		t.Errorf("Scan() = %v, want %v", got, want)
	}
}

func TestScanEmpty(t *testing.T) {
	if got := scannedPaths(t, t.TempDir()); len(got) != 0 {
		t.Errorf("Scan() = %v, want no repositories", got)
	}
}

// makeWorkspace строит дерево из groups групп по perGroup репозиториев,
// внутри каждого репозитория — dirsPerRepo обычных каталогов.
func makeWorkspace(b *testing.B, groups, perGroup, dirsPerRepo int) string {
	root := b.TempDir()

	for g := range groups {
		for r := range perGroup {
			repo := fmt.Sprintf("group%d/service%d", g, r)
			makeRepo(b, root, repo)

			for d := range dirsPerRepo {
				makeDir(b, root, fmt.Sprintf("%s/pkg%d", repo, d))
			}
		}
	}

	return root
}

func BenchmarkScan(b *testing.B) {
	for _, repos := range []int{10, 100} {
		root := makeWorkspace(b, 5, repos/5, 20)

		b.Run(fmt.Sprintf("repos=%d", repos), func(b *testing.B) {
			for b.Loop() {
				if _, err := Scan(root); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
