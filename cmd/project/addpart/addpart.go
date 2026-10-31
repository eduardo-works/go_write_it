package addpart

import (
	"fmt"

	cmdproject "github.com/eduardo-works/writeit/cmd/project"
	"github.com/eduardo-works/writeit/pkg/project"
)

func ProjectAddPart(p project.Project) (project.Project, error) {
	if !p.IsParted() {
		return project.Project{}, fmt.Errorf("%w: project is not parted", project.ErrInvalidValue)
	}

	p.Parts++

	_, err := cmdproject.CreateContextPartFolder(&p, p.Parts)
	if err != nil {
		return project.Project{}, err
	}

	if err := p.Save(); err != nil {
		return project.Project{}, fmt.Errorf("failed to save project: %w", err)
	}

	return p, nil
}
