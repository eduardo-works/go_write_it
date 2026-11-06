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
		} else {
			terminal.Error(err.Error())
		}

		cmd, _, findErr := rootCmd.Find(os.Args[1:])
		if findErr == nil && cmd != nil {
			fmt.Fprintln(os.Stderr)
			cmd.SetOut(os.Stderr)
			_ = cmd.Usage()
		}

		os.Exit(1)
	}
}
