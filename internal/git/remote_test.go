package git

import (
	"path/filepath"
	"testing"

	"github.com/amatsu-tech/amr/internal/gittest"
)

func TestRemoteURL(t *testing.T) {
	gittest.Isolate(t)

	remote := gittest.NewRemote(t)

	cloned := filepath.Join(t.TempDir(), "cloned")
	gittest.Clone(t, remote.URL, cloned)

	local := filepath.Join(t.TempDir(), "local")
	gittest.InitRepo(t, local)

	tests := []struct {
		name string
		dir  string
		want string
	}{
		{name: "with origin", dir: cloned, want: remote.URL},
		{name: "without origin", dir: local, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RemoteURL(tt.dir)
			if err != nil {
				t.Fatalf("RemoteURL() error = %v", err)
			}

			if got != tt.want {
				t.Errorf("RemoteURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRemoteURLNotRepository(t *testing.T) {
	gittest.Isolate(t)

	if _, err := RemoteURL(filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Error("RemoteURL() for a missing directory succeeded, want error")
	}
}
