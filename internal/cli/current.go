package cli

import (
	"errors"

	"github.com/amatsu-tech/amr/internal/cli/output"
	"github.com/amatsu-tech/amr/internal/workspace"
)

func currentWorkspace() (string, error) {
	root, err := workspace.Root()
	if err != nil {
		if errors.Is(err, workspace.ErrNotFound) {
			return "", output.WithHint(err, "Initialize one here with:\n  amr init")
		}

		return "", err
	}

	return root, nil
}
