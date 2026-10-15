package git

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func IsRepository(directory string) (bool, error) {
	info, err := os.Stat(filepath.Join(directory, ".git"))
	if err == nil {
		return info.IsDir(), nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err

}

// TopLevel возвращает корень git-репозитория, внутри которого лежит directory.
func TopLevel(directory string) (string, error) {
	output, err := run(directory, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", errors.New(firstLine(output))
	}

	return output, nil
}

func firstLine(text string) string {
	line, _, _ := strings.Cut(text, "\n")
	return line
}
