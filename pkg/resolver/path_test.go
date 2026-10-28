package resolver

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/eduardo-works/writeit/pkg/project"
)

func TestExpandPath(t *testing.T) {
	t.Run("empty path", func(t *testing.T) {
		_, err := ExpandPath("   ")
		if err == nil {
			t.Fatal("expected error for empty path")
		}
	})

	t.Run("absolute path", func(t *testing.T) {
		dir := t.TempDir()

		got, err := ExpandPath(dir)
		if err != nil {
			t.Fatal(err)
		}

		if got != dir {
			t.Fatalf("expected %q, got %q", dir, got)
		}
	})

	t.Run("relative path", func(t *testing.T) {
		got, err := ExpandPath(".")
		if err != nil {
			t.Fatal(err)
		}

		expected, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}

		if got != expected {
			t.Fatalf("expected %q, got %q", expected, got)
		}
	})

	t.Run("home directory", func(t *testing.T) {
		home, err := os.UserHomeDir()
		if err != nil {
			t.Fatal(err)
		}

		got, err := ExpandPath("~/books")
		if err != nil {
			t.Fatal(err)
		}

		expected := filepath.Join(home, "books")
		if got != expected {
			t.Fatalf("expected %q, got %q", expected, got)
		}
	})
}

func TestResolvePath(t *testing.T) {
	t.Run("existing directory", func(t *testing.T) {
		dir := t.TempDir()

		got, err := ResolvePath(dir)
		if err != nil {
			t.Fatal(err)
		}

		if got != dir {
			t.Fatalf("expected %q, got %q", dir, got)
		}
	})

	t.Run("nonexistent path", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "missing")

		_, err := ResolvePath(path)
		if !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("expected os.ErrNotExist, got %v", err)
		}
	})
}

func TestResolveFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "chapter.md")

	if err := os.WriteFile(file, []byte("Chapter"), 0600); err != nil {
		t.Fatal(err)
	}

	t.Run("valid file", func(t *testing.T) {
		got, err := ResolveFile(file)
		if err != nil {
			t.Fatal(err)
		}

		if got != file {
			t.Fatalf("expected %q, got %q", file, got)
		}
	})

	t.Run("directory is not a file", func(t *testing.T) {
		_, err := ResolveFile(dir)
		if !errors.Is(err, ErrNotFile) {
			t.Fatalf("expected ErrNotFile, got %v", err)
		}
	})

	t.Run("nonexistent file", func(t *testing.T) {
		_, err := ResolveFile(filepath.Join(dir, "missing.md"))
		if !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("expected os.ErrNotExist, got %v", err)
		}
	})
}

func TestResolveDir(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "chapter.md")

	if err := os.WriteFile(file, []byte("Chapter"), 0600); err != nil {
		t.Fatal(err)
	}

	t.Run("valid directory", func(t *testing.T) {
		got, err := ResolveDir(dir)
		if err != nil {
			t.Fatal(err)
		}

		if got != dir {
			t.Fatalf("expected %q, got %q", dir, got)
		}
	})

	t.Run("file is not a directory", func(t *testing.T) {
		_, err := ResolveDir(file)
		if !errors.Is(err, ErrNotDir) {
			t.Fatalf("expected ErrNotDir, got %v", err)
		}
	})
}

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
		if !errors.Is(err, ErrNotProjectRoot) {
			t.Fatalf("expected ErrNotProjectRoot, got %v", err)
		}
	})
}
