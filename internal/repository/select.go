package repository

import (
	"fmt"
	"slices"
	"strings"
)

// Select оставляет из repositories те, что указаны путём или входят в одну из групп.
// Без путей и групп возвращает все. Порядок — как в конфиге, без повторов.
// Пути считаются относительно текущей папки, как в Add. Путь внутри репозитория
// (например, "." из его подпапки) выбирает этот репозиторий.
func Select(root string, repositories []Repository, paths, groups []string) ([]Repository, error) {
	if len(paths) == 0 && len(groups) == 0 {
		return repositories, nil
	}

	wanted := make(map[string]bool)

	for _, rawPath := range paths {
		path, err := resolvePath(root, rawPath)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", err, rawPath)
		}

		repo, ok := containingRepository(repositories, path)
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrRepositoryNotTracked, path)
		}

		wanted[repo.Path] = true
	}

	for _, group := range groups {
		found := false

		for _, repo := range repositories {
			if slices.Contains(repo.Groups, group) {
				wanted[repo.Path] = true
				found = true
			}
		}

		if !found {
			return nil, fmt.Errorf("%w: %s", ErrUnknownGroup, group)
		}
	}

	return slices.DeleteFunc(slices.Clone(repositories), func(repo Repository) bool {
		return !wanted[repo.Path]
	}), nil
}

// containingRepository находит репозиторий, внутри которого лежит path.
// Если репозитории вложены друг в друга, берётся самый глубокий.
func containingRepository(repositories []Repository, path string) (Repository, bool) {
	var found Repository

	for _, repo := range repositories {
		inside := path == repo.Path || strings.HasPrefix(path, repo.Path+"/")
		if inside && len(repo.Path) > len(found.Path) {
			found = repo
		}
	}

	return found, found.Path != ""
}
