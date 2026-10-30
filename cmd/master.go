package main

import (
	"fmt"

	"github.com/eduardo-works/writeit/pkg/resolver"
	"github.com/spf13/cobra"
)

var masterCmd = &cobra.Command{
	Use:   "master",
	Short: "Actions related to master file.",
}

func configureMaster() {
	masterCmd.AddCommand(
		masterAssembleCmd,
		masterDryRunCmd,
		masterCompileCmd,
		masterBuildCmd,
	)

	rootCmd.AddCommand(masterCmd)
}

var masterAssembleCmd = &cobra.Command{
	Use:     "assemble <rootFolder>",
	Short:   "Assemble the project into a single Markdown file.",
	Example: `writeit master assemble ./my-book`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := resolver.ResolveProject(args[0])
		if err != nil {
			return err
		}

		fmt.Println("Root path:", p.RootFolder)
		return nil
	},
}

var masterDryRunCmd = &cobra.Command{
	Use:     "dry-run <rootFolder>",
	Short:   "List the files that will be assembled into the single Markdown file.",
	Example: `writeit master dry-run ./my-book`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := resolver.ResolveProject(args[0])
		if err != nil {
			return err
		}

		fmt.Println("Root path:", p.RootFolder)
		return nil
	},
}

var masterCompileCmd = &cobra.Command{
	Use:     "compile <rootFolder>",
	Short:   "Compile the assembled Markdown into a LibreOffice document.",
	Example: `writeit master compile ./my-book`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := resolver.ResolveProject(args[0])
		if err != nil {
			return err
		}

		fmt.Println("Root path:", p.RootFolder)
		return nil
	},
}

var masterBuildCmd = &cobra.Command{
	Use:     "build <rootFolder>",
	Short:   "Assemble and compile the project.",
	Example: `writeit master build ./my-book`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := resolver.ResolveProject(args[0])
		if err != nil {
			return err
		}

		fmt.Println("Root path:", p.RootFolder)
		return nil
	},
}
