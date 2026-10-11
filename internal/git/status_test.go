package git

import (
	"fmt"
	"strings"
	"testing"
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
