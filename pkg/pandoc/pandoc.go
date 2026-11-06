package pandoc

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Pandoc(source, target, template string) (*exec.Cmd, error) {
	path, err := exec.LookPath("pandoc")
	if err != nil {
		return nil, fmt.Errorf("pandoc not found: %w", err)
	}

	if err := os.Remove(target); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	options := []string{
		"markdown",
		"-auto_identifiers",
		"-hard_line_breaks",
	}

	args := []string{
		source,
		"--from=" + strings.Join(options, ""),
		"--metadata=lang=pt-PT",
	}

	if template != "" {
		args = append(args, "--reference-doc="+template)
	}

	args = append(args, "--output="+target)

	cmd := exec.Command(path, args...)

	return cmd, nil
}
