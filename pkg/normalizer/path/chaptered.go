package path

import (
	"fmt"
	"path/filepath"

	"github.com/eduardo-works/writeit/pkg/naming"
	"github.com/eduardo-works/writeit/pkg/parser"
	"github.com/eduardo-works/writeit/pkg/project"
	"github.com/eduardo-works/writeit/pkg/validator"
)

func NormalizeChaptered(p *project.Project, path string) NormalizedPath {
	result := NormalizedPath{
		Project: p,
		Source:  path,
		Target:  path,
	}

	if !p.IsChaptered() {
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

	name, err := naming.Manuscript(naming.NamingParameters{
		Project: p,
		Index:   parsed.Index,
		Version: parsed.Version,
		Title:   parsed.Slug,
	})
	if err != nil {
		result.Issue = err
		return result
	}

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
