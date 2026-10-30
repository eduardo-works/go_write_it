package project

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"
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

func (p Project) IsChaptered() bool {
	return IsProjectType(p.ProjectType, Chaptered)
}

func (p Project) IsParted() bool {
	return IsProjectType(p.ProjectType, Parted)
}

func (p Project) IsStandalone() bool {
	return IsProjectType(p.ProjectType, Standalone)
}

func (p Project) IsEmpty() bool {
	return p.RootFolder == ""
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

func (p *Project) AddFolder(folder string) (string, error) {
	relativePath, err := filepath.Rel(p.RootFolder, folder)
	if err != nil {
		return "", err
	}

	if relativePath == ".." || strings.HasPrefix(relativePath, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("folder is outside the project root: %s", folder)
	}

	if slices.Contains(p.Folders, relativePath) {
		return "", fmt.Errorf("folder already exists: %s", relativePath)
	}

	p.Folders = append(p.Folders, relativePath)
	slices.Sort(p.Folders)

	if err := p.Save(); err != nil {
		return "", err
	}

	return relativePath, nil
}

func (p *Project) AddFile(file string) (string, error) {
	relativePath, err := filepath.Rel(p.RootFolder, file)
	if err != nil {
		return "", err
	}

	if filepath.Ext(file) != ".md" {
		return "", nil
	}

	parts := strings.Split(relativePath, string(filepath.Separator))
	if len(parts) < 2 {
		return "", nil
	}

	folder := parts[0]

	validFolder := folder == ContextStandalone ||
		folder == ContextChaptered ||
		regexp.MustCompile(`^parte_[0-9]+$`).MatchString(folder)

	if !validFolder {
		return "", nil
	}

	if regexp.MustCompile(`_parte-[0-9]+\.md$`).MatchString(filepath.Base(file)) {
		return "", nil
	}

	if p.Files == nil {
		p.Files = make(map[string]string)
	}

	if _, exists := p.Files[relativePath]; exists {
		return relativePath, nil
	}

	p.Files[relativePath] = time.Now().Format("2006-01-02")

	if err := p.Save(); err != nil {
		return "", err
	}

	return file, nil
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
