package file

import (
	"fmt"

	cmdproject "github.com/eduardo-works/writeit/cmd/project"
	"github.com/eduardo-works/writeit/pkg/naming"
	"github.com/eduardo-works/writeit/pkg/project"
	"github.com/eduardo-works/writeit/pkg/resolver"
	"github.com/eduardo-works/writeit/pkg/template"
)

func FileCreate(p project.Project, contextFolder, title string) (string, error) {
	fileName, err := resolveName(p, contextFolder, title)
	if err != nil {
		return "", err
	}

	return saveFile(p, contextFolder, title, fileName)
}

func saveFile(p project.Project, contextFolder, title, fileName string) (string, error) {
	content := template.Content(template.StandaloneTemplate, map[string]string{
		"title": title,
	})

	path, err := cmdproject.WriteFile(
		&p,
		content,
		contextFolder,
		fileName,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create manuscript: %w", err)
	}

	return path, nil
}

func resolveName(p project.Project, contextFolder, title string) (string, error) {
	index, err := resolver.ResolveNextIndex(p, contextFolder)
	if err != nil {
		return "", fmt.Errorf("failed to resolve next index: %w", err)
	}

	part, err := resolver.ResolvePart(p, contextFolder)
	if err != nil {
		return "", fmt.Errorf("failed to resolve part: %w", err)
	}

	params := naming.NamingParameters{
		Project: &p,
		Part:    part,
		Index:   index,
		Version: 1,
		Title:   title,
	}

	fileName, err := naming.Manuscript(params)
	if err != nil {
		return "", fmt.Errorf("failed to generate file name: %w", err)
	}

	return fileName, nil
}
