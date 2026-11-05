package main

import (
	"fmt"

	pkgimage "github.com/eduardo-works/writeit/cmd/image"
	"github.com/spf13/cobra"
)

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Actions related to images.",
}

func configureImage() {
	imageCmd.AddCommand(imageConvertCmd)
	rootCmd.AddCommand(imageCmd)
}

var imageConvertCmd = &cobra.Command{
	Use:     "convert <imagePath> <extension>",
	Short:   "Convert an image to another format.",
	Example: `writeit image convert ./images/photo.webp jpg`,
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
