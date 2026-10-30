package naming

import "fmt"

// Pattern: v00_another-slug.md
func StandaloneManuscript(version int, title string) string {
	if version <= 0 {
		version = 1
	}

	return fmt.Sprintf("v%02d_%s.md", version, Slugify(title))
}

// Pattern: i0000_v00_another-slug.md
func ChapteredManuscript(index int, version int, title string) string {
	if index < 0 {
		return ""
	}

	if version <= 0 {
		version = 1
	}

	return fmt.Sprintf("i%04d_v%02d_%s.md", index, version, Slugify(title))
}

// Pattern: p00_i0000_another-slug.md
func PartedPart(part int) string {
	if part <= 0 {
		return ""
	}

	return fmt.Sprintf("p%02d_i0000_parte-%02d.md", part, part)
}

// Pattern: p00_i0000_v00_another-slug.md
func PartedManuscript(part int, index int, version int, title string) string {
	if part <= 0 || index < 0 {
		return ""
	}

	if version <= 0 {
		version = 1
	}

	return fmt.Sprintf("p%02d_i%04d_v%02d_%s.md", part, index, version, Slugify(title))
}
