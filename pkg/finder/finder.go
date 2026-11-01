package finder

import (
	"path/filepath"
	"sort"
)

func Find(folder string, pattern string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(folder, pattern))
	if err != nil {
		return nil, err
	}

	sort.Sort(sort.Reverse(sort.StringSlice(files)))

	return files, nil
}
