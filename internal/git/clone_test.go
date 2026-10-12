package git

import (
	"path/filepath"
	"testing"

	"github.com/amatsu-tech/amr/internal/gittest"
)

func TestClone(t *testing.T) {
	gittest.Isolate(t)

	remote := gittest.NewRemote(t)
	dir := filepath.Join(t.TempDir(), "group", "cloned-repo")

	if output, err := Clone(remote.URL, dir); err != nil {
		t.Fatalf("Clone() error = %v, output: %s", err, output)
	}

	isRepository, err := IsRepository(dir)
	if err != nil {
		t.Fatal(err)
	}

	if !isRepository {
		t.Fatal("Clone() did not create a git repository")
	}

	if got := gittest.Run(t, dir, "config", "--get", "remote.origin.url"); got != remote.URL {
		t.Errorf("origin = %q, want %q", got, remote.URL)
	}
}

func TestCloneFailure(t *testing.T) {
	gittest.Isolate(t)

	missing := filepath.Join(t.TempDir(), "missing.git")
	dir := filepath.Join(t.TempDir(), "repo")

	output, err := Clone(missing, dir)
	if err == nil {
		t.Fatal("Clone() of a missing repository succeeded, want error")
	}

	if output == "" {
		t.Error("Clone() returned no git output to explain the failure")
	}
}
