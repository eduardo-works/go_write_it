package resolver

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/eduardo-works/writeit/pkg/project"
)

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

	return "", fmt.Errorf("path is not a project root: %s: %w", path, project.ErrNotProjectRoot)
}

func ResolveProject(path string) (project.Project, error) {
	rootFolder, err := ResolveRootFolder(path)
	if err != nil {
		return project.Project{}, err
	}

	p, err := project.LoadProject(rootFolder)
	if err != nil {
		return project.Project{}, err
	}

	return p, nil
}

func ResolveContext(p project.Project, parts ...string) (string, error) {
	if p.IsStandalone() {
		return filepath.Join(p.RootFolder, string(project.Standalone)), nil
	}

	if p.IsChaptered() {
		return filepath.Join(p.RootFolder, string(project.Chaptered)), nil
	}

	if !p.IsParted() {
		return "", fmt.Errorf("invalid project type: %s", p.ProjectType)
	}

	if len(parts) == 0 {
		return filepath.Join(p.RootFolder, project.ContextParted), nil
	}

	part, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", fmt.Errorf("part must be an integer: %w", err)
	}

	if part < 1 || part > p.Parts {
		return "", fmt.Errorf("part must be between 1 and %d", p.Parts)
	}

	partName := fmt.Sprintf("%s%02d", project.ContextParted, part)

	return filepath.Join(p.RootFolder, partName), nil
}
