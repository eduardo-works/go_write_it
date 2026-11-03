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

	return "", fmt.Errorf("%s: %w", path, project.ErrNotFile)
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

	return "", fmt.Errorf("%s: %w", path, project.ErrNotDir)
}

func ResolveContainer(path string) (string, error) {
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

	if info.Mode().IsRegular() {
		return filepath.Dir(path), nil
	}

	return "", fmt.Errorf("%s: %w", path, project.ErrNotDir)
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
