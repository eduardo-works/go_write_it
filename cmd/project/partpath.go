package project

import (
	"fmt"

	"github.com/eduardo-works/writeit/pkg/naming"
	"github.com/eduardo-works/writeit/pkg/project"
	"github.com/eduardo-works/writeit/pkg/resolver"
	"github.com/eduardo-works/writeit/pkg/template"
)

func CreateAllContextPartFolders(p *project.Project) error {
	if !p.IsParted() {
		return nil
	}

	for part := 1; part <= p.Parts; part++ {
		if _, err := CreateContextPartFolder(p, part); err != nil {
			return err
		}
	}

	return nil
}

func CreateContextPartFolder(p *project.Project, part int) (string, error) {
	if !p.IsParted() {
		return "", nil
	}

	contextFolder, err := resolver.ResolveContext(*p, fmt.Sprint(part))
	if err != nil {
		return "", fmt.Errorf("failed to resolve the context folder for part %d: %w", part, err)
	}

	contextFolder, err = MakeDir(p, contextFolder)
	if err != nil {
		return "", fmt.Errorf("failed to create the context folder for part %d: %w", part, err)
	}

	file, err := AddPartFile(p, contextFolder, part)
	if err != nil {
		return "", err
	}

	return file, nil
}

func AddPartFile(p *project.Project, contextFolder string, part int) (string, error) {
	title := fmt.Sprintf("Parte %s", naming.PortugueseNumbersToWords(part))

	content := template.Content(template.PartedTemplate, map[string]string{
		"title": title,
	})

	path, err := WriteFile(
		p,
		content,
		contextFolder,
		naming.PartedPart(part),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create part file for part %d: %w", part, err)
	}

	return path, nil
}
