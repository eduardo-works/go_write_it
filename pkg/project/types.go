package project

type ProjectType string

const (
	Standalone  ProjectType = "standalone"
	Chaptered   ProjectType = "chaptered"
	Parted      ProjectType = "parted"
	ProjectFile string      = ".writit_project"
)

type Project struct {
	Title       string            `json:"title"`
	Slug        string            `json:"slug"`
	ProjectType ProjectType       `json:"project_type"`
	Parts       int               `json:"parts"`
	Locale      string            `json:"locale"`
	Author      string            `json:"author"`
	Email       string            `json:"email"`
	Phone       string            `json:"phone"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
	Folders     []string          `json:"folders"`
	Files       map[string]string `json:"files"`
}

func (p Project) IsStandalone() bool {
	return IsProjectType(p.ProjectType, Standalone)
}

func (p Project) IsChaptered(value any) bool {
	return IsProjectType(p.ProjectType, Chaptered)
}

func (p Project) IsParted(value any) bool {
	return IsProjectType(p.ProjectType, Parted)
}

func IsValid(value any) bool {
	return IsProjectType(value, Standalone) ||
		IsProjectType(value, Chaptered) ||
		IsProjectType(value, Parted)
}

func IsProjectType(value any, expected ProjectType) bool {
	switch v := value.(type) {
	case string:
		return v == string(expected)
	case ProjectType:
		return v == expected
	default:
		return false
	}
}
