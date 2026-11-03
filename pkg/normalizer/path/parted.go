package path

import (
	"fmt"
	"path/filepath"

	"github.com/eduardo-works/writeit/pkg/naming"
	"github.com/eduardo-works/writeit/pkg/parser"
	"github.com/eduardo-works/writeit/pkg/project"
	"github.com/eduardo-works/writeit/pkg/validator"
)

func NormalizeParted(p *project.Project, path string) NormalizedPath {
	result := NormalizedPath{
		Project: p,
		Source:  path,
		Target:  path,
	}

	if !p.IsParted() {
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

	var normalizedPath string

	if parsed.IsPart {
		normalizedPath = normalizeIsPart(parsed)
	} else {
		normalizedPath, err = normalizeIsNotPart(p, parsed)
		if err != nil {
			result.Issue = err
			return result
		}
	}

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

func normalizeIsPart(parsed parser.ParsedFiled) string {
	name := naming.PartedPart(parsed.Part)
	return filepath.Join(filepath.Dir(parsed.Original), name)
}

func normalizeIsNotPart(p *project.Project, parsed parser.ParsedFiled) (string, error) {
	name, err := naming.Manuscript(naming.NamingParameters{
		Project: p,
		Part:    parsed.Part,
		Index:   parsed.Index,
		Version: parsed.Version,
		Title:   parsed.Slug,
	})
	if err != nil {
		return "", err
	}

	return filepath.Join(filepath.Dir(parsed.Original), name), nil
}
