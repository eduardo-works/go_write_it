package main

import (
	"fmt"

	"github.com/eduardo-works/writeit/pkg/resolver"
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
			return err
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

		projectType, parts, err := handleTypeAndPart(cmd)
		if err != nil {
			return err
		}

		fmt.Println("Project path:", rootFolder)
		fmt.Println("Project title:", title)
		fmt.Println("Project type:", projectType)

		if parts > 0 {
			fmt.Println("Parts:", parts)
		}

		fmt.Println("Author:", author)

		if email != "" {
			fmt.Println("Email:", email)
		}

		if phone != "" {
			fmt.Println("Phone:", phone)
		}

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
			return err
		}

		title, err := handlePositionalTitle(args, 1)
		if err != nil {
			return err
		}

		fmt.Println("Project path:", p.RootFolder)
		fmt.Println("New title:", title)

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
			return err
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

		fmt.Println("Project path:", p.RootFolder)

		if author != "" {
			fmt.Println("Author:", author)
		}

		if email != "" {
			fmt.Println("Email:", email)
		}

		if phone != "" {
			fmt.Println("Phone:", phone)
		}

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
			return err
		}

		fmt.Println("Project path:", p.RootFolder)

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
			return err
		}

		fmt.Println("Project path:", p.RootFolder)

		return nil
	},
}

var projectInfoCmd = &cobra.Command{
	Use:     "info <rootFolder>",
	Short:   "Show information about a project.",
	Example: `writeit project info ./my-book`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := resolver.ResolveProject(args[0])
		if err != nil {
			return err
		}

		fmt.Println("Project path:", p.RootFolder)

		return nil
	},
}
