package create

import "github.com/eduardo-works/writeit/pkg/project"

type Options struct {
	Title       string
	ProjectType project.ProjectType
	Parts       int
	Author      string
	Email       string
	Phone       string
}
