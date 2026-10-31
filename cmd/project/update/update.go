package update

import (
	"fmt"
	"strings"

	"github.com/eduardo-works/writeit/pkg/project"
	"github.com/eduardo-works/writeit/pkg/validator"
)

func ProjectUpdate(p project.Project, o Options) (project.Project, error) {
	handleAuthor(&p, o)

	if err := handleEmail(&p, o); err != nil {
		return project.Project{}, err
	}

	if err := handlePhone(&p, o); err != nil {
		return project.Project{}, err
	}

	if err := p.Save(); err != nil {
		return project.Project{}, fmt.Errorf("saving project: %w", err)
	}

	return p, nil
}

func handleAuthor(p *project.Project, o Options) {
	author := strings.TrimSpace(o.Author)
	if author == "" {
		return
	}

	p.Author = author
}

func handleEmail(p *project.Project, o Options) error {
	email := strings.TrimSpace(o.Email)
	if email == "" {
		return nil
	}

	email, err := validator.ValidateEmail(email)
	if err != nil {
		return fmt.Errorf("email: %w", err)
	}

	p.Email = email

	return nil
}

func handlePhone(p *project.Project, o Options) error {
	phone := strings.TrimSpace(o.Phone)
	if phone == "" {
		return nil
	}

	phone, err := validator.ValidatePhone(phone)
	if err != nil {
		return fmt.Errorf("phone: %w", err)
	}

	p.Phone = phone

	return nil
}
