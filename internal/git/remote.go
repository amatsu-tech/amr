package git

import (
	"errors"
	"os/exec"
	"strings"
)

// RemoteURL возвращает адрес origin. Если origin не настроен, возвращает "".
func RemoteURL(directory string) (string, error) {
	output, err := exec.Command("git", "-C", directory, "config", "--get", "remote.origin.url").Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			return "", nil
		}

		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}
