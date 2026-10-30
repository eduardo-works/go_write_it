package project

import (
	"fmt"
	"os"
	"path/filepath"

	pkgProject "github.com/eduardo-works/writeit/pkg/project"
	"github.com/eduardo-works/writeit/pkg/terminal"
)

func MakeDir(p *pkgProject.Project, segments ...string) (string, error) {
	path := filepath.Join(segments...)

	if err := os.Mkdir(path, 0755); err != nil {
		return "", fmt.Errorf("cannot create directory %s: %w", path, err)
	}

	if p != nil {
		if _, err := p.AddFolder(path); err != nil {
			return "", fmt.Errorf("cannot register folder %s: %w", path, err)
		}
	}

	return path, nil
}

func WriteFile(p *pkgProject.Project, content string, segments ...string) (string, error) {
	path := filepath.Join(segments...)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("cannot write file %s: %w", path, err)
	}

	if p != nil {
		if _, err := p.AddFile(path); err != nil {
			return "", fmt.Errorf("cannot register file %s: %w", path, err)
		}
	}

	return path, nil
}

func AddGitKeep(path string) {
	if _, err := WriteFile(nil, "", path, ".gitkeep"); err != nil {
		terminal.Warning(fmt.Sprintf("Failed to create .gitkeep: %v", err))
	}
}
