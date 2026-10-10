package workspace

import "errors"

var (
	ErrAlreadyInitialized       = errors.New("workspace is already initialized")
	ErrNotFound                 = errors.New("workspace not found")
	ErrUnsupportedConfigVersion = errors.New("unsupported workspace config version")
)
