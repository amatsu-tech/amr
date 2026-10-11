package repository

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestResolvePath(t *testing.T) {
	root := t.TempDir()

	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr error
	}{
		{name: "top level", raw: filepath.Join(root, "top-level-repo"), want: "top-level-repo"},
		{name: "nested", raw: filepath.Join(root, "group", "nested-repo"), want: "group/nested-repo"},
		{name: "not clean", raw: filepath.Join(root, "group", "..", "other-group", "repo"), want: "other-group/repo"},
		{name: "root itself", raw: root, wantErr: ErrOutsideWorkspace},
		{name: "parent", raw: filepath.Dir(root), wantErr: ErrOutsideWorkspace},
		{name: "sibling", raw: filepath.Join(filepath.Dir(root), "outside-repo"), wantErr: ErrOutsideWorkspace},
		{name: "dot-dot prefix is still inside", raw: filepath.Join(root, "..cache"), want: "..cache"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolvePath(root, tt.raw)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("resolvePath() error = %v, want %v", err, tt.wantErr)
			}

			if got != tt.want {
				t.Errorf("resolvePath() = %q, want %q", got, tt.want)
			}
		})
	}
}
