package resolver

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/eduardo-works/writeit/pkg/project"
)

func ResolveFile(path string) (string, error) {
	path, err := ResolvePath(path)
	if err != nil {
		return "", err
	}

	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if info.Mode().IsRegular() {
		return path, nil
	}

	return "", fmt.Errorf("%s: %w", path, ErrNotFile)
}

func ResolveDir(path string) (string, error) {
	path, err := ResolvePath(path)
	if err != nil {
		return "", err
	}

	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if info.IsDir() {
		return path, nil
	}

	return "", fmt.Errorf("%s: %w", path, ErrNotDir)
}

func ResolvePath(path string) (string, error) {
	path, err := ExpandPath(path)
	if err != nil {
		return "", err
	}

	_, err = os.Stat(path)

	if errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("path does not exist: %s: %w", path, err)
	}

	if errors.Is(err, os.ErrPermission) {
		return "", fmt.Errorf("permission denied: %s: %w", path, err)
	}

	if err != nil {
		return "", fmt.Errorf("cannot access path %s: %w", path, err)
	}

	return path, nil
}

func ExpandPath(path string) (string, error) {
	path = strings.TrimSpace(path)

	if path == "" {
		return "", errors.New("path cannot be empty")
	}

	if path == "~" || strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		path = filepath.Join(home, strings.TrimPrefix(path, "~/"))
	}

	return filepath.Abs(path)
}

func ResolveRootFolder(path string) (string, error) {
	resolved, err := ResolvePath(path)
	if err != nil {
		return "", err
	}

	info, err := os.Stat(resolved)
	if err != nil {
		return "", err
	}

	if !info.IsDir() {
		resolved = filepath.Dir(resolved)
	}

	for range 2 {
		candidate := filepath.Join(resolved, project.ProjectFile)

		_, err := ResolveFile(candidate)
		if err == nil {
			return resolved, nil
		}

		if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}

		parent := filepath.Dir(resolved)
		if parent == resolved {
			break
		}

		resolved = parent
	}

	return "", fmt.Errorf("path is not a project root: %s: %w", path, ErrNotProjectRoot)
}
