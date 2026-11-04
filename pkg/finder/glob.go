package finder

import (
	"fmt"

	"github.com/eduardo-works/writeit/pkg/project"
)

func ManuscriptGlob(p *project.Project, part, index int) (string, error) {
	if p == nil {
		return "", fmt.Errorf("no project found: %w", project.ErrNoProject)
	}

	switch {
	case p.IsStandalone():
		return StandaloneManuscriptGlob(), nil

	case p.IsChaptered():
		if index <= 0 {
			return "", fmt.Errorf("invalid index: %w", project.ErrInvalidValue)
		}
		return ChapteredManuscriptGlob(index), nil

	case p.IsParted():
		if part < 1 || part > p.Parts {
			return "", fmt.Errorf("invalid part: %w", project.ErrInvalidValue)
		}
		if index <= 0 {
			return "", fmt.Errorf("invalid index: %w", project.ErrInvalidValue)
		}
		return PartedManuscriptGlob(part, index), nil
	}

	return "", fmt.Errorf("invalid project type: %w", project.ErrInvalidValue)
}

// Pattern: v*_*.md
func StandaloneManuscriptGlob() string {
	return "v*_*.md"
}

// Pattern: i0000_v*_*.md
func ChapteredManuscriptGlob(index int) string {
	return fmt.Sprintf("i%04d_v*_*.md", index)
}

// Pattern: p00_i0000_parte-00.md
func PartedPartGlob(part int) string {
	if part <= 0 {
		return ""
	}

	return fmt.Sprintf("p%02d_i0000_parte-%02d.md", part, part)
}

// Pattern: p00_i0000_v*_*.md
func PartedManuscriptGlob(part, index int) string {
	return fmt.Sprintf("p%02d_i%04d_v*_*.md", part, index)
}
