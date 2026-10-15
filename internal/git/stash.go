package git

import (
	"bytes"
	"os/exec"
	"strings"
)

// Stash — запись из git stash list.
type Stash struct {
	// Ссылка на стеш, например "stash@{0}". Номер сдвигается при каждом новом стеше.
	Ref     string
	Message string
}

// StashPush прячет все изменения, включая новые файлы. Игнорируемые файлы не трогаются.
func StashPush(directory, message string) (string, error) {
	return run(directory, "stash", "push", "--include-untracked", "--message", message)
}

// StashList возвращает стеши от новых к старым.
func StashList(directory string) ([]Stash, error) {
	output, err := run(directory, "stash", "list", "--format=%gd%x00%gs")
	if err != nil {
		return nil, err
	}

	var stashes []Stash

	for line := range strings.Lines(output) {
		ref, subject, ok := strings.Cut(strings.TrimRight(line, "\n"), "\x00")
		if !ok {
			continue
		}

		// subject выглядит как "On main: сообщение"; имя ветки не может содержать ':'.
		_, message, found := strings.Cut(subject, ": ")
		if !found {
			message = subject
		}

		stashes = append(stashes, Stash{Ref: ref, Message: message})
	}

	return stashes, nil
}

func StashPop(directory, ref string) (string, error) {
	return run(directory, "stash", "pop", ref)
}

func StashDrop(directory, ref string) (string, error) {
	return run(directory, "stash", "drop", ref)
}

// run выполняет git в directory и возвращает его вывод без пробелов по краям.
func run(directory string, args ...string) (string, error) {
	var buf bytes.Buffer

	cmd := exec.Command("git", append([]string{"-C", directory}, args...)...)
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err := cmd.Run()

	return strings.TrimSpace(buf.String()), err
}
