package path

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/eduardo-works/writeit/pkg/parser"
	"github.com/eduardo-works/writeit/pkg/project"
	"github.com/eduardo-works/writeit/pkg/validator"
)

func Normalize(p *project.Project, path string) NormalizedPath {
	switch {
	case p.IsStandalone():
		return NormalizeStandAlone(p, path)

	case p.IsChaptered():
		return NormalizeChaptered(p, path)

	case p.IsParted():
		return NormalizeParted(p, path)

	default:
		return NormalizedPath{
			Project: p,
			Source:  path,
			Target:  path,
			Issue: fmt.Errorf(
				"invalid project type (%s): %w",
				p.ProjectType,
				project.ErrInvalidValue,
			),
		}
	}
}

func NormalizeStandAlone(p *project.Project, path string) NormalizedPath {
	result := NormalizedPath{
		Project: p,
		Source:  path,
		Target:  path,
	}

	if !p.IsStandalone() {
		result.Prepare()
		return result
	}

	validatedPath, err := validator.ValidateFile(path)
	if err != nil {
		result.Issue = err
		return result
	}

	result.Source = validatedPath
	result.Target = validatedPath

	parsed, err := parser.FileParse(*p, validatedPath)
	if err != nil {
		result.Issue = err
		return result
	}

	name := fmt.Sprintf("v%02d_%s.md", parsed.Version, p.Slug)
	normalizedPath := filepath.Join(filepath.Dir(validatedPath), name)

	can, err := canNormalize(validatedPath, normalizedPath)
	if !can {
		result.Issue = fmt.Errorf(
			"cannot normalize %s: %w",
			validatedPath,
			err,
		)
		return result
	}

	result.Target = normalizedPath
	result.Prepare()

	return result
}

func canNormalize(path, normalizedPath string) (bool, error) {
	if path == normalizedPath {
		return true, nil
	}

	_, err := os.Stat(normalizedPath)

	if errors.Is(err, os.ErrNotExist) {
		return true, nil
	}

	if err != nil {
		return false, err
	}

	return false, fmt.Errorf(
		"destination file already exists: %s",
		normalizedPath,
	)
}
