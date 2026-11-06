package main

import (
	"fmt"
	"os"

	"github.com/eduardo-works/writeit/pkg/terminal"
	"github.com/spf13/cobra"
)

var debug bool

var rootCmd = &cobra.Command{
	Use:           "writeit",
	Short:         "A project manager for writers.",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(
		&debug,
		"debug",
		"d",
		false,
		"Show detailed errors",
	)
}

func main() {
	configureProject()
	configureFiles()
	configureMaster()
	configureConvert()
	if err := rootCmd.Execute(); err != nil {
		if debug {
			terminal.Error(fmt.Sprintf("%+v", err))
			os.Exit(1)
		}

		terminal.Error(err.Error())
		os.Exit(1)
	}
}
