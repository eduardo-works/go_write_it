package main

import (
	"fmt"

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
		rootFolder, err := handleRootFolderCandidate(args)
		if err != nil {
			return err
		}

		title := args[1]

		fmt.Println("Project path:", rootFolder)
		fmt.Println("Project title:", title)
		return nil
	},
}

var projectRenameCmd = &cobra.Command{
	Use:     "rename <rootFolder> <title>",
	Short:   "Rename the project.",
	Example: `writeit project rename ./my-book "My New Book Title"`,
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		rootFolder, err := handleRootFolder(args)
		if err != nil {
			return err
		}

		title := args[1]

		fmt.Println("Project path:", rootFolder)
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
		rootFolder, err := handleRootFolder(args)
		if err != nil {
			return err
		}

		fmt.Println("Project path:", rootFolder)
		return nil
	},
}

var projectAddPartCmd = &cobra.Command{
	Use:     "add-part <rootFolder>",
	Short:   "Add a new part to the project.",
	Example: `writeit project add-part ./my-book`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rootFolder, err := handleRootFolder(args)
		if err != nil {
			return err
		}

		fmt.Println("Project path:", rootFolder)
		return nil
	},
}

var projectNormalizeCmd = &cobra.Command{
	Use:     "normalize <rootFolder>",
	Short:   "Normalize the project file names.",
	Example: `writeit project normalize ./my-book`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rootFolder, err := handleRootFolder(args)
		if err != nil {
			return err
		}

		fmt.Println("Project path:", rootFolder)
		return nil
	},
}

var projectInfoCmd = &cobra.Command{
	Use:     "info <rootFolder>",
	Short:   "Show information about a project.",
	Example: `writeit project info ./my-book`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rootFolder, err := handleRootFolder(args)
		if err != nil {
			return err
		}

		fmt.Println("Project path:", rootFolder)
		return nil
	},
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
