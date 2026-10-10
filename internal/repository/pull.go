package repository

import "github.com/amatsu-tech/amr/internal/git"

type PullOutcome int

const (
	PullUpdated PullOutcome = iota
	PullUpToDate
	PullSkipped
	PullFailed
)

type PullResult struct {
	Outcome PullOutcome
	Reason  string
}

func Pull(repo Repository) PullResult {
	status, err := git.Status(repo.Directory)
	if err != nil {
		return PullResult{Outcome: PullFailed, Reason: err.Error()}
	}

	if reason := dirtyReason(status); reason != "" {
		return PullResult{Outcome: PullSkipped, Reason: reason}
	}

	if !status.HasUpstream {
		return PullResult{Outcome: PullFailed, Reason: "no upstream branch configured"}
	}

	if output, err := git.Fetch(repo.Directory); err != nil {
		return PullResult{Outcome: PullFailed, Reason: firstNonEmpty(output, err.Error())}
	}

	status, err = git.Status(repo.Directory)
	if err != nil {
		return PullResult{Outcome: PullFailed, Reason: err.Error()}
	}

	if status.Behind == 0 {
		return PullResult{Outcome: PullUpToDate}
	}

	if output, err := git.MergeFastForward(repo.Directory); err != nil {
		return PullResult{Outcome: PullFailed, Reason: firstNonEmpty(output, err.Error())}
	}

	return PullResult{Outcome: PullUpdated}
}

func dirtyReason(status git.StatusSummary) string {
	if status.Conflicted > 0 {
		return "merge conflicts"
	}

	if status.Staged > 0 || status.Modified > 0 || status.Untracked > 0 {
		return "modified files"
	}

	return ""
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}

	return ""
}
