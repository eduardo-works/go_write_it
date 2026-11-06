package compile

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/eduardo-works/writeit/cmd/master"
	"github.com/eduardo-works/writeit/cmd/master/assemble"
	"github.com/eduardo-works/writeit/pkg/pandoc"
	"github.com/eduardo-works/writeit/pkg/project"
	"github.com/eduardo-works/writeit/pkg/template/odt"
)

func MasterCompile(p project.Project) (string, error) {
	masterFile := master.MasterFile(p)
	compiledFile := master.CompiledFile(p)
	template, err := odt.PrepareTemplate(p)
	if err != nil {
		return "", err
	}

	cmd, err := pandoc.Pandoc(masterFile, compiledFile, template)
	if err != nil {
		return "", err
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("pandoc compilation failed: %w: %s", err, strings.TrimSpace(string(output)))
	}

	return compiledFile, nil
}

func MasterBuild(p project.Project) (string, string, error) {
	master, err := assemble.MasterAssemble(p)
	if err != nil {
		return "", "", err
	}

	compiled, err := MasterCompile(p)
	if err != nil {
		return "", "", err
	}

	return master, compiled, nil
}

func ConvertFile(p project.Project, path string) (string, error) {
	target := strings.TrimSuffix(path, filepath.Ext(path)) + ".odt"
	template, err := odt.PrepareTemplate(odt.Clean)
	if err != nil {
		return "", err
	}

	cmd, err := pandoc.Pandoc(path, target, template)
	if err != nil {
		return "", err
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("pandoc compilation failed: %w: %s", err, strings.TrimSpace(string(output)))
	}

	return target, nil
}
