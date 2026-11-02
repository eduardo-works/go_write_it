package main

import (
	"fmt"
	"maps"
	"slices"

	"github.com/eduardo-works/writeit/cmd/project/addpart"
	"github.com/eduardo-works/writeit/cmd/project/create"
	"github.com/eduardo-works/writeit/cmd/project/info"
	"github.com/eduardo-works/writeit/cmd/project/rename"
	"github.com/eduardo-works/writeit/cmd/project/update"
	"github.com/eduardo-works/writeit/pkg/normalizer"
	"github.com/eduardo-works/writeit/pkg/resolver"
	"github.com/eduardo-works/writeit/pkg/terminal"
	"github.com/eduardo-works/writeit/pkg/validator"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Actions related to the project itself.",
}

func configureProject() {
	projectCmd.AddCommand(
		projectCreateCmd,
		projectRenameCmd,
		projectUpdateCmd,
		projectAddPartCmd,
		projectNormalizeCmd,
		projectDryRunCmd,
		projectInfoCmd,
	)

	rootCmd.AddCommand(projectCmd)
}

func init() {
	projectCreateCmd.Flags().StringP("author", "a", "", "Author name")
	projectCreateCmd.Flags().StringP("email", "e", "", "Author email")
	projectCreateCmd.Flags().StringP("phone", "p", "", "Author phone")
	projectCreateCmd.Flags().StringP("type", "t", "standalone", "Project type: standalone, chaptered or parted")
	projectCreateCmd.Flags().Int("parts", 0, "Number of project parts")

	projectUpdateCmd.Flags().StringP("author", "a", "", "Author name")
	projectUpdateCmd.Flags().StringP("email", "e", "", "Author email")
	projectUpdateCmd.Flags().StringP("phone", "p", "", "Author phone")
}

var projectCreateCmd = &cobra.Command{
	Use:   "create <rootFolder> <title>",
	Short: "Create a new project.",
	Example: `writeit project create ./my-book "My Book Title" ` +
		`--author="John Doe" ` +
		`--email="john.doe@some-email.com" ` +
		`--phone="+351 222 333 444" ` +
		`--type=parted --parts=3`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		rootFolder, err := validator.ValidateRootFolderCandidate(args[0])
		if err != nil {
			return fmt.Errorf("rootFolder: %w", err)
		}

		title, err := handlePositionalTitle(args, 1)
		if err != nil {
			return err
		}

		author, err := handleAuthor(cmd)
		if err != nil {
			return err
		}

		author, err = validator.ValidateNotEmpty(author)
		if err != nil {
			return fmt.Errorf("author: %w", err)
		}

		email, err := handleEmail(cmd)
		if err != nil {
			return err
		}

		phone, err := handlePhone(cmd)
		if err != nil {
			return err
		}

		projectType, parts, err := handleTypeAndPart(cmd)
		if err != nil {
			return err
		}

		p, err := create.ProjectCreate(rootFolder, create.Options{
			Title:       title,
			Author:      author,
			Email:       email,
			Phone:       phone,
			ProjectType: projectType,
			Parts:       parts,
		})
		if err != nil {
			return fmt.Errorf("creating project: %w", err)
		}

		terminal.Info(fmt.Sprintf("Project created: %s", p.RootFolder))

		return nil
	},
}

var projectRenameCmd = &cobra.Command{
	Use:     "rename <rootFolder> <title>",
	Short:   "Rename the project.",
	Example: `writeit project rename ./my-book "My New Book Title"`,
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := resolver.ResolveProject(args[0])
		if err != nil {
			return fmt.Errorf("rootFolder: %w", err)
		}

		title, err := handlePositionalTitle(args, 1)
		if err != nil {
			return err
		}

		p, err = rename.ProjectRename(p, title)
		if err != nil {
			return fmt.Errorf("renaming project: %w", err)
		}

		terminal.Info(fmt.Sprintf("Project renamed: %s", p.RootFolder))

		return nil
	},
}

