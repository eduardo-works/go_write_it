package main

import (
	"github.com/eduardo-works/writeit/pkg/validator"
)

func handleRootFolderCandidate(args []string) (string, error) {
	return validator.ValidateRootFolderCandidate(args[0])
}

func handleRootFolder(args []string) (string, error) {
	return validator.ValidateRootFolder(args[0])
}
