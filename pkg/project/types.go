package project

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ProjectType string

const (
	Standalone        ProjectType = "standalone"
	Chaptered         ProjectType = "chaptered"
	Parted            ProjectType = "parted"
	ProjectFile       string      = ".writit_project"
	MetaPrompts       string      = "prompts"
	MetaMeta          string      = "meta"
	ContextStandalone string      = "texto"
	ContextChaptered  string      = "capitulos"
	ContextParted     string      = "parte_"
)

type Project struct {
	RootFolder  string            `json:"-"`
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

func (p Project) IsChaptered() bool {
	return IsProjectType(p.ProjectType, Chaptered)
}

func (p Project) IsParted() bool {
	return IsProjectType(p.ProjectType, Parted)
}

func (p Project) Save() error {
	file := filepath.Join(p.RootFolder, ProjectFile)

	data, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		return fmt.Errorf("cannot serialize project: %w", err)
	}

	if err := os.WriteFile(file, data, 0644); err != nil {
		return fmt.Errorf("cannot save project file %s: %w", file, err)
	}

	return nil
}

func (p Project) Contain(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	rel, err := filepath.Rel(p.RootFolder, path)
	if err != nil {
		return "", err
	}

	if rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("the path does not belong to project: %s", path)
	}

	return path, nil
}

func IsValid(value any) bool {
	return IsProjectType(value, Standalone) ||
		IsProjectType(value, Chaptered) ||
		IsProjectType(value, Parted)
}

func LoadProject(rootFolder string) (Project, error) {
	info, err := os.Stat(rootFolder)
	if err != nil {
		return Project{}, err
	}

	if !info.IsDir() {
		rootFolder = filepath.Dir(rootFolder)
	}

	projectFilePath := filepath.Join(rootFolder, ProjectFile)

	data, err := os.ReadFile(projectFilePath)
	if err != nil {
		return Project{}, fmt.Errorf("cannot read project file %s: %w", projectFilePath, err)
	}

	var p Project

	if err := json.Unmarshal(data, &p); err != nil {
		return Project{}, fmt.Errorf("invalid project JSON %s: %w", projectFilePath, err)
	}

	p.RootFolder = filepath.Clean(rootFolder)

	return p, nil
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
