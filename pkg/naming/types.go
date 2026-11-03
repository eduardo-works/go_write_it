package naming

import (
	"fmt"

	"github.com/eduardo-works/writeit/pkg/project"
	"github.com/eduardo-works/writeit/pkg/validator"
)

type NamingParameters struct {
	Project *project.Project
	Part    int
	Index   int
	Version int
	Title   string
}

func (n *NamingParameters) Validate() error {
	if n.Project == nil {
		return fmt.Errorf("no project found: %w", project.ErrNoProject)
	}

	if _, err := validator.ValidateNotEmpty(n.Title); err != nil {
		return fmt.Errorf("title is empty: %w", err)
	}

	if err := n.validateStandalone(); err != nil {
		return err
	}

	if err := n.validateChaptered(); err != nil {
		return err
	}

	if err := n.validateParted(); err != nil {
		return err
	}

	return nil
}

func (n *NamingParameters) validateStandalone() error {
	if !n.Project.IsStandalone() {
		return nil
	}

	if n.Version <= 0 {
		n.Version = 1
	}

	return nil
}

func (n *NamingParameters) validateChaptered() error {
	if !n.Project.IsChaptered() {
		return nil
	}

	if n.Index <= 0 {
		return fmt.Errorf("index must be greater than zero: %w", project.ErrInvalidValue)
	}

	if n.Version <= 0 {
		n.Version = 1
	}

	return nil
}

func (n *NamingParameters) validateParted() error {
	if !n.Project.IsParted() {
		return nil
	}

	if n.Part < 1 || n.Part > n.Project.Parts {
		return fmt.Errorf(
			"part %d is outside project range (1-%d): %w",
			n.Part,
			n.Project.Parts,
			project.ErrInvalidValue,
		)
	}

	if n.Index <= 0 {
		return fmt.Errorf("index must be greater than zero: %w", project.ErrInvalidValue)
	}

	if n.Version <= 0 {
		n.Version = 1
	}

	return nil
}
