package normalize

import (
	"path/filepath"
	"time"

	"github.com/eduardo-works/writeit/pkg/project"
)

type Normalized struct {
	Source    string
	SourceRel string
	Target    string
	TargetRel string
	Date      string
	Issue     error
}

func (n Normalized) IsChanged() bool {
	return n.Source != n.Target
}

func (n Normalized) Ignore() bool {
	return n.Issue != nil
}

func (n *Normalized) Prepare(p project.Project) {
	n.Issue = nil

	relative, err := filepath.Rel(p.RootFolder, n.Source)
	if err != nil {
		n.Issue = err
		return
	}
	n.SourceRel = relative

	relative, err = filepath.Rel(p.RootFolder, n.Target)
	if err != nil {
		n.Issue = err
		return
	}
	n.TargetRel = relative

	n.Date = findDate(p, n.SourceRel)
}

func findDate(p project.Project, sourceRel string) string {
	date := p.Files[sourceRel]
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	return date
}
