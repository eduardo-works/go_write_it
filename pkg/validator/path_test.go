package validator

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/eduardo-works/writeit/pkg/project"
	"github.com/eduardo-works/writeit/pkg/resolver"
)

func TestExpandPath(t *testing.T) {
	t.Run("empty path", func(t *testing.T) {
		_, err := resolver.ExpandPath("   ")
		if err == nil {
			t.Fatal("expected error for empty path")
		}
	})

	t.Run("absolute path", func(t *testing.T) {
		dir := t.TempDir()

		got, err := resolver.ExpandPath(dir)
		if err != nil {
			t.Fatal(err)
		}

		if got != dir {
			t.Fatalf("expected %q, got %q", dir, got)
		}
	})

	t.Run("relative path", func(t *testing.T) {
		got, err := resolver.ExpandPath(".")
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

		got, err := resolver.ExpandPath("~/books")
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

		got, err := resolver.ResolvePath(dir)
		if err != nil {
			t.Fatal(err)
		}

		if got != dir {
			t.Fatalf("expected %q, got %q", dir, got)
		}
	})

	t.Run("nonexistent path", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "missing")

		_, err := resolver.ResolvePath(path)
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
		got, err := resolver.ResolveFile(file)
		if err != nil {
			t.Fatal(err)
		}

		if got != file {
			t.Fatalf("expected %q, got %q", file, got)
		}
	})

	t.Run("directory is not a file", func(t *testing.T) {
		_, err := resolver.ResolveFile(dir)
		if !errors.Is(err, project.ErrNotFile) {
			t.Fatalf("expected ErrNotFile, got %v", err)
		}
	})

	t.Run("nonexistent file", func(t *testing.T) {
		_, err := resolver.ResolveFile(filepath.Join(dir, "missing.md"))
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
		got, err := resolver.ResolveDir(dir)
		if err != nil {
			t.Fatal(err)
		}

		if got != dir {
			t.Fatalf("expected %q, got %q", dir, got)
		}
	})

	t.Run("file is not a directory", func(t *testing.T) {
		_, err := resolver.ResolveDir(file)
		if !errors.Is(err, project.ErrNotDir) {
			t.Fatalf("expected ErrNotDir, got %v", err)
		}
	})
}
