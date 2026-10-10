package workspace

import "path/filepath"

func workspaceDirectory(root string) string {
	return filepath.Join(root, DirectoryName)
}

func configPath(root string) string {
	return filepath.Join(workspaceDirectory(root), ConfigFileName)
}
