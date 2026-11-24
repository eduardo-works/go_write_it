package main

import (
	"fmt"

	"github.com/eduardo-works/writeit/pkg/template/odt"
	"github.com/eduardo-works/writeit/pkg/terminal"
	"github.com/spf13/cobra"
)

var toOdtCmd = &cobra.Command{
	Use:     "to-odt <filePath>",
	Short:   "Convert a Markdown file to ODT.",
	Example: "writeit to-odt ./manuscript.md",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		destination, err := odt.Convert(args[0])
		if err != nil {
			return fmt.Errorf("file conversion failed: %w", err)
		}

		terminal.Info(fmt.Sprintf("File converted: %s", destination))
		return nil
	},
}
