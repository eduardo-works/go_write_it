package validator

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/eduardo-works/writeit/pkg/project"
	"github.com/eduardo-works/writeit/pkg/resolver"
)

func TestValidateRootFolder(t *testing.T) {
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
			got, err := ValidateRootFolder(tt.path)
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

		_, err := resolver.ResolveRootFolder(dir)
		if !errors.Is(err, project.ErrNotProjectRoot) {
			t.Fatalf("expected ErrNotProjectRoot, got %v", err)
		}
	})
}

func TestValidateRootFolderCandidate(t *testing.T) {
	t.Run("valid directory", func(t *testing.T) {
		dir := t.TempDir()

		got, err := ValidateRootFolderCandidate(dir)
		if err != nil {
			t.Fatal(err)
		}

		if got != dir {
			t.Fatalf("expected %q, got %q", dir, got)
		}
	})

	t.Run("directory is already a project", func(t *testing.T) {
		dir := t.TempDir()
		marker := filepath.Join(dir, project.ProjectFile)

		if err := os.WriteFile(marker, []byte("{}"), 0600); err != nil {
			t.Fatal(err)
		}

		_, err := ValidateRootFolderCandidate(dir)
		if err == nil {
			t.Fatal("expected error for existing project")
		}
	})

	t.Run("directory inside existing project", func(t *testing.T) {
		root := t.TempDir()
		child := filepath.Join(root, "parte_01")

		if err := os.Mkdir(child, 0700); err != nil {
			t.Fatal(err)
		}

		marker := filepath.Join(root, project.ProjectFile)
		if err := os.WriteFile(marker, []byte("{}"), 0600); err != nil {
			t.Fatal(err)
		}

		_, err := ValidateRootFolderCandidate(child)
		if err == nil {
			t.Fatal("expected error for directory inside existing project")
		}
	})

	t.Run("nested directory inside existing project", func(t *testing.T) {
		root := t.TempDir()
		nested := filepath.Join(root, "parte_01", "chapters")

		if err := os.MkdirAll(nested, 0700); err != nil {
			t.Fatal(err)
		}

		marker := filepath.Join(root, project.ProjectFile)
		if err := os.WriteFile(marker, []byte("{}"), 0600); err != nil {
			t.Fatal(err)
		}

		_, err := ValidateRootFolderCandidate(nested)
		if err == nil {
			t.Fatal("expected error for nested directory inside existing project")
		}
	})

	t.Run("nonexistent directory", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "missing")

		_, err := ValidateRootFolderCandidate(path)
		if !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("expected os.ErrNotExist, got %v", err)
		}
	})

	t.Run("file is not a directory", func(t *testing.T) {
		dir := t.TempDir()
		file := filepath.Join(dir, "chapter.md")

		if err := os.WriteFile(file, []byte("Chapter"), 0600); err != nil {
			t.Fatal(err)
		}

		_, err := ValidateRootFolderCandidate(file)
		if !errors.Is(err, project.ErrNotDir) {
			t.Fatalf("expected ErrNotDir, got %v", err)
		}
	})
}

func TestValidateManuscript(t *testing.T) {
	root := t.TempDir()
	marker := filepath.Join(root, project.ProjectFile)

	if err := os.WriteFile(marker, []byte("{}"), 0600); err != nil {
		t.Fatal(err)
	}

	valid := filepath.Join(root, "chapter.md")
	invalidExtension := filepath.Join(root, "chapter.txt")
	uppercaseExtension := filepath.Join(root, "chapter.MD")
	outside := filepath.Join(t.TempDir(), "outside.md")

	for _, path := range []string{valid, invalidExtension, uppercaseExtension, outside} {
		if err := os.WriteFile(path, []byte("Chapter"), 0600); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("valid markdown file", func(t *testing.T) {
		got, err := ValidateManuscript(valid)
		if err != nil {
			t.Fatal(err)
		}

		if got != valid {
			t.Fatalf("expected %q, got %q", valid, got)
		}
	})

	t.Run("invalid extension", func(t *testing.T) {
		_, err := ValidateManuscript(invalidExtension)
		if !errors.Is(err, project.ErrNotManuscriptFile) {
			t.Fatalf("expected ErrNotManuscriptFile, got %v", err)
		}
	})

	t.Run("uppercase extension", func(t *testing.T) {
		_, err := ValidateManuscript(uppercaseExtension)
		if !errors.Is(err, project.ErrNotManuscriptFile) {
			t.Fatalf("expected ErrNotManuscriptFile, got %v", err)
		}
	})

	t.Run("nonexistent file", func(t *testing.T) {
		path := filepath.Join(root, "missing.md")

		_, err := ValidateManuscript(path)
		if !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("expected os.ErrNotExist, got %v", err)
		}
	})

	t.Run("directory instead of file", func(t *testing.T) {
		_, err := ValidateManuscript(root)
		if err == nil {
			t.Fatal("expected error for directory")
		}
	})

	t.Run("file outside project", func(t *testing.T) {
		_, err := ValidateManuscript(outside)
		if err == nil {
			t.Fatal("expected error for file outside project")
		}
	})
}

func TestValidatePart(t *testing.T) {
	tests := []struct {
		name        string
		projectType project.ProjectType
		parts       int
		folder      string
		want        int
		wantErr     bool
	}{
		{"parted first part", project.Parted, 3, "parte_01", 1, false},
		{"parted last part", project.Parted, 3, "parte_03", 3, false},
		{"parted invalid part", project.Parted, 3, "parte_04", 0, true},
		{"parted zero part", project.Parted, 3, "parte_00", 0, true},
		{"parted without part", project.Parted, 3, "", 0, true},
		{"chaptered without part", project.Chaptered, 0, "", 0, false},
		{"chaptered with part", project.Chaptered, 0, "parte_01", 0, true},
		{"standalone without part", project.Standalone, 0, "", 0, false},
		{"standalone with part", project.Standalone, 0, "parte_01", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()

			p := project.Project{
				RootFolder:  root,
				ProjectType: tt.projectType,
				Parts:       tt.parts,
			}

			if err := p.Save(); err != nil {
				t.Fatal(err)
			}

			dir := root
			if tt.folder != "" {
				dir = filepath.Join(root, tt.folder)
				if err := os.Mkdir(dir, 0700); err != nil {
					t.Fatal(err)
				}
			}

			path := filepath.Join(dir, "chapter.md")

			if err := os.WriteFile(path, []byte("Chapter"), 0600); err != nil {
				t.Fatal(err)
			}

			got, err := ValidatePart(path)

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
				t.Fatalf("expected part %d, got %d", tt.want, got)
			}
		})
	}
}

func TestProjectContain(t *testing.T) {
	root := t.TempDir()

	p := project.Project{
		RootFolder:  root,
		ProjectType: project.Standalone,
	}

	if err := p.Save(); err != nil {
		t.Fatal(err)
	}

	dir := filepath.Join(root, "texto")
	if err := os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(dir, "lala.md")
	if err := os.WriteFile(path, []byte("Manuscript"), 0600); err != nil {
		t.Fatal(err)
	}

	got, err := p.Contain(path)
	if err != nil {
		t.Fatal(err)
	}

	if got != path {
		t.Fatalf("expected %q, got %q", path, got)
	}
}
