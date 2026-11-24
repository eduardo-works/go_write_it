package odt

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/eduardo-works/writeit/pkg/pandoc"
	"github.com/eduardo-works/writeit/pkg/project"
	"github.com/eduardo-works/writeit/pkg/resolver"
)

func Convert(path string) (string, error) {
	path, err := resolver.ExpandPath(path)
	if err != nil {
		return "", err
	}

	ext := filepath.Ext(path)
	if !strings.EqualFold(ext, ".md") {
		return "", fmt.Errorf("only markdown files allowed: %w", project.ErrNotManuscriptFile)
	}

	name := strings.TrimSuffix(path, ext)
	date := time.Now().Format("20060102")
	destination := fmt.Sprintf("%s_%s.odt", name, date)

	template, err := PrepareTemplate(StandaloneODT)
	if err != nil {
		return "", err
	}

	cmd, err := pandoc.Pandoc(path, destination, template)
	if err != nil {
		return "", err
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("pandoc compilation failed: %w: %s", err, strings.TrimSpace(string(output)))
	}

	return destination, nil
}
