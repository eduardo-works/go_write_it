package main

import (
	"fmt"

	"github.com/eduardo-works/writeit/pkg/resolver"
	"github.com/eduardo-works/writeit/pkg/validator"
	"github.com/spf13/cobra"
)

var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "Actions related to project files.",
}

func configureFiles() {
	fileCmd.AddCommand(
		fileCreateCmd,
		fileVersionCmd,
	)

	rootCmd.AddCommand(fileCmd)
}

func init() {
	fileVersionCmd.Flags().String("title", "", "Title for the new version")
}

var fileCreateCmd = &cobra.Command{
	Use:     "create <contextFolder> <title>",
	Short:   "Create a new project file.",
	Example: `writeit file create ./my-book/part_01 "This Is a Nice Chapter"`,
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := resolver.ResolveProject(args[0])
		if err != nil {
			return err
		}

		contextFolder, err := resolver.ResolveContext(p)
		if err != nil {
			return err
		}

		title, err := handlePositionalTitle(args, 1)
		if err != nil {
			return err
		}

		fmt.Println("Context path:", contextFolder)
		fmt.Println("Chapter title:", title)
		return nil
	},
}

var fileVersionCmd = &cobra.Command{
	Use:     "version <filePath>",
	Short:   "Create a new version of a project file.",
	Example: `writeit file version ./my-book/part_01/some-chapter.md --title="Change The Title"`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := validator.ValidateManuscript(args[0])
		if err != nil {
			return err
		}

		title, err := handleTitle(cmd)
		if err != nil {
			return err
		}

		fmt.Println("Manuscript path:", path)
		if title != "" {
			fmt.Println("Chapter title:", title)
		}
		return nil
	},
}
