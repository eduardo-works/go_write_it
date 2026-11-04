package tree

import (
	"path/filepath"
	"sort"

	"github.com/eduardo-works/writeit/pkg/terminal"
)

func FileTree(files []string) {
	folders := make(map[string]*Node)

	for _, file := range files {
		folder := filepath.Base(filepath.Dir(file))
		name := filepath.Base(file)

		if folders[folder] == nil {
			folders[folder] = &Node{
				Name:     folder,
				Children: make(map[string]*Node),
			}
		}

		folders[folder].Children[name] = &Node{
			Name:   name,
			IsFile: true,
		}
	}

	names := make([]string, 0, len(folders))
	for name := range folders {
		names = append(names, name)
	}

	sort.Strings(names)

	for i, name := range names {
		if i > 0 {
			terminal.WriteLine("")
		}

		root := folders[name]

		terminal.WriteLine("{BrightBlue}{Bold}" + root.Name + "/")
		printChildren(root, "")
	}
}
