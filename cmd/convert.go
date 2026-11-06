package main

import (
	"fmt"

	pkgimage "github.com/eduardo-works/writeit/cmd/image"
	"github.com/eduardo-works/writeit/cmd/master/compile"
	"github.com/eduardo-works/writeit/pkg/resolver"
	"github.com/eduardo-works/writeit/pkg/validator"
	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Actions related to file conversion.",
}

func configureConvert() {
	convertCmd.AddCommand(
		imageConvertCmd,
		fileConvertCmd,
	)

	rootCmd.AddCommand(convertCmd)
}

var imageConvertCmd = &cobra.Command{
	Use:     "image <imagePath> <extension>",
	Short:   "Convert an image to another format.",
	Example: `writeit convert image ./images/photo.webp jpg`,
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := pkgimage.Convert(args[0], args[1])
		if err != nil {
			return fmt.Errorf("failed to convert image: %w", err)
		}

		fmt.Println("Converted to:", file)
		return nil
	},
}

var fileConvertCmd = &cobra.Command{
	Use:     "file <filePath>",
	Short:   "Convert a Markdown file to ODT.",
	Example: `writeit convert file ./files/something.md`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := resolver.ResolveProject(args[0])
		if err != nil {
			return err
		}

		path, err := validator.ValidateManuscript(args[0])
		if err != nil {
			return err
		}

		file, err := compile.ConvertFile(p, path)
		if err != nil {
			return fmt.Errorf("failed to convert file: %w", err)
		}

		fmt.Println("Converted to:", file)
		return nil
	},
}
