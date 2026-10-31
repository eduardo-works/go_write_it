package tree

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/eduardo-works/writeit/pkg/project"
	"github.com/eduardo-works/writeit/pkg/terminal"
)

type Node struct {
	Name     string
	Children map[string]*Node
	IsFile   bool
}

func Print(p project.Project) {
	root := &Node{
		Name:     filepath.Base(filepath.Clean(p.RootFolder)),
		Children: make(map[string]*Node),
	}

	for _, folder := range p.Folders {
		addPath(root, folder, false)
	}

	for file := range p.Files {
		addPath(root, file, true)
	}

	terminal.WriteLine("{BrightBlue}{Bold}" + root.Name + "/")
	printChildren(root, "")
}

func addPath(root *Node, path string, isFile bool) {
	path = filepath.Clean(path)

	if path == "." || filepath.IsAbs(path) {
		return
	}

	parts := strings.Split(filepath.ToSlash(path), "/")
	current := root

	for i, part := range parts {
		if part == "" || part == "." || part == ".." {
			return
		}

		if current.Children[part] == nil {
			current.Children[part] = &Node{
				Name:     part,
				Children: make(map[string]*Node),
			}
		}

		current = current.Children[part]

		if i == len(parts)-1 && isFile {
			current.IsFile = true
		}
	}
}

func printChildren(node *Node, prefix string) {
	names := sortedChildren(node)

	for i, name := range names {
		child := node.Children[name]
		last := i == len(names)-1

		branch := "├── "
		nextPrefix := prefix + "│   "

		if last {
			branch = "└── "
			nextPrefix = prefix + "    "
		}

		displayName := child.Name
		color := "{White}"

		if !child.IsFile {
			displayName += "/"
			color = "{BrightBlue}{Bold}"
		}

		terminal.WriteLine(
			"{BrightBlack}" + prefix + branch +
				color + displayName,
		)

		if !child.IsFile {
			printChildren(child, nextPrefix)
		}
	}
}

func sortedChildren(node *Node) []string {
	names := make([]string, 0, len(node.Children))

	for name := range node.Children {
		names = append(names, name)
	}

	sort.Slice(names, func(i, j int) bool {
		a := node.Children[names[i]]
		b := node.Children[names[j]]

		if a.IsFile != b.IsFile {
			return !a.IsFile
		}

		return names[i] < names[j]
	})

	return names
}

func PrintLabel(label string, value any) {
	terminal.WriteLine(fmt.Sprintf(
		"{BrightBlue}{Bold}%-12s{Reset} {White}%v{Reset}",
		strings.TrimSpace(label)+":",
		value,
	))
}
