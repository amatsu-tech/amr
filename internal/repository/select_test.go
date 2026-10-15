package repository

import (
	"errors"
	"path/filepath"
	"slices"
	"testing"
)

func TestSelect(t *testing.T) {
	root := t.TempDir()

	repositories := []Repository{
		{Path: "frontend", Groups: []string{"site"}},
		{Path: "backend/catalog", Groups: []string{"backend", "idp"}},
		{Path: "backend/integrations", Groups: []string{"backend"}},
		{Path: "infra"},
	}

	abs := func(path string) string { return filepath.Join(root, filepath.FromSlash(path)) }

	tests := []struct {
		name    string
		paths   []string
		groups  []string
		want    []string
		wantErr error
	}{
		{name: "nothing selected means all", want: []string{"frontend", "backend/catalog", "backend/integrations", "infra"}},
		{name: "paths", paths: []string{abs("infra"), abs("frontend")}, want: []string{"frontend", "infra"}},
		{name: "group", groups: []string{"backend"}, want: []string{"backend/catalog", "backend/integrations"}},
		{name: "several groups", groups: []string{"site", "idp"}, want: []string{"frontend", "backend/catalog"}},
		{name: "path and group are combined", paths: []string{abs("infra")}, groups: []string{"backend"}, want: []string{"backend/catalog", "backend/integrations", "infra"}},
		{name: "overlap is not duplicated", paths: []string{abs("backend/catalog")}, groups: []string{"idp"}, want: []string{"backend/catalog"}},
		{name: "unknown path", paths: []string{abs("missing")}, wantErr: ErrRepositoryNotTracked},
		{name: "path outside workspace", paths: []string{filepath.Dir(root)}, wantErr: ErrOutsideWorkspace},
		{name: "unknown group", groups: []string{"typo"}, wantErr: ErrUnknownGroup},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selected, err := Select(root, repositories, tt.paths, tt.groups)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Select() error = %v, want %v", err, tt.wantErr)
			}

			got := make([]string, len(selected))
			for i, repo := range selected {
				got[i] = repo.Path
			}

			if tt.wantErr == nil && !slices.Equal(got, tt.want) {
				t.Errorf("Select() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelectRelativeToCurrentDirectory(t *testing.T) {
	root := t.TempDir()
	makeDir(t, root, "backend")
	t.Chdir(filepath.Join(root, "backend"))

	repositories := []Repository{{Path: "backend/catalog"}, {Path: "frontend"}}

	selected, err := Select(root, repositories, []string{"catalog"}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(selected) != 1 || selected[0].Path != "backend/catalog" {
		t.Errorf("Select() = %+v, want backend/catalog", selected)
	}
}

func TestSelectInsideRepository(t *testing.T) {
	root := t.TempDir()
	makeDir(t, root, "services/api/internal/pkg")
	t.Chdir(filepath.Join(root, "services", "api", "internal"))

	repositories := []Repository{{Path: "services/api"}, {Path: "services/api-docs"}, {Path: "frontend"}}

	selected, err := Select(root, repositories, []string{"."}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(selected) != 1 || selected[0].Path != "services/api" {
		t.Errorf("Select(.) = %+v, want services/api", selected)
	}
}

func TestSelectNestedRepositories(t *testing.T) {
	root := t.TempDir()
	makeDir(t, root, "platform/tools/cli")
	t.Chdir(filepath.Join(root, "platform", "tools", "cli"))

	repositories := []Repository{{Path: "platform"}, {Path: "platform/tools"}}

	selected, err := Select(root, repositories, []string{"."}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(selected) != 1 || selected[0].Path != "platform/tools" {
		t.Errorf("Select(.) = %+v, want the deepest repository platform/tools", selected)
	}
}
