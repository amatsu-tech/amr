package repository

import (
	"fmt"
	"slices"
)

// Select оставляет из repositories те, что указаны путём или входят в одну из групп.
// Без путей и групп возвращает все. Порядок — как в конфиге, без повторов.
// Пути считаются относительно текущей папки, как в Add.
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

		if !slices.ContainsFunc(repositories, func(repo Repository) bool { return repo.Path == path }) {
			return nil, fmt.Errorf("%w: %s", ErrRepositoryNotTracked, path)
		}

		wanted[path] = true
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
