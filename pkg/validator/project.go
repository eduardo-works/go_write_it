package validator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/eduardo-works/writeit/pkg/project"
	"github.com/eduardo-works/writeit/pkg/resolver"
)

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

func ValidateManuscript(path string) (string, error) {
	path, err := ValidateFile(path)
	if err != nil {
		return "", err
	}

	root, err := resolver.ResolveRootFolder(path)
	if err != nil {
		return "", err
	}

	p, err := project.LoadProject(root)
	if err != nil {
		return "", err
	}

	path, err = p.Contain(path)
	if err != nil {
		return "", err
	}

	if filepath.Ext(path) != ".md" {
		return "", fmt.Errorf("The manuscript is not a markdown file: %w", project.ErrNotManuscriptFile)
	}

	return path, nil
}

func ValidatePart(path string) (int, error) {
	root, err := resolver.ResolveRootFolder(path)
	if err != nil {
		return 0, err
	}

	p, err := project.LoadProject(root)
	if err != nil {
		return 0, err
	}

	partRegex := regexp.MustCompile(`/parte_(\d{2})/[^/]+$`)
	match := partRegex.FindStringSubmatch(path)

	part := 0

	if match != nil {
		part, err = strconv.Atoi(match[1])
		if err != nil {
			return 0, err
		}
	}

	if p.IsParted() {
		if part < 1 || part > p.Parts {
			return 0, fmt.Errorf("invalid part number: %d", part)
		}
		return part, nil

	}

	if part != 0 {
		return 0, fmt.Errorf("project does not support parts: %s", path)
	}

	return part, nil
}

func ValidateTypeAndParts(projectType project.ProjectType, parts int) (project.ProjectType, int, error) {
	if parts < 0 {
		return project.Standalone, 0, fmt.Errorf("%w: number of parts cannot be negative", project.ErrInvalidValue)
	}

	if projectType == project.Standalone || projectType == project.Chaptered {
		if parts == 0 {
			return projectType, parts, nil
		}

		return project.Standalone, 0, fmt.Errorf("%w: this project type cannot have parts", project.ErrInvalidValue)
	}

	if projectType != project.Parted {
		return project.Standalone, 0, fmt.Errorf("%w: unknown project type", project.ErrInvalidValue)
	}

	if parts == 0 {
		return project.Standalone, 0, fmt.Errorf("%w: a parted project must have at least one part", project.ErrInvalidValue)
	}

	return projectType, parts, nil
}

func ValidateContext(p project.Project, path string) (string, error) {
	resolved, err := resolver.ResolvePath(path)
	if err != nil {
		return "", err
	}

	info, err := os.Stat(resolved)
	if err != nil {
		return "", err
	}

	candidate := resolved
	if info.Mode().IsRegular() {
		candidate = filepath.Dir(resolved)
	}

	candidate, err = p.Contain(candidate)
	if err != nil {
		return "", err
	}

	contextFolder, err := resolver.ResolveContext(p)
	if err != nil {
		return "", err
	}

	rel, err := filepath.Rel(contextFolder, candidate)
	if err != nil {
		return "", err
	}

	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("the context is invalid: %s", path)
	}

	if info.Mode().IsRegular() {
		return filepath.Join(candidate, filepath.Base(resolved)), nil
	}

	return candidate, nil
}
