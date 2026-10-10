package git

import (
	"bytes"
	"os/exec"
	"strings"
)

// MergeFastForward fast-forwards the current branch to its upstream. It
// refuses (and returns an error) if the merge would not be a fast-forward.
func MergeFastForward(directory string) (string, error) {
	var buf bytes.Buffer

	cmd := exec.Command("git", "-C", directory, "merge", "--ff-only", "@{upstream}")
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err := cmd.Run()

	return strings.TrimSpace(buf.String()), err
}
