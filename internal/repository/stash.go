package repository

import (
	"strings"
	"time"

	"github.com/amatsu-tech/amr/internal/git"
)

// Префикс сообщения стешей amr: отличает их от обычных git stash.
const stashPrefix = "amr: "

// DefaultStashName — имя стеша, если пользователь его не задал.
func DefaultStashName(now time.Time) string {
	return "amr-" + now.Format("2006-01-02-150405")
}

type StashOutcome int

const (
	StashSaved StashOutcome = iota
	StashNoChanges
	StashFailed
)

type StashResult struct {
	Outcome StashOutcome
	Reason  string
}

// StashSave прячет все изменения репозитория, включая новые файлы, под именем name.
func StashSave(repo Repository, name string) StashResult {
	status, err := git.Status(repo.Directory)
	if err != nil {
		return StashResult{Outcome: StashFailed, Reason: err.Error()}
	}

	if status.Clean() {
		return StashResult{Outcome: StashNoChanges}
	}

	if output, err := git.StashPush(repo.Directory, stashPrefix+name); err != nil {
		return StashResult{Outcome: StashFailed, Reason: firstNonEmpty(output, err.Error())}
	}

	return StashResult{Outcome: StashSaved}
}

// StashNames возвращает имена стешей amr в репозитории, от новых к старым.
func StashNames(repo Repository) ([]string, error) {
	stashes, err := git.StashList(repo.Directory)
	if err != nil {
		return nil, err
	}

	var names []string

	for _, stash := range stashes {
		if name, ok := strings.CutPrefix(stash.Message, stashPrefix); ok {
			names = append(names, name)
		}
	}

	return names, nil
}

type StashActionOutcome int

const (
	StashDone StashActionOutcome = iota
	StashNotFound
	StashActionFailed
)

type StashActionResult struct {
	Outcome StashActionOutcome
	Reason  string
}

// StashPop возвращает изменения из стеша name и удаляет его.
// Пустое name — последний стеш amr.
func StashPop(repo Repository, name string) StashActionResult {
	return withStash(repo, name, git.StashPop)
}

// StashDrop удаляет стеш name, не возвращая изменения.
func StashDrop(repo Repository, name string) StashActionResult {
	return withStash(repo, name, git.StashDrop)
}

func withStash(repo Repository, name string, action func(directory, ref string) (string, error)) StashActionResult {
	stashes, err := git.StashList(repo.Directory)
	if err != nil {
		return StashActionResult{Outcome: StashActionFailed, Reason: err.Error()}
	}

	for _, stash := range stashes {
		stashName, ok := strings.CutPrefix(stash.Message, stashPrefix)
		if !ok || (name != "" && stashName != name) {
			continue
		}

		if output, err := action(repo.Directory, stash.Ref); err != nil {
			return StashActionResult{Outcome: StashActionFailed, Reason: firstNonEmpty(output, err.Error())}
		}

		return StashActionResult{Outcome: StashDone}
	}

	return StashActionResult{Outcome: StashNotFound}
}
