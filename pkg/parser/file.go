package parser

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/eduardo-works/writeit/pkg/project"
	"github.com/eduardo-works/writeit/pkg/resolver"
)

func FileParse(p project.Project, path string) (ParsedFiled, error) {
	filename := filepath.Base(path)
	context := filepath.Base(filepath.Dir(path))
	if context == "." {
		context = ""
	}

	segments := strings.Split(filename, "_")

	parsed := ParsedFiled{
		Original: path,
		Context:  context,
		Part:     -1,
		Index:    -1,
		Version:  -1,
		Slug:     "",
		IsPart:   false,
	}

	fields := map[string]*int{
		"p": &parsed.Part,
		"i": &parsed.Index,
		"v": &parsed.Version,
	}

	for _, segment := range segments {
		if slug, ok := strings.CutSuffix(segment, ".md"); ok {
			parsed.Slug = slug
		}

		for prefix, field := range fields {
			matched, _ := regexp.MatchString(`^`+prefix+`\d+$`, segment)
			if !matched {
				continue
			}

			value, err := strconv.Atoi(strings.TrimPrefix(segment, prefix))
			if err != nil {
				return ParsedFiled{}, fmt.Errorf(
					"invalid numeric segment %q (%s): %w",
					segment,
					path,
					project.ErrInvalidValue,
				)
			}

			*field = value
			break
		}
	}

	matched, _ := regexp.MatchString(`^parte-\d+$`, parsed.Slug)
	parsed.IsPart = p.IsParted() && matched

	var err error

	parsed, err = ensurePart(p, parsed)
	if err != nil {
		return ParsedFiled{}, err
	}

	parsed, err = ensureContext(p, parsed)
	if err != nil {
		return ParsedFiled{}, err
	}

	parsed, err = ensureVersion(parsed)
	if err != nil {
		return ParsedFiled{}, err
	}

	parsed, err = ensureIndex(p, parsed)
	if err != nil {
		return ParsedFiled{}, err
	}

	return parsed, nil
}

func ensurePart(p project.Project, parsed ParsedFiled) (ParsedFiled, error) {
	if !p.IsParted() {
		if parsed.Part >= 0 {
			return ParsedFiled{}, fmt.Errorf(
				"file should not have parts (%s): %w",
				parsed.Original,
				project.ErrInvalidValue,
			)
		}

		return parsed, nil
	}

	if parsed.Part > 0 {
		if parsed.Part > p.Parts {
			return ParsedFiled{}, fmt.Errorf(
				"part %d is outside the project range (%s): %w",
				parsed.Part,
				parsed.Original,
				project.ErrInvalidValue,
			)
		}

		return parsed, nil
	}

	if matched, _ := regexp.MatchString(`^parte_\d+$`, parsed.Context); matched {
		part, _ := strings.CutPrefix(parsed.Context, "parte_")

		value, err := strconv.Atoi(part)
		if err == nil && value > 0 && value <= p.Parts {
			parsed.Part = value
			return parsed, nil
		}
	}

	if parsed.IsPart {
		part, _ := strings.CutPrefix(parsed.Slug, "parte-")

		value, err := strconv.Atoi(part)
		if err == nil && value > 0 && value <= p.Parts {
			parsed.Part = value
			return parsed, nil
		}
	}

	return ParsedFiled{}, fmt.Errorf(
		"impossible to determine a valid part (%s): %w",
		parsed.Original,
		project.ErrInvalidValue,
	)
}

func ensureContext(p project.Project, parsed ParsedFiled) (ParsedFiled, error) {
	if parsed.Context != "" {
		return parsed, nil
	}

	context, err := resolver.ResolveContext(p, strconv.Itoa(parsed.Part))
	if err != nil {
		return ParsedFiled{}, fmt.Errorf(
			"impossible to determine the context (%s): %w",
			parsed.Original,
			project.ErrInvalidValue,
		)
	}

	parsed.Context = filepath.Base(context)

	return parsed, nil
}

func ensureVersion(parsed ParsedFiled) (ParsedFiled, error) {
	if parsed.IsPart {
		if parsed.Version != -1 {
			return ParsedFiled{}, fmt.Errorf(
				"part files cannot have versions (%s): %w",
				parsed.Original,
				project.ErrInvalidValue,
			)
		}

		return parsed, nil
	}

	if parsed.Version <= 0 {
		return ParsedFiled{}, fmt.Errorf(
			"file must have a valid version (%s): %w",
			parsed.Original,
			project.ErrInvalidValue,
		)
	}

	return parsed, nil
}

func ensureIndex(p project.Project, parsed ParsedFiled) (ParsedFiled, error) {
	if p.IsStandalone() {
		parsed.Index = -1
		return parsed, nil
	}

	if parsed.IsPart {
		parsed.Index = 0
		return parsed, nil
	}

	if parsed.Index <= 0 {
		return ParsedFiled{}, fmt.Errorf(
			"file must have a positive index (%s): %w",
			parsed.Original,
			project.ErrInvalidValue,
		)
	}

	return parsed, nil
}