var projectUpdateCmd = &cobra.Command{
	Use:   "update <rootFolder>",
	Short: "Update the project metadata.",
	Example: `writeit project update ./my-book ` +
		`--author="John Doe" ` +
		`--email="john.doe@some-email.com" ` +
		`--phone="+351 222 333 444"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := resolver.ResolveProject(args[0])
		if err != nil {
			return fmt.Errorf("rootFolder: %w", err)
		}

		author, err := handleAuthor(cmd)
		if err != nil {
			return err
		}

		email, err := handleEmail(cmd)
		if err != nil {
			return err
		}

		phone, err := handlePhone(cmd)
		if err != nil {
			return err
		}

		p, err = update.ProjectUpdate(p, update.Options{
			Author: author,
			Email:  email,
			Phone:  phone,
		})
		if err != nil {
			return fmt.Errorf("updating project: %w", err)
		}

		terminal.Info(fmt.Sprintf("Project updated: %s", p.RootFolder))

		return nil
	},
}

var projectAddPartCmd = &cobra.Command{
	Use:     "add-part <rootFolder>",
	Short:   "Add a new part to the project.",
	Example: `writeit project add-part ./my-book`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := resolver.ResolveProject(args[0])
		if err != nil {
			return fmt.Errorf("rootFolder: %w", err)
		}

		p, err = addpart.ProjectAddPart(p)
		if err != nil {
			return fmt.Errorf("adding project part: %w", err)
		}

		terminal.Info(fmt.Sprintf("Part added to project: %s", p.RootFolder))

		return nil
	},
}

var projectNormalizeCmd = &cobra.Command{
	Use:     "normalize <rootFolder>",
	Short:   "Normalize the project file names.",
	Example: `writeit project normalize ./my-book`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := resolver.ResolveProject(args[0])
		if err != nil {
			return fmt.Errorf("rootFolder: %w", err)
		}

		issues, count, err := normalizer.Normalize(&p)

		if count > 0 {
			terminal.Info(fmt.Sprintf("%d files normalized", count))
		} else if err == nil {
			terminal.Info("No files were normalized")
		}

		printNormalizationIssues(issues)

		if err != nil {
			return fmt.Errorf("normalizing project: %w", err)
		}

		return nil
	},
}

var projectDryRunCmd = &cobra.Command{
	Use:     "dry-run <rootFolder>",
	Short:   "List the project files that would be normalized.",
	Example: `writeit project dry-run ./my-book`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := resolver.ResolveProject(args[0])
		if err != nil {
			return fmt.Errorf("rootFolder: %w", err)
		}

		files, issues, count, err := normalizer.DryRun(&p)

		if count > 0 {
			terminal.Info(fmt.Sprintf("%d files would be normalized", count))
		} else if err == nil {
			terminal.Info("No files need to be normalized")
		}

		if len(files) > 0 {
			terminal.WriteLine("{BrightCyan}{Bold}Files to be normalized")
			terminal.WriteLine("{BrightCyan}----------------------------------------")
			terminal.WriteLine("")

			for _, file := range slices.Sorted(maps.Keys(files)) {
				terminal.WriteLine(fmt.Sprintf("   %s: %s", files[file], file))
			}
		}

		printNormalizationIssues(issues)

		if err != nil {
			return fmt.Errorf("project dry run: %w", err)
		}

		return nil
	},
}

func printNormalizationIssues(issues []string) {
	if len(issues) == 0 {
		return
	}

	slices.Sort(issues)

	terminal.WriteLine("{BrightCyan}{Bold}Normalization issues")
	terminal.WriteLine("{BrightCyan}----------------------------------------")
	terminal.WriteLine("")

	for index, issue := range issues {
		terminal.WriteLine(fmt.Sprintf("   #%04d: %s", index+1, issue))
	}
}

var projectInfoCmd = &cobra.Command{
	Use:     "info <rootFolder>",
	Short:   "Show information about a project.",
	Example: `writeit project info ./my-book`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := resolver.ResolveProject(args[0])
		if err != nil {
			return fmt.Errorf("rootFolder: %w", err)
		}

		info.ProjectInfo(p)

		return nil
	},
}
