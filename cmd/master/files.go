package master

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/eduardo-works/writeit/pkg/project"
)

// Pattern: master.md
func MasterFile(p project.Project) string {
	return filepath.Join(p.RootFolder, "master.md")
}

// Pattern: some-title_2006-01-02.odt.md
func CompiledFile(p project.Project) string {
	return filepath.Join(
		p.RootFolder,
		project.MetaMeta,
		fmt.Sprintf("%s_%s.odt", p.Slug, time.Now().Format("2006-01-02")),
	)
}
