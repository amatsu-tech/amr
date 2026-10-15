package repository

import "errors"

var (
	ErrRepositoryNotFound   = errors.New("repository not found")
	ErrRepositoryExists     = errors.New("repository already exists")
	ErrRepositoryNotTracked = errors.New("repository is not tracked in this workspace")
	ErrNotGitRepository     = errors.New("directory is not a git repository")
	ErrOutsideWorkspace     = errors.New("path is outside the workspace")
	ErrWorkspaceRoot        = errors.New("path is the workspace root, not one of its repositories")
	ErrUnknownGroup         = errors.New("no repositories in group")
)
