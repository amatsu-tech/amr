package workspace

import (
	"errors"
	"os"
	"path/filepath"
)

// Gitignore оставляет в git только конфиг workspace: сами репозитории
// не коммитятся, их скачивает amr clone.
const Gitignore = `# amr workspace: only the config is tracked, repositories are cloned with "amr clone".
/*
!/.amr/
!/.gitignore
!/README.md
`

// CreateGitignore записывает Gitignore в корень workspace и возвращает true, если записал.
// Не трогает существующий .gitignore и не создаёт его, если корень уже git-репозиторий:
// иначе файл спрятал бы от git весь код этого репозитория.
func CreateGitignore(root string) (bool, error) {
	for _, name := range []string{".gitignore", ".git"} {
		_, err := os.Stat(filepath.Join(root, name))
		if err == nil {
			return false, nil
		}

		if !errors.Is(err, os.ErrNotExist) {
			return false, err
		}
	}

	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte(Gitignore), 0o644); err != nil {
		return false, err
	}

	return true, nil
}
