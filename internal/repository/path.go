package repository

import (
	"path/filepath"
	"strings"
)

// resolvePath превращает путь, переданный пользователем (относительно текущей
// папки), в путь относительно корня workspace в формате конфига ("services/api").
func resolvePath(root, raw string) (string, error) {
	absolute, err := filepath.Abs(raw)
	if err != nil {
		return "", err
	}

	relative, err := filepath.Rel(root, absolute)
	if err != nil {
		return "", err
	}

	if relative == "." || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", ErrOutsideWorkspace
	}

	return filepath.ToSlash(relative), nil
}

func directory(root, path string) string {
	return filepath.Join(root, filepath.FromSlash(path))
}
