package finder

import (
	"path/filepath"
	"sort"

	"github.com/eduardo-works/writeit/pkg/project"
)

func Find(folder string, pattern string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(folder, pattern))
	if err != nil {
		return nil, err
	}

	sort.Sort(sort.Reverse(sort.StringSlice(files)))

	return files, nil
}

func FindMaxVersion(folder string, p *project.Project, part, index int) (string, error) {
	pattern, err := ManuscriptGlob(p, part, index)
	if err != nil {
		return "", err
	}

	files, err := filepath.Glob(filepath.Join(folder, pattern))
	if err != nil {
		return "", err
	}

	sort.Sort(sort.Reverse(sort.StringSlice(files)))
	if len(files) == 0 {
		return "", nil
	}

	return files[0], nil
}
