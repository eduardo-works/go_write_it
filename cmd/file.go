package main

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/eduardo-works/writeit/cmd/file"
	"github.com/eduardo-works/writeit/pkg/resolver"
	"github.com/eduardo-works/writeit/pkg/terminal"
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
	Example: `writeit file create ./my-book/parte_01 "This Is a Nice Chapter"`,
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		p, err := resolver.ResolveProject(path)
		if err != nil {
			return err
		}

		contextFolder, err := resolver.ResolveContainer(path)
		if err != nil {
			return err
		}

		contextFolder, err = validator.ValidateContext(p, contextFolder)
		if err != nil {
			return err
		}

		title, err := handlePositionalTitle(args, 1)
		if err != nil {
			return err
		}

		createdFile, err := file.FileCreate(p, contextFolder, title)
		if err != nil {
			return fmt.Errorf("file creation failed: %w", err)
		}

		terminal.Info(fmt.Sprintf("File created: %s", createdFile))
		return nil
	},
}

var fileVersionCmd = &cobra.Command{
	Use:     "version <filePath>",
	Short:   "Create a new version of a project file.",
	Example: `writeit file version ./my-book/parte_01/some-chapter.md --title="Change The Title"`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		p, err := resolver.ResolveProject(path)
		if err != nil {
			return err
		}

		title, err := handleTitle(cmd)
		if err != nil {
			return err
		}

		createdFile, sourceFile, err := file.FileVersion(p, path, title)
		if err != nil {
			return fmt.Errorf("file version creation failed: %w", err)
		}

		terminal.Info(fmt.Sprintf("File version created: %s", createdFile))

		output, err := exec.Command("diff", "--", sourceFile, createdFile).CombinedOutput()

		if err == nil {
			terminal.Info(fmt.Sprintf(
				"%s is equal to %s",
				filepath.Base(createdFile),
				filepath.Base(sourceFile),
			))
			return nil
		}

		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			fmt.Print(string(output))
			return nil
		}

		return fmt.Errorf("failed to compare file versions: %w: %s", err, output)
	},
}
