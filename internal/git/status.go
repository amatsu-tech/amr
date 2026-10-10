package git

import (
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var ErrDirectoryNotFound = errors.New("repository directory not found")

type StatusSummary struct {
	Branch      string
	HasUpstream bool
	Staged      int
	Modified    int
	Untracked   int
	Conflicted  int
	Ahead       int
	Behind      int
}

func (s StatusSummary) Clean() bool {
	return s.Staged == 0 && s.Modified == 0 && s.Untracked == 0 && s.Conflicted == 0
}

var (
	aheadPattern  = regexp.MustCompile(`ahead (\d+)`)
	behindPattern = regexp.MustCompile(`behind (\d+)`)
)

func isChangedStatus(status byte) bool {
	return status != ' ' && status != '?'
}

func parseBranch(header string) string {
	rest := strings.TrimPrefix(header, "## ")

	if rest == "HEAD (no branch)" {
		return "detached"
	}

	if branch, found := strings.CutPrefix(rest, "No commits yet on "); found {
		return branch
	}

	if i := strings.Index(rest, "..."); i != -1 {
		return rest[:i]
	}

	return rest
}

func Status(directory string) (StatusSummary, error) {
	if _, err := os.Stat(directory); err != nil {
		if os.IsNotExist(err) {
			return StatusSummary{}, ErrDirectoryNotFound
		}

		return StatusSummary{}, err
	}

	output, err := exec.Command(
		"git",
		"-C",
		directory,
		"status",
		"--short",
		"--branch",
	).Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && len(exitErr.Stderr) > 0 {
			return StatusSummary{}, errors.New(strings.TrimSpace(string(exitErr.Stderr)))
		}

		return StatusSummary{}, err
	}

	var summary StatusSummary

	trimmed := strings.TrimSpace(string(output))
	if trimmed == "" {
		return summary, nil
	}

	for _, line := range strings.Split(trimmed, "\n") {
		if strings.HasPrefix(line, "## ") {
			summary.Branch = parseBranch(line)
			summary.HasUpstream = strings.Contains(line, "...")

			if match := aheadPattern.FindStringSubmatch(line); match != nil {
				summary.Ahead, _ = strconv.Atoi(match[1])
			}

			if match := behindPattern.FindStringSubmatch(line); match != nil {
				summary.Behind, _ = strconv.Atoi(match[1])
			}

			continue
		}

		index, worktree := line[0], line[1]

		switch {
		case index == '?' && worktree == '?':
			summary.Untracked++
		case index == 'U' || worktree == 'U' ||
			(index == 'A' && worktree == 'A') ||
			(index == 'D' && worktree == 'D'):
			summary.Conflicted++
		default:
			if isChangedStatus(index) {
				summary.Staged++
			}
			if isChangedStatus(worktree) {
				summary.Modified++
			}
		}
	}

	return summary, nil
}
