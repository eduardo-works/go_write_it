package file

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/eduardo-works/writeit/pkg/finder"
	"github.com/eduardo-works/writeit/pkg/naming"
	"github.com/eduardo-works/writeit/pkg/parser"
	"github.com/eduardo-works/writeit/pkg/project"
	"github.com/eduardo-works/writeit/pkg/resolver"
	"github.com/eduardo-works/writeit/pkg/validator"
)

func FileVersion(p project.Project, source, title string) (string, string, error) {
	source, err := sourceName(p, source)
	if err != nil {
		return "", "", err
	}

	target, err := targetName(p, source, title)
	if err != nil {
		return "", "", err
	}

	content, err := readContent(source, title)
	if err != nil {
		return "", "", err
	}

	if err := saveTarget(target, content); err != nil {
		return "", "", err
	}

	if _, err := p.AddFile(target); err != nil {
		return "", "", removeTarget(target, err)
	}

	return target, source, nil
}

func sourceName(p project.Project, source string) (string, error) {
	source, err := validator.ValidateManuscript(source)
	if err != nil {
		return "", fmt.Errorf("invalid manuscript: %w", err)
	}

	if _, err := p.Contain(source); err != nil {
		return "", fmt.Errorf("manuscript does not belong to project: %w", err)
	}

	parsed, err := parser.FileParse(p, source)
	if err != nil {
		return "", fmt.Errorf("failed to parse manuscript: %w", err)
	}

	newSourceName, err := finder.FindMaxVersion(
		filepath.Dir(source),
		&p,
		parsed.Part,
		parsed.Index,
	)
	if err != nil {
		return "", fmt.Errorf("failed to find latest manuscript version: %w", err)
	}

	if newSourceName == "" {
		return "", fmt.Errorf("manuscript version not found for %s: %w", source, project.ErrNotManuscriptFile)
	}

	return newSourceName, nil
}

func targetName(p project.Project, source, title string) (string, error) {
	nextVersion, err := resolver.ResolveNextVersion(p, source)
	if err != nil {
		return "", fmt.Errorf("failed to resolve next version: %w", err)
	}

	parsed, err := parser.FileParse(p, source)
	if err != nil {
		return "", fmt.Errorf("failed to parse manuscript: %w", err)
	}

	params, err := naming.FromParsed(&p, parsed)
	if err != nil {
		return "", fmt.Errorf("invalid manuscript naming parameters: %w", err)
	}

	if strings.TrimSpace(title) != "" {
		params.Title = title
	}

	params.Version = nextVersion

	fileName, err := naming.Manuscript(params)
	if err != nil {
		return "", fmt.Errorf("failed to generate target filename: %w", err)
	}

	return filepath.Join(filepath.Dir(source), fileName), nil
}

func readContent(source, title string) (string, error) {
	data, err := os.ReadFile(source)
	if err != nil {
		return "", fmt.Errorf("failed to read manuscript: %w", err)
	}

	content := string(data)

	if strings.TrimSpace(title) == "" {
		return content, nil
	}

	lines := strings.SplitAfter(content, "\n")

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "# ") {
			newline := ""
			if strings.HasSuffix(line, "\r\n") {
				newline = "\r\n"
			} else if strings.HasSuffix(line, "\n") {
				newline = "\n"
			}

			lines[i] = "# " + strings.TrimSpace(title) + newline
			break
		}
	}

	return strings.Join(lines, ""), nil
}

func saveTarget(target, content string) error {
	file, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return fmt.Errorf("failed to create target file: %w", err)
	}

	if _, err := file.WriteString(content); err != nil {
		closeErr := file.Close()
		writeErr := fmt.Errorf("failed to write target file: %w", err)
		return removeTarget(target, errors.Join(writeErr, closeErr))
	}

	if err := file.Close(); err != nil {
		return removeTarget(target, fmt.Errorf("failed to close target file: %w", err))
	}

	return nil
}

func removeTarget(target string, originalErr error) error {
	if errors.Is(originalErr, os.ErrExist) {
		return originalErr
	}

	if err := os.Remove(target); err != nil && !errors.Is(err, os.ErrNotExist) {
		return errors.Join(originalErr, fmt.Errorf("failed to remove target: %w", err))
	}

	return originalErr
}
