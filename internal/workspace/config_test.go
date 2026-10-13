package workspace

import (
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
)

func writeRawConfig(t *testing.T, root, content string) {
	t.Helper()

	if err := os.WriteFile(configPath(root), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestConfigRoundTrip(t *testing.T) {
	root := newWorkspace(t)

	want := Config{
		Version: WorkspaceConfigVersion,
		Repositories: []RepositoryEntry{
			{Path: "top-level-repo", URL: "git@example.com:org/top-level-repo.git"},
			{Path: "group/nested-repo", URL: "https://example.com/org/nested-repo.git"},
			{Path: "group/repo-without-origin"},
		},
	}

	if err := WriteConfig(root, want); err != nil {
		t.Fatalf("WriteConfig() error = %v", err)
	}

	got, err := ReadConfig(root)
	if err != nil {
		t.Fatalf("ReadConfig() error = %v", err)
	}

	if got.Version != want.Version || !reflect.DeepEqual(got.Repositories, want.Repositories) {
		t.Errorf("ReadConfig() = %+v, want %+v", got, want)
	}
}

func TestWriteConfigOmitsEmptyURL(t *testing.T) {
	root := newWorkspace(t)

	config := DefaultConfig()
	config.Repositories = []RepositoryEntry{{Path: "repo-without-origin"}}

	if err := WriteConfig(root, config); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(configPath(root))
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(data), "url") {
		t.Errorf("config contains url for a repository without origin:\n%s", data)
	}
}

func TestReadConfigVersion(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr error
		wantMsg string
	}{
		{
			name:    "current version",
			content: "version: 1\nrepositories: []\n",
		},
		{
			name:    "newer version",
			content: "version: 2\nrepositories: []\n",
			wantErr: ErrUnsupportedConfigVersion,
			wantMsg: "unsupported workspace config version: 2 (this amr CLI supports up to version 1)",
		},
		{
			name:    "missing version",
			content: "repositories: []\n",
			wantErr: ErrUnsupportedConfigVersion,
			wantMsg: "unsupported workspace config version: 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := newWorkspace(t)
			writeRawConfig(t, root, tt.content)

			_, err := ReadConfig(root)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ReadConfig() error = %v, want %v", err, tt.wantErr)
			}

			if tt.wantMsg != "" && err.Error() != tt.wantMsg {
				t.Errorf("error message = %q, want %q", err.Error(), tt.wantMsg)
			}
		})
	}
}

func TestReadConfigInvalidYAML(t *testing.T) {
	root := newWorkspace(t)
	writeRawConfig(t, root, "version: [not a number\n")

	if _, err := ReadConfig(root); err == nil {
		t.Fatal("ReadConfig() with broken YAML succeeded, want error")
	}
}

func TestReadConfigMissingFile(t *testing.T) {
	root := newWorkspace(t)

	if err := os.Remove(configPath(root)); err != nil {
		t.Fatal(err)
	}

	_, err := ReadConfig(root)
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("ReadConfig() error = %v, want %v", err, os.ErrNotExist)
	}
}
