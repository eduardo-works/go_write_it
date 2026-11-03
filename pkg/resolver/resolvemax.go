package resolver

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/eduardo-works/writeit/pkg/finder"
	"github.com/eduardo-works/writeit/pkg/project"
)

// ResolveMaxIndex returns the highest manuscript index found in the given context.
// Returns -1 for standalone projects, which do not use indexes.
// Returns 0 when no indexed manuscripts are found.
// Existing indexes are returned as-is, even if they are not multiples of 10.
func ResolveMaxIndex(p project.Project, path string) (int, error) {
	if p.IsStandalone() {
		return -1, nil
	}

	folder, err := validateDirForSearch(p, path)
	if err != nil {
		return 0, err
	}

	files, err := finder.Find(folder, "*.md")
	if err != nil {
		return 0, err
	}

	return iterateFilesForMaxValue(files, `(?:^|_)i([0-9]+)(?:_|$)`), nil
}

// ResolveNextIndex returns the next available manuscript index.
// Indexes start at 10 and advance to the next multiple of 10.
// For example, 0 becomes 10, 17 becomes 20, and 20 becomes 30.
// Returns -1 for standalone projects, which do not use indexes.
func ResolveNextIndex(p project.Project, path string) (int, error) {
	maxIndex, err := ResolveMaxIndex(p, path)
	if err != nil {
		return -1, err
	}

	if maxIndex == -1 {
		return -1, nil
	}

	return (maxIndex/10 + 1) * 10, nil
}

// ResolveMaxVersion returns the highest version of a manuscript.
// When given a file, it searches for versions of the same manuscript,
// matching its index and part when present, regardless of the slug.
// When given a directory, it returns 0, indicating a new manuscript.
// Returns 0 when no matching versions are found.
func ResolveMaxVersion(p project.Project, path string) (int, error) {
	resolved, err := ResolvePath(path)
	if err != nil {
		return 0, err
	}

	info, err := os.Stat(resolved)
	if err != nil {
		return 0, err
	}

	folder, err := validateDirForSearch(p, resolved)
	if err != nil {
		return 0, err
	}

	if info.IsDir() {
		return 0, nil
	}

	regex, err := versionRegex(resolved)
	if err != nil {
		return 0, err
	}

	files, err := finder.Find(folder, "*.md")
	if err != nil {
		return 0, err
	}

	return iterateFilesForMaxValue(files, regex), nil
}

// ResolveNextVersion returns the next version number for a manuscript.
// Versions start at 1 and increase by 1.
// When given a directory, it returns 1 for a new manuscript.
// When given a file, it returns the highest matching version plus 1.
// The manuscript slug may change between versions.
func ResolveNextVersion(p project.Project, path string) (int, error) {
	maxVersion, err := ResolveMaxVersion(p, path)
	if err != nil {
		return -1, err
	}

	if maxVersion == -1 {
		return -1, nil
	}

	return maxVersion + 1, nil
}

func validateDirForSearch(p project.Project, path string) (string, error) {
	folder, err := ResolveContainer(path)
	if err != nil {
		return "", err
	}

	folder, err = p.Contain(folder)
	if err != nil {
		return "", err
	}

	if folder == p.RootFolder {
		return "", fmt.Errorf("cannot search for indexes or versions in the root folder: %w", project.ErrInvalidValue)
	}

	return folder, nil

}

func iterateFilesForMaxValue(files []string, regex string) int {
	maxValue := 0
	re := regexp.MustCompile(regex)

	for _, file := range files {
		matches := re.FindStringSubmatch(filepath.Base(file))
		if matches == nil {
			continue
		}

		selected, err := strconv.Atoi(matches[1])
		if err != nil {
			continue
		}

		if selected > maxValue {
			maxValue = selected
		}
	}

	return maxValue
}

func versionRegex(filename string) (string, error) {
	re := regexp.MustCompile(`^(?:(p[0-9]+_i[0-9]+|i[0-9]+)_)?v[0-9]+_.*\.md$`)

	matches := re.FindStringSubmatch(filepath.Base(filename))
	if matches == nil {
		return "", fmt.Errorf("invalid filename: %s", filename)
	}

	if matches[1] == "" {
		return `^v([0-9]+)_.*\.md$`, nil
	}

	return fmt.Sprintf(
		`^%s_v([0-9]+)_.*\.md$`,
		regexp.QuoteMeta(matches[1]),
	), nil
}
