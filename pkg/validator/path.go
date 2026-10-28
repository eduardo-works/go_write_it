package validator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/eduardo-works/writeit/pkg/project"
	"github.com/eduardo-works/writeit/pkg/resolver"
)

func ValidateExists(path string) (string, error) {
	return resolver.ResolvePath(path)
}

func ValidateFile(path string) (string, error) {
	return resolver.ResolveFile(path)
}

func ValidateDir(path string) (string, error) {
	return resolver.ResolveDir(path)
}

func ValidateRootFolder(path string) (string, error) {
	return resolver.ResolveRootFolder(path)
}

func ValidateRootFolderCandidate(path string) (string, error) {
	resolved, err := resolver.ResolveDir(path)
	if err != nil {
		return "", err
	}

	current := resolved

	for {
		candidate := filepath.Join(current, project.ProjectFile)

		_, err := os.Stat(candidate)
		if err == nil {
			return "", fmt.Errorf("directory is inside an existing project: %s", current)
		}

		if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}

		current = parent
	}

	return resolved, nil
}
