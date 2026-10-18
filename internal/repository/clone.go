package repository

import (
	"errors"
	"os"

	"github.com/amatsu-tech/amr/internal/git"
)

type CloneOutcome int

const (
	CloneCloned CloneOutcome = iota
	CloneExists
	CloneSkipped
	CloneFailed
)

type CloneResult struct {
	Outcome CloneOutcome
	Reason  string
}

func Clone(repo Repository) CloneResult {
	if _, err := os.Stat(repo.Directory); err == nil {
		return CloneResult{Outcome: CloneExists}
	} else if !errors.Is(err, os.ErrNotExist) {
		return CloneResult{Outcome: CloneFailed, Reason: err.Error()}
	}

	if repo.URL == "" {
		return CloneResult{Outcome: CloneSkipped, Reason: "no url in workspace config"}
	}

	if output, err := git.Clone(repo.URL, repo.Directory); err != nil {
		return CloneResult{Outcome: CloneFailed, Reason: gitFailure(output, err)}
	}

	return CloneResult{Outcome: CloneCloned}
}
