package assemble

import (
	"fmt"
	"slices"
	"sort"
	"strconv"

	"github.com/eduardo-works/writeit/pkg/finder"
	"github.com/eduardo-works/writeit/pkg/parser"
	"github.com/eduardo-works/writeit/pkg/project"
	"github.com/eduardo-works/writeit/pkg/resolver"
	"github.com/eduardo-works/writeit/pkg/terminal"
	"github.com/eduardo-works/writeit/pkg/tree"
)

func MasterDryRun(p project.Project) error {
	terminal.WriteLine("{BrightCyan}{Bold}FILES TO ASSEMBLE")
	terminal.WriteLine("{BrightCyan}----------------------------------------")
	fmt.Println("")
	files, err := FilterLatestVersions(p)
	if err != nil {
		return err
	}

	tree.FileTree(files)
	return nil
}

func FilterLatestVersions(p project.Project) ([]string, error) {
	files, err := findFiles(p)
	if err != nil {
		return nil, err
	}

	filtered := []string{}
	part := -1
	index := -1
	version := -1
	latest := ""

	for _, file := range files {
		parsed, err := parser.FileParse(p, file)
		if err != nil {
			return nil, err
		}

		if parsed.Part != part || parsed.Index != index {
			if latest != "" {
				filtered = append(filtered, latest)
			}

			part = parsed.Part
			index = parsed.Index
			version = parsed.Version
			latest = file
			continue
		}

		if parsed.Version > version {
			version = parsed.Version
			latest = file
		}
	}

	if latest != "" {
		filtered = append(filtered, latest)
	}

	slices.Reverse(filtered)

	return filtered, nil
}

func findFiles(p project.Project) ([]string, error) {
	if p.IsStandalone() || p.IsChaptered() {
		contextFolder, err := resolver.ResolveContext(p)
		if err != nil {
			return nil, err
		}

		return finder.Find(contextFolder, "*.md")
	}

	files := []string{}

	for part := range p.Parts {
		contextFolder, err := resolver.ResolveContext(p, strconv.Itoa(part+1))
		if err != nil {
			return nil, err
		}

		found, err := finder.Find(contextFolder, "*.md")
		if err != nil {
			return nil, err
		}

		files = append(files, found...)
	}

	sort.Sort(sort.Reverse(sort.StringSlice(files)))
	return files, nil
}
