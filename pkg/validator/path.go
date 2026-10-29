package validator

import (
	"github.com/eduardo-works/writeit/pkg/resolver"
)

func ValidateExists(path string) (string, error) {
	return resolver.ResolvePath(path)
}

func ValidateFile(path string) (string, error) {
	return resolver.ResolveFile(path)
}

func ValidateDir(path string) (string, error) {
	return resolver.ResolveDir(path)
}
