package naming

import (
	"fmt"

	"github.com/eduardo-works/writeit/pkg/project"
)

func Manuscript(n NamingParameters) (string, error) {
	if err := n.Validate(); err != nil {
		return "", err
	}

	switch {
	case n.Project.IsStandalone():
		return StandaloneManuscript(n), nil

	case n.Project.IsChaptered():
		return ChapteredManuscript(n), nil

	case n.Project.IsParted():
		return PartedManuscript(n), nil
	}

	return "", fmt.Errorf("invalid project type: %w", project.ErrInvalidValue)
}

// Pattern: v00_another-slug.md
func StandaloneManuscript(n NamingParameters) string {
	return fmt.Sprintf("v%02d_%s.md", n.Version, Slugify(n.Title))
}

// Pattern: i0000_v00_another-slug.md
func ChapteredManuscript(n NamingParameters) string {
	return fmt.Sprintf(
		"i%04d_v%02d_%s.md",
		n.Index,
		n.Version,
		Slugify(n.Title),
	)
}

// Pattern: p00_i0000_another-slug.md
func PartedPart(part int) string {
	if part <= 0 {
		return ""
	}

	return fmt.Sprintf("p%02d_i0000_parte-%02d.md", part, part)
}

// Pattern: p00_i0000_v00_another-slug.md
func PartedManuscript(n NamingParameters) string {
	return fmt.Sprintf(
		"p%02d_i%04d_v%02d_%s.md",
		n.Part,
		n.Index,
		n.Version,
		Slugify(n.Title),
	)
}
