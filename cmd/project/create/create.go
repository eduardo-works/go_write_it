package create

import (
	"errors"
	"fmt"
	"os"
	"time"

	cmdproject "github.com/eduardo-works/writeit/cmd/project"
	"github.com/eduardo-works/writeit/pkg/naming"
	"github.com/eduardo-works/writeit/pkg/project"
	"github.com/eduardo-works/writeit/pkg/resolver"
	"github.com/eduardo-works/writeit/pkg/template"
)

func ProjectCreate(rootFolder string, o Options) (p project.Project, err error) {
	slug := naming.Slugify(o.Title)

	projectFolder, err := cmdproject.MakeDir(nil, rootFolder, slug)
	if err != nil {
		return project.Project{}, fmt.Errorf("failed to create project folder: %w", err)
	}

	defer func() {
		if err == nil {
			return
		}

		if removeErr := os.RemoveAll(projectFolder); removeErr != nil {
			err = errors.Join(
				err,
				fmt.Errorf("failed to clean up project folder: %w", removeErr),
			)
		}

		p = project.Project{}
	}()

	p, err = createProjectFile(projectFolder, slug, o)
	if err != nil {
		return p, err
	}

	for _, folder := range []string{project.MetaMeta, project.MetaPrompts} {
		if err = createMetaFolder(&p, folder); err != nil {
			return p, err
		}
	}

	if err = createContextChapteredFolder(&p); err != nil {
		return p, err
	}

	if err = createContextStandaloneFolder(&p); err != nil {
		return p, err
	}

	if err = cmdproject.CreateAllContextPartFolders(&p); err != nil {
		return p, fmt.Errorf("failed to create part folders: %w", err)
	}

	return p, nil
}

func createProjectFile(rootFolder string, slug string, o Options) (project.Project, error) {
	now := time.Now().Format("2006-01-02")

	p := project.Project{
		RootFolder:  rootFolder,
		Title:       o.Title,
		Slug:        slug,
		ProjectType: o.ProjectType,
		Parts:       o.Parts,
		Locale:      "pt_PT",
		Author:      o.Author,
		Email:       o.Email,
		Phone:       o.Phone,
		CreatedAt:   now,
		UpdatedAt:   now,
		Folders:     []string{},
		Files:       make(map[string]string),
	}

	if err := p.Save(); err != nil {
		return project.Project{}, fmt.Errorf("failed to save project: %w", err)
	}

	return p, nil
}

func createContextChapteredFolder(p *project.Project) error {
	if !p.IsChaptered() {
		return nil
	}

	contextFolder, err := resolver.ResolveContext(*p)
	if err != nil {
		return fmt.Errorf("failed to resolve chaptered context: %w", err)
	}

	if _, err := cmdproject.MakeDir(p, contextFolder); err != nil {
		return fmt.Errorf("failed to create chaptered context folder: %w", err)
	}

	cmdproject.AddGitKeep(contextFolder)

	return nil
}

func createContextStandaloneFolder(p *project.Project) error {
	if !p.IsStandalone() {
		return nil
	}

	contextFolder, err := resolver.ResolveContext(*p)
	if err != nil {
		return fmt.Errorf("failed to resolve standalone context: %w", err)
	}

	if _, err := cmdproject.MakeDir(p, contextFolder); err != nil {
		return fmt.Errorf("failed to create standalone context folder: %w", err)
	}

	if _, err := addStandaloneManuscriptFile(p, contextFolder); err != nil {
		return err
	}

	return nil
}

func createMetaFolder(p *project.Project, folderName string) error {
	folder, err := cmdproject.MakeDir(p, p.RootFolder, folderName)
	if err != nil {
		return fmt.Errorf("failed to create folder %s: %w", folderName, err)
	}

	cmdproject.AddGitKeep(folder)

	return nil
}

func addStandaloneManuscriptFile(p *project.Project, contextFolder string) (string, error) {
	content := template.Content(template.StandaloneTemplate, map[string]string{
		"title": p.Title,
	})

	name, err := naming.Manuscript(naming.NamingParameters{
		Project: p,
		Title:   p.Title,
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate standalone manuscript name: %w", err)
	}

	path, err := cmdproject.WriteFile(
		p,
		content,
		contextFolder,
		name,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create standalone manuscript: %w", err)
	}

	return path, nil
}
