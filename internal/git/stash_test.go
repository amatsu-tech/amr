package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/amatsu-tech/amr/internal/gittest"
)

func TestStash(t *testing.T) {
	gittest.Isolate(t)

	dir := filepath.Join(t.TempDir(), "repo")
	gittest.InitRepo(t, dir)

	write := func(name, content string) {
		t.Helper()

		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write("README.md", "changed\n")
	write("new-file.txt", "new\n")
	write(".gitignore", "ignored.txt\n")
	gittest.Run(t, dir, "add", ".gitignore")
	gittest.Run(t, dir, "commit", "-q", "-m", "ignore")
	write("ignored.txt", "secret\n")

	if output, err := StashPush(dir, "first"); err != nil {
		t.Fatalf("StashPush() error = %v, output: %s", err, output)
	}

	// Новый файл спрятан вместе с изменениями, игнорируемый остался на месте.
	status, err := Status(dir)
	if err != nil {
		t.Fatal(err)
	}

	if !status.Clean() {
		t.Errorf("after StashPush status = %+v, want clean", status)
	}

	if _, err := os.Stat(filepath.Join(dir, "ignored.txt")); err != nil {
		t.Error("ignored file was stashed, want it left alone")
	}

	write("README.md", "changed again\n")

	if _, err := StashPush(dir, "second"); err != nil {
		t.Fatal(err)
	}

	stashes, err := StashList(dir)
	if err != nil {
		t.Fatal(err)
	}

	want := []Stash{{Ref: "stash@{0}", Message: "second"}, {Ref: "stash@{1}", Message: "first"}}
	if len(stashes) != 2 || stashes[0] != want[0] || stashes[1] != want[1] {
		t.Fatalf("StashList() = %+v, want %+v", stashes, want)
	}

	if output, err := StashPop(dir, "stash@{1}"); err != nil {
		t.Fatalf("StashPop() error = %v, output: %s", err, output)
	}

	if _, err := os.Stat(filepath.Join(dir, "new-file.txt")); err != nil {
		t.Error("new file was not restored by StashPop")
	}

	if _, err := StashDrop(dir, "stash@{0}"); err != nil {
		t.Fatal(err)
	}

	if stashes, _ := StashList(dir); len(stashes) != 0 {
		t.Errorf("StashList() after pop and drop = %+v, want empty", stashes)
	}
}

func TestStashListEmpty(t *testing.T) {
	gittest.Isolate(t)

	dir := filepath.Join(t.TempDir(), "repo")
	gittest.InitRepo(t, dir)

	stashes, err := StashList(dir)
	if err != nil {
		t.Fatal(err)
	}

	if len(stashes) != 0 {
		t.Errorf("StashList() = %+v, want empty", stashes)
	}
}
