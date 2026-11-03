package create

import (
	"fmt"
	"os"
	"time"

	cmdproject "github.com/eduardo-works/writeit/cmd/project"
	"github.com/eduardo-works/writeit/pkg/naming"
	"github.com/eduardo-works/writeit/pkg/project"
	"github.com/eduardo-works/writeit/pkg/resolver"
	"github.com/eduardo-works/writeit/pkg/template"
	"github.com/eduardo-works/writeit/pkg/terminal"
)

func ProjectCreate(rootFolder string, o Options) (project.Project, error) {
	slug := naming.Slugify(o.Title)

	projectFolder, err := cmdproject.MakeDir(nil, rootFolder, slug)
	if err != nil {
		return project.Project{}, err
	}

	p, err := createProjectFile(projectFolder, slug, o)
	if err != nil {
		return project.Project{}, err
	}

	for _, folder := range []string{project.MetaMeta, project.MetaPrompts} {
		createMetaFolder(&p, folder)
	}

	createContextChapteredFolder(&p)
	createContextStandaloneFolder(&p)
	err = cmdproject.CreateAllContextPartFolders(&p)
	if err != nil {
		terminal.Warning(err.Error())
		return p, err
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
		if removeErr := os.RemoveAll(rootFolder); removeErr != nil {
			return project.Project{}, fmt.Errorf(
				"failed to save project: %w; failed to clean up folder: %v",
				err,
				removeErr,
			)
		}

		return project.Project{}, err
	}

	return p, nil
}

func createContextChapteredFolder(p *project.Project) {
	if !p.IsChaptered() {
		return
	}

	contextFolder, err := resolver.ResolveContext(*p)
	if err != nil {
		terminal.Warning(fmt.Sprintf("Failed to resolve the context folder: %v", err))
		return
	}

	if _, err := cmdproject.MakeDir(p, contextFolder); err != nil {
		terminal.Warning(fmt.Sprintf("Failed to create the context folder: %v", err))
		return
	}

	cmdproject.AddGitKeep(contextFolder)
}

func createContextStandaloneFolder(p *project.Project) {
	if !p.IsStandalone() {
		return
	}

	contextFolder, err := resolver.ResolveContext(*p)
	if err != nil {
		terminal.Warning(fmt.Sprintf("Failed to resolve the context folder: %v", err))
		return
	}

	if _, err := cmdproject.MakeDir(p, contextFolder); err != nil {
		terminal.Warning(fmt.Sprintf("Failed to create the context folder: %v", err))
		return
	}

	addStandaloneManuscriptFile(p, contextFolder)
}

func createMetaFolder(p *project.Project, folderName string) {
	folder, err := cmdproject.MakeDir(p, p.RootFolder, folderName)
	if err != nil {
		terminal.Warning(fmt.Sprintf(
			"Failed to create the folder %s: %v",
			folderName,
			err,
		))
		return
	}

	cmdproject.AddGitKeep(folder)
}

func addStandaloneManuscriptFile(p *project.Project, contextFolder string) string {
	content := template.Content(template.StandaloneTemplate, map[string]string{
		"title": p.Title,
	})

	name, err := naming.Manuscript(naming.NamingParameters{
		Project: p,
		Title:   p.Title,
	})
	if err != nil {
		terminal.Warning(fmt.Sprintf(
			"Failed to generate standalone manuscript name: %v",
			err,
		))
		return ""
	}

	path, err := cmdproject.WriteFile(
		p,
		content,
		contextFolder,
		name,
	)
	if err != nil {
		terminal.Warning(fmt.Sprintf(
			"Failed to create standalone manuscript: %v",
			err,
		))
		return ""
	}

	return path
}
