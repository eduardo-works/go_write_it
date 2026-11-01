package path

import (
	"path/filepath"
	"time"

	"github.com/eduardo-works/writeit/pkg/project"
)

type NormalizedPath struct {
	Project   *project.Project
	Source    string
	SourceRel string
	Target    string
	TargetRel string
	Date      string
	Issue     error
}

func (n NormalizedPath) IsChanged() bool {
	return n.Source != n.Target
}

func (n NormalizedPath) Ignore() bool {
	return n.Issue != nil
}

func (n *NormalizedPath) SetDate() string {
	date := n.Project.Files[n.SourceRel]
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	return date
}

func (n *NormalizedPath) Prepare() {
	n.Issue = nil

	relative, err := filepath.Rel(n.Project.RootFolder, n.Source)
	if err != nil {
		n.Issue = err
		return
	}
	n.SourceRel = relative

	relative, err = filepath.Rel(n.Project.RootFolder, n.Target)
	if err != nil {
		n.Issue = err
		return
	}
	n.TargetRel = relative

	n.Date = n.SetDate()
}
