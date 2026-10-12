package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// newWorkspace создаёт инициализированный workspace во временном каталоге.
func newWorkspace(t *testing.T) string {
	t.Helper()

	root := t.TempDir()

	if err := Init(root); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	return root
}

func mkdirAll(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
}

func TestLocate(t *testing.T) {
	root := newWorkspace(t)
	deep := filepath.Join(root, "group", "repo", "internal", "pkg")
	mkdirAll(t, deep)

	tests := []struct {
		name  string
		start string
	}{
		{name: "from root", start: root},
		{name: "from deep subdirectory", start: deep},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Locate(tt.start)
			if err != nil {
				t.Fatalf("Locate() error = %v", err)
			}

			if got != root {
				t.Errorf("Locate() = %q, want %q", got, root)
			}
		})
	}
}

func TestLocateNearestWorkspace(t *testing.T) {
	outer := newWorkspace(t)
	inner := filepath.Join(outer, "nested-workspace")
	mkdirAll(t, inner)

	if err := Init(inner); err != nil {
		t.Fatal(err)
	}

	start := filepath.Join(inner, "repo")
	mkdirAll(t, start)

	got, err := Locate(start)
	if err != nil {
		t.Fatalf("Locate() error = %v", err)
	}

	if got != inner {
		t.Errorf("Locate() = %q, want nearest workspace %q", got, inner)
	}
}

func TestLocateRelativePath(t *testing.T) {
	root := newWorkspace(t)
	mkdirAll(t, filepath.Join(root, "repo"))
	t.Chdir(root)

	got, err := Locate("repo")
	if err != nil {
		t.Fatalf("Locate() error = %v", err)
	}

	// t.TempDir на macOS может лежать за символической ссылкой (/var → /private/var),
	// поэтому сравниваем реальные пути.
	if !samePath(t, got, root) {
		t.Errorf("Locate() = %q, want %q", got, root)
	}
}

func TestLocateNotFound(t *testing.T) {
	start := t.TempDir()

	_, err := Locate(start)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Locate() error = %v, want %v", err, ErrNotFound)
	}

	if !strings.Contains(err.Error(), start) {
		t.Errorf("error %q does not mention the start directory %q", err, start)
	}
}

func TestRoot(t *testing.T) {
	root := newWorkspace(t)
	sub := filepath.Join(root, "repo")
	mkdirAll(t, sub)
	t.Chdir(sub)

	got, err := Root()
	if err != nil {
		t.Fatalf("Root() error = %v", err)
	}

	if !samePath(t, got, root) {
		t.Errorf("Root() = %q, want %q", got, root)
	}
}

func samePath(t *testing.T, a, b string) bool {
	t.Helper()

	realA, err := filepath.EvalSymlinks(a)
	if err != nil {
		t.Fatal(err)
	}

	realB, err := filepath.EvalSymlinks(b)
	if err != nil {
		t.Fatal(err)
	}

	return realA == realB
}
