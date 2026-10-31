package rename

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/eduardo-works/writeit/pkg/naming"
	"github.com/eduardo-works/writeit/pkg/project"
	"github.com/eduardo-works/writeit/pkg/validator"
)

const projectFile = ".writeit_project"

func ProjectRename(p project.Project, title string) (project.Project, error) {

	title, err := validator.ValidateNotEmpty(title)
	if err != nil {
		return project.Project{}, err
	}

	if p.Title == title {
		return p, nil
	}

	slug := naming.Slugify(title)
	if slug == "" {
		return project.Project{}, fmt.Errorf(
			"%w: title produces an empty slug",
			project.ErrInvalidValue,
		)
	}

	sourceFolder := p.RootFolder
	destinationFolder := filepath.Join(filepath.Dir(sourceFolder), slug)

	needsRename := sourceFolder != destinationFolder

	if needsRename {
		if err := validateDestination(destinationFolder); err != nil {
			return project.Project{}, err
		}
	}

	stagingFolder, err := os.MkdirTemp(sourceFolder, ".rename-*")
	if err != nil {
		return project.Project{}, err
	}

	defer func() {
		os.RemoveAll(stagingFolder)
	}()

	updated := p
	updated.Title = title
	updated.Slug = slug
	updated.RootFolder = stagingFolder

	if err := updated.Save(); err != nil {
		return project.Project{}, fmt.Errorf(
			"preparing project: %w", err,
		)
	}

	stagedFile := filepath.Join(stagingFolder, projectFile)

	if _, err := os.Stat(stagedFile); err != nil {
		return project.Project{}, fmt.Errorf(
			"prepared project file: %w", err,
		)
	}

	if needsRename {
		if err := os.Rename(sourceFolder, destinationFolder); err != nil {
			return project.Project{}, fmt.Errorf(
				"renaming project folder: %w", err,
			)
		}

		stagingFolder = filepath.Join(
			destinationFolder,
			filepath.Base(stagingFolder),
		)
		stagedFile = filepath.Join(stagingFolder, projectFile)
	}

	destinationFile := filepath.Join(destinationFolder, projectFile)

	if err := os.Rename(stagedFile, destinationFile); err != nil {
		if needsRename {
			rollbackErr := os.Rename(destinationFolder, sourceFolder)

			if rollbackErr != nil {
				return project.Project{}, fmt.Errorf(
					"replacing project file: %v; restoring folder: %w",
					err, rollbackErr,
				)
			}

			stagingFolder = filepath.Join(
				sourceFolder,
				filepath.Base(stagingFolder),
			)
		}

		return project.Project{}, fmt.Errorf(
			"replacing project file: %w", err,
		)
	}

	updated.RootFolder = destinationFolder

	return updated, nil
}

func validateDestination(folder string) error {
	_, err := os.Lstat(folder)

	if err == nil {
		return fmt.Errorf(
			"%w: destination already exists: %s",
			project.ErrInvalidValue,
			folder,
		)
	}

	if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}
