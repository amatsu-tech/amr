package repository

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	isolateGit(t)

	dir := filepath.Join(t.TempDir(), "repo")
	initRepo(t, dir)
	repo := Repository{Path: "repo", Directory: dir}

	var stdout, stderr bytes.Buffer

	if err := Run(repo, "git", []string{"rev-parse", "--abbrev-ref", "HEAD"}, &stdout, &stderr); err != nil {
		t.Fatalf("Run() error = %v, stderr: %s", err, stderr.String())
	}

	if got := strings.TrimSpace(stdout.String()); got != "main" {
		t.Errorf("Run() output = %q, want %q", got, "main")
	}

	if err := Run(repo, "git", []string{"no-such-command"}, &stdout, &stderr); err == nil {
		t.Error("Run() with a failing command succeeded, want error")
	}
}
