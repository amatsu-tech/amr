package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/amatsu-tech/amr/internal/gittest"
)

func TestParseStatus(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   StatusSummary
	}{
		{
			name:   "clean with upstream",
			output: "## main...origin/main\n",
			want:   StatusSummary{Branch: "main", HasUpstream: true},
		},
		{
			name:   "no upstream",
			output: "## feature\n",
			want:   StatusSummary{Branch: "feature"},
		},
		{
			name:   "ahead and behind",
			output: "## main...origin/main [ahead 2, behind 3]\n",
			want:   StatusSummary{Branch: "main", HasUpstream: true, Ahead: 2, Behind: 3},
		},
		{
			name:   "detached head",
			output: "## HEAD (no branch)\n",
			want:   StatusSummary{Branch: "detached"},
		},
		{
			name:   "no commits yet",
			output: "## No commits yet on main\n",
			want:   StatusSummary{Branch: "main"},
		},
		{
			name: "changes of every kind",
			output: "## main...origin/main\n" +
				"M  staged.go\n" +
				" M modified.go\n" +
				"MM both.go\n" +
				"?? new.go\n" +
				"UU conflict.go\n" +
				"AA added-by-both.go\n",
			want: StatusSummary{
				Branch:      "main",
				HasUpstream: true,
				Staged:      2,
				Modified:    2,
				Untracked:   1,
				Conflicted:  2,
			},
		},
		{
			name:   "empty output",
			output: "",
			want:   StatusSummary{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseStatus(tt.output); got != tt.want {
				t.Errorf("parseStatus() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

// statusOutput строит вывод git status с заданным числом изменённых файлов.
func statusOutput(files int) string {
	var b strings.Builder

	b.WriteString("## main...origin/main [ahead 1, behind 2]\n")

	for i := range files {
		switch i % 3 {
		case 0:
			fmt.Fprintf(&b, " M internal/pkg/file%d.go\n", i)
		case 1:
			fmt.Fprintf(&b, "M  internal/pkg/file%d.go\n", i)
		default:
			fmt.Fprintf(&b, "?? internal/pkg/file%d.go\n", i)
		}
	}

	return b.String()
}

func BenchmarkParseStatus(b *testing.B) {
	for _, files := range []int{0, 10, 1000} {
		output := statusOutput(files)

		b.Run(fmt.Sprintf("files=%d", files), func(b *testing.B) {
			for b.Loop() {
				parseStatus(output)
			}
		})
	}
}

func TestStatus(t *testing.T) {
	gittest.Isolate(t)

	remote := gittest.NewRemote(t)
	dir := filepath.Join(t.TempDir(), "repo")
	gittest.Clone(t, remote.URL, dir)

	clean, err := Status(dir)
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}

	want := StatusSummary{Branch: "main", HasUpstream: true}
	if clean != want {
		t.Errorf("Status() of a fresh clone = %+v, want %+v", clean, want)
	}

	if !clean.Clean() {
		t.Error("Clean() = false for a fresh clone")
	}

	gittest.CommitFile(t, dir, "committed.txt", "x\n")

	if err := os.WriteFile(filepath.Join(dir, "untracked.txt"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("changed\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	changed, err := Status(dir)
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}

	want = StatusSummary{Branch: "main", HasUpstream: true, Ahead: 1, Modified: 1, Untracked: 1}
	if changed != want {
		t.Errorf("Status() with changes = %+v, want %+v", changed, want)
	}

	if changed.Clean() {
		t.Error("Clean() = true for a repository with changes")
	}
}

func TestStatusMissingDirectory(t *testing.T) {
	_, err := Status(filepath.Join(t.TempDir(), "missing"))
	if !errors.Is(err, ErrDirectoryNotFound) {
		t.Errorf("Status() error = %v, want %v", err, ErrDirectoryNotFound)
	}
}

func TestStatusNotRepository(t *testing.T) {
	gittest.Isolate(t)

	_, err := Status(t.TempDir())
	if err == nil {
		t.Fatal("Status() of a plain directory succeeded, want error")
	}

	if !strings.Contains(err.Error(), "not a git repository") {
		t.Errorf("Status() error = %q, want git's own explanation", err)
	}
}
