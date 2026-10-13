package repository

import (
	"errors"
	"os"

	"github.com/amatsu-tech/amr/internal/git"
)

type SyncOutcome int

const (
	SyncCloned SyncOutcome = iota
	SyncUpdated
	SyncUpToDate
	SyncSkipped
	SyncFailed
)

type SyncResult struct {
	Outcome SyncOutcome
	Reason  string
}

// Sync приводит репозиторий в актуальное состояние: клонирует, если его нет на диске,
// иначе подтягивает изменения так же, как Pull.
func Sync(repo Repository) SyncResult {
	if _, err := os.Stat(repo.Directory); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return SyncResult{Outcome: SyncFailed, Reason: err.Error()}
		}

		switch result := Clone(repo); result.Outcome {
		case CloneCloned:
			return SyncResult{Outcome: SyncCloned}
		case CloneSkipped:
			return SyncResult{Outcome: SyncSkipped, Reason: result.Reason}
		case CloneFailed:
			return SyncResult{Outcome: SyncFailed, Reason: result.Reason}
		}
	}

	switch result := Pull(repo); result.Outcome {
	case PullUpdated:
		return SyncResult{Outcome: SyncUpdated}
	case PullUpToDate:
		return SyncResult{Outcome: SyncUpToDate}
	case PullSkipped:
		return SyncResult{Outcome: SyncSkipped, Reason: result.Reason}
	default:
		return SyncResult{Outcome: SyncFailed, Reason: result.Reason}
	}
}

type WorkspaceOutcome int

const (
	// WorkspaceNotRepository — корень workspace не git-репозиторий, обновлять нечего.
	WorkspaceNotRepository WorkspaceOutcome = iota
	WorkspaceUpdated
	WorkspaceUpToDate
	WorkspaceSkipped
	WorkspaceFailed
)

type WorkspaceResult struct {
	Outcome WorkspaceOutcome
	Reason  string
}

// UpdateWorkspace подтягивает изменения в корне workspace, если он git-репозиторий
// (например, общий репозиторий с amr.yaml). Без upstream или с локальными
// изменениями корень не трогается.
func UpdateWorkspace(root string) WorkspaceResult {
	isRepository, err := git.IsRepository(root)
	if err != nil {
		return WorkspaceResult{Outcome: WorkspaceFailed, Reason: err.Error()}
	}

	if !isRepository {
		return WorkspaceResult{Outcome: WorkspaceNotRepository}
	}

	status, err := git.Status(root)
	if err != nil {
		return WorkspaceResult{Outcome: WorkspaceFailed, Reason: err.Error()}
	}

	if !status.HasUpstream {
		return WorkspaceResult{Outcome: WorkspaceSkipped, Reason: "no upstream branch configured"}
	}

	switch result := Pull(Repository{Directory: root}); result.Outcome {
	case PullUpdated:
		return WorkspaceResult{Outcome: WorkspaceUpdated}
	case PullUpToDate:
		return WorkspaceResult{Outcome: WorkspaceUpToDate}
	case PullSkipped:
		return WorkspaceResult{Outcome: WorkspaceSkipped, Reason: result.Reason}
	default:
		return WorkspaceResult{Outcome: WorkspaceFailed, Reason: result.Reason}
	}
}
