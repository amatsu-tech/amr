package repository

import (
	"errors"
	"strings"
	"testing"
)

func TestGitFailure(t *testing.T) {
	exitErr := errors.New("exit status 1")

	tests := []struct {
		name   string
		output string
		want   string
	}{
		{
			name: "untracked files already exist",
			output: "a.png already exists, no checkout\n" +
				"b.png already exists, no checkout\n" +
				"error: could not restore untracked files from stash\n" +
				"On branch main\n\tmodified:   README.md\n" +
				"The stash entry is kept in case you need it again.\n",
			want: "restored partly: 2 new files from the stash already exist in the folder, the stash is kept",
		},
		{
			name: "merge conflicts",
			output: "Auto-merging README.md\n" +
				"CONFLICT (content): Merge conflict in README.md\n" +
				"The stash entry is kept in case you need it again.\n",
			want: "1 file has merge conflicts, the stash is kept",
		},
		{
			name:   "one untracked file already exists",
			output: "a.png already exists, no checkout\nerror: could not restore untracked files from stash\n",
			want:   "restored partly: 1 new file from the stash already exists in the folder, the stash is kept",
		},
		{
			name:   "repository not found",
			output: "ERROR: Repository not found.\nfatal: Could not read from remote repository.\n",
			want:   "Repository not found.",
		},
		{
			name:   "fast-forward impossible",
			output: "hint: Diverging branches can't be fast-forwarded\nfatal: Not possible to fast-forward, aborting.\n",
			want:   "Not possible to fast-forward, aborting.",
		},
		{
			name:   "plain message",
			output: "\nsomething went wrong\nmore details\n",
			want:   "something went wrong",
		},
		{
			name:   "no output",
			output: "",
			want:   "exit status 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gitFailure(tt.output, exitErr); got != tt.want {
				t.Errorf("gitFailure() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGitFailureLongLine(t *testing.T) {
	got := gitFailure("fatal: "+strings.Repeat("x", 500), errors.New("exit status 1"))

	if len([]rune(got)) != maxReasonLength {
		t.Errorf("reason length = %d, want %d", len([]rune(got)), maxReasonLength)
	}

	if !strings.HasSuffix(got, "…") {
		t.Errorf("long reason %q does not end with …", got)
	}
}
