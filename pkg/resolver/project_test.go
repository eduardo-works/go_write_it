package resolver

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/eduardo-works/writeit/pkg/project"
)

func TestResolveRootFolder(t *testing.T) {
	root := t.TempDir()
	child := filepath.Join(root, "parte_01")
	file := filepath.Join(child, "chapter.md")
	marker := filepath.Join(root, project.ProjectFile)

	if err := os.Mkdir(child, 0700); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(file, []byte("Chapter"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(marker, []byte("{}"), 0600); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		path string
		want string
	}{
		{"project root", root, root},
		{"child directory", child, root},
		{"file inside child directory", file, root},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveRootFolder(tt.path)
			if err != nil {
				t.Fatal(err)
			}

			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}

	t.Run("not inside a project", func(t *testing.T) {
		dir := t.TempDir()

		_, err := ResolveRootFolder(dir)
		if !errors.Is(err, project.ErrNotProjectRoot) {
			t.Fatalf("expected ErrNotProjectRoot, got %v", err)
		}
	})
}

func TestResolveProject(t *testing.T) {
	root := t.TempDir()

	p := project.Project{
		RootFolder:  root,
		Title:       "Test Project",
		Slug:        "test-project",
		ProjectType: project.Standalone,
	}

	if err := p.Save(); err != nil {
		t.Fatal(err)
	}

	child := filepath.Join(root, "standalone")
	if err := os.Mkdir(child, 0700); err != nil {
		t.Fatal(err)
	}

	file := filepath.Join(child, "chapter.md")
	if err := os.WriteFile(file, []byte("Chapter"), 0600); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		path string
	}{
		{"project root", root},
		{"child directory", child},
		{"file inside child directory", file},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveProject(tt.path)
			if err != nil {
				t.Fatal(err)
			}

			if got.RootFolder != root {
				t.Errorf("expected root %q, got %q", root, got.RootFolder)
			}

			if got.Title != p.Title {
				t.Errorf("expected title %q, got %q", p.Title, got.Title)
			}

			if got.ProjectType != p.ProjectType {
				t.Errorf("expected type %q, got %q", p.ProjectType, got.ProjectType)
			}
		})
	}

	t.Run("not a project", func(t *testing.T) {
		_, err := ResolveProject(t.TempDir())

		if !errors.Is(err, project.ErrNotProjectRoot) {
			t.Fatalf("expected ErrNotProjectRoot, got %v", err)
		}
	})

	t.Run("invalid project JSON", func(t *testing.T) {
		dir := t.TempDir()
		marker := filepath.Join(dir, project.ProjectFile)

		if err := os.WriteFile(marker, []byte("{invalid"), 0600); err != nil {
			t.Fatal(err)
		}

		_, err := ResolveProject(dir)
		if err == nil {
			t.Fatal("expected error for invalid project JSON")
		}
	})
}

func TestResolveContext(t *testing.T) {
	root := t.TempDir()

	tests := []struct {
		name        string
		projectType project.ProjectType
		parts       int
		args        []string
		want        string
		wantErr     bool
	}{
		{
			name:        "standalone",
			projectType: project.Standalone,
			want:        filepath.Join(root, "standalone"),
		},
		{
			name:        "chaptered",
			projectType: project.Chaptered,
			want:        filepath.Join(root, "chaptered"),
		},
		{
			name:        "parted without part",
			projectType: project.Parted,
			parts:       3,
			want:        filepath.Join(root, project.ContextParted),
		},
		{
			name:        "parted first part",
			projectType: project.Parted,
			parts:       3,
			args:        []string{"1"},
			want:        filepath.Join(root, "parte_01"),
		},
		{
			name:        "parted last part",
			projectType: project.Parted,
			parts:       3,
			args:        []string{"3"},
			want:        filepath.Join(root, "parte_03"),
		},
		{
			name:        "parted invalid zero",
			projectType: project.Parted,
			parts:       3,
			args:        []string{"0"},
			wantErr:     true,
		},
		{
			name:        "parted exceeds parts",
			projectType: project.Parted,
			parts:       3,
			args:        []string{"4"},
			wantErr:     true,
		},
		{
			name:        "parted non numeric",
			projectType: project.Parted,
			parts:       3,
			args:        []string{"abc"},
			wantErr:     true,
		},
		{
			name:        "invalid project type",
			projectType: project.ProjectType("invalid"),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := project.Project{
				RootFolder:  root,
				ProjectType: tt.projectType,
				Parts:       tt.parts,
			}

			got, err := ResolveContext(p, tt.args...)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}

			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestResolveRootFolderTooDeep(t *testing.T) {
	root := t.TempDir()

	if err := os.WriteFile(filepath.Join(root, project.ProjectFile), []byte("{}"), 0600); err != nil {
		t.Fatal(err)
	}

	nested := filepath.Join(root, "parte_01", "chapters")
	if err := os.MkdirAll(nested, 0700); err != nil {
		t.Fatal(err)
	}

	_, err := ResolveRootFolder(nested)
	if !errors.Is(err, project.ErrNotProjectRoot) {
		t.Fatalf("expected ErrNotProjectRoot, got %v", err)
	}
}
