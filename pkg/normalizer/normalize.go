package normalizer

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"

	"github.com/eduardo-works/writeit/pkg/finder"
	normalizerpath "github.com/eduardo-works/writeit/pkg/normalizer/path"
	"github.com/eduardo-works/writeit/pkg/project"
	"github.com/eduardo-works/writeit/pkg/resolver"
)

func Normalize(p *project.Project) ([]string, int, error) {
	normalized := make(map[string]string)
	tnFiles := []normalizerpath.NormalizedPath{}
	issues := []string{}

	if err := fillToNormalize(p, &tnFiles); err != nil {
		return issues, 0, err
	}

	count := normalizeFiles(tnFiles, normalized, &issues)

	p.Files = normalized

	if err := p.Save(); err != nil {
		return issues, count, err
	}

	slices.Sort(issues)
	return issues, count, nil
}

func DryRun(p *project.Project) (map[string]string, []string, int, error) {
	normalized := make(map[string]string)
	tnFiles := []normalizerpath.NormalizedPath{}
	issues := []string{}
	count := 0

	if err := fillToNormalize(p, &tnFiles); err != nil {
		return normalized, issues, 0, err
	}

	for _, item := range tnFiles {
		prepareItem(&item)

		if item.Ignore() {
			normalized[item.SourceRel] = item.Date
			addIssue(&issues, item, item.Issue)
			continue
		}

		if item.IsChanged() {
			count = count + 1
			normalized[item.TargetRel] = item.Date
		}
	}
	slices.Sort(issues)
	return normalized, issues, count, nil
}

func normalizeFiles(
	tnFiles []normalizerpath.NormalizedPath,
	normalized map[string]string,
	issues *[]string,
) int {
	count := 0
	for _, item := range tnFiles {
		prepareItem(&item)

		if item.Ignore() {
			normalized[item.SourceRel] = item.Date
			addIssue(issues, item, item.Issue)
			continue
		}

		if !item.IsChanged() {
			normalized[item.SourceRel] = item.Date
			continue
		}

		count = count + 1
		rename(item, normalized, issues)
	}
	return count
}

func prepareItem(item *normalizerpath.NormalizedPath) {
	if item.SourceRel == "" {
		relative, err := filepath.Rel(item.Project.RootFolder, item.Source)
		if err != nil {
			item.Issue = err
			return
		}

		item.SourceRel = relative
	}

	if item.Date == "" {
		item.Date = item.SetDate()
	}
}

func addIssue(
	issues *[]string,
	item normalizerpath.NormalizedPath,
	err error,
) {
	*issues = append(*issues, fmt.Sprintf(
		"[%s] %v",
		filepath.Base(item.Source),
		err,
	))
}

func rename(
	item normalizerpath.NormalizedPath,
	normalized map[string]string,
	issues *[]string,
) {
	err := os.Rename(item.Source, item.Target)
	if err != nil {
		normalized[item.SourceRel] = item.Date
		addIssue(issues, item, fmt.Errorf(
			"failed to rename to %s: %w",
			filepath.Base(item.Target),
			err,
		))
		return
	}

	normalized[item.TargetRel] = item.Date
}

func fillToNormalize(
	p *project.Project,
	tnFiles *[]normalizerpath.NormalizedPath,
) error {
	if !p.IsParted() {
		folder, err := resolver.ResolveContext(*p)
		if err != nil {
			return err
		}

		return fillForContext(p, folder, tnFiles)
	}

	for part := 1; part <= p.Parts; part++ {
		folder, err := resolver.ResolveContext(*p, strconv.Itoa(part))
		if err != nil {
			return err
		}

		if err := fillForContext(p, folder, tnFiles); err != nil {
			return err
		}
	}

	return nil
}

func fillForContext(
	p *project.Project,
	folder string,
	tnFiles *[]normalizerpath.NormalizedPath,
) error {
	foundFiles, err := finder.Find(folder, "*.md")
	if err != nil {
		return err
	}

	for _, foundFile := range foundFiles {
		*tnFiles = append(
			*tnFiles,
			normalizerpath.Normalize(p, foundFile),
		)
	}

	return nil
}
