package git

import (
	"bytes"
	"os/exec"
	"strings"
)

func Fetch(directory string) (string, error) {
	var buf bytes.Buffer

	cmd := exec.Command("git", "-C", directory, "fetch")
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err := cmd.Run()

	return strings.TrimSpace(buf.String()), err
}
