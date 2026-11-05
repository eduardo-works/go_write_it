package compile

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/eduardo-works/writeit/cmd/master"
	"github.com/eduardo-works/writeit/cmd/master/assemble"
	"github.com/eduardo-works/writeit/pkg/project"
)

func MasterCompile(p project.Project) (string, error) {
	masterFile := master.MasterFile(p)
	compiledFile := master.CompiledFile(p)

	cmd, err := pandoc(masterFile, compiledFile)
	if err != nil {
		return "", err
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("pandoc compilation failed: %w: %s", err, strings.TrimSpace(string(output)))
	}

	return compiledFile, nil
}

func MasterBuild(p project.Project) (string, error) {
	_, err := assemble.MasterAssemble(p)
	if err != nil {
		return "", err
	}

	return MasterCompile(p)
}

func pandoc(source, destination string) (*exec.Cmd, error) {
	path, err := exec.LookPath("pandoc")
	if err != nil {
		return nil, fmt.Errorf("pandoc not found: %w", err)
	}

	options := []string{
		"markdown",
		"-auto_identifiers",
		"-hard_line_breaks",
	}

	cmd := exec.Command(
		path,
		source,
		"--from="+strings.Join(options, ""),
		"--output="+destination,
	)

	return cmd, nil
}
