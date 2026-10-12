package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/amatsu-tech/amr/internal/gittest"
)

func TestFetchAndMergeFastForward(t *testing.T) {
	gittest.Isolate(t)

	remote := gittest.NewRemote(t)
	dir := filepath.Join(t.TempDir(), "repo")
	gittest.Clone(t, remote.URL, dir)
	remote.Push(t, "new-file.txt")

	if output, err := Fetch(dir); err != nil {
		t.Fatalf("Fetch() error = %v, output: %s", err, output)
	}

	status, err := Status(dir)
	if err != nil {
		t.Fatal(err)
	}

	if status.Behind != 1 {
		t.Fatalf("after Fetch() Behind = %d, want 1", status.Behind)
	}

	if output, err := MergeFastForward(dir); err != nil {
		t.Fatalf("MergeFastForward() error = %v, output: %s", err, output)
	}

	if _, err := os.Stat(filepath.Join(dir, "new-file.txt")); err != nil {
		t.Errorf("merged file is missing: %v", err)
	}
}

func TestMergeFastForwardRefusesDivergedHistory(t *testing.T) {
	gittest.Isolate(t)

	remote := gittest.NewRemote(t)
	dir := filepath.Join(t.TempDir(), "repo")
	gittest.Clone(t, remote.URL, dir)

	gittest.CommitFile(t, dir, "local-commit.txt", "local\n")
	remote.Push(t, "remote-commit.txt")

	if output, err := Fetch(dir); err != nil {
		t.Fatalf("Fetch() error = %v, output: %s", err, output)
	}

	head := gittest.Run(t, dir, "rev-parse", "HEAD")

	output, err := MergeFastForward(dir)
	if err == nil {
		t.Fatal("MergeFastForward() of diverged history succeeded, want error")
	}

	if output == "" {
		t.Error("MergeFastForward() returned no git output to explain the failure")
	}

	if got := gittest.Run(t, dir, "rev-parse", "HEAD"); got != head {
		t.Error("MergeFastForward() moved HEAD despite failing")
	}
}

func TestFetchFailure(t *testing.T) {
	gittest.Isolate(t)

	remote := gittest.NewRemote(t)
	dir := filepath.Join(t.TempDir(), "repo")
	gittest.Clone(t, remote.URL, dir)

	// "Сервер" пропал: удаляем bare-репозиторий.
	if err := os.RemoveAll(remote.URL); err != nil {
		t.Fatal(err)
	}

	if _, err := Fetch(dir); err == nil {
		t.Error("Fetch() from a deleted remote succeeded, want error")
	}
}
