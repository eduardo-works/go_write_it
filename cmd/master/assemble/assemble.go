package assemble

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/eduardo-works/writeit/cmd/master"
	"github.com/eduardo-works/writeit/pkg/project"
	"github.com/eduardo-works/writeit/pkg/template"
	"github.com/eduardo-works/writeit/pkg/terminal"
)

func MasterAssemble(p project.Project) (string, error) {
	content := assemble(p)
	if content == "" {
		return "", nil
	}

	file := master.MasterFile(p)

	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		terminal.Error(err.Error())
		return "", err
	}

	return file, nil
}

func assemble(p project.Project) string {
	files, err := FilterLatestVersions(p)
	if err != nil {
		terminal.Error(err.Error())
		return ""
	}

	content := strings.TrimSpace(generateMasterContent(p, files))
	if content == "" {
		terminal.Warning("The project has no content.")
		return ""
	}

	count := countContent(content)

	var authorInfo strings.Builder

	if p.Author != "" {
		fmt.Fprintf(&authorInfo, "Autor: %s%s", p.Author, template.NewLine)
	}

	if p.Email != "" {
		fmt.Fprintf(&authorInfo, "Email: %s%s", p.Email, template.NewLine)
	}

	if p.Phone != "" {
		fmt.Fprintf(&authorInfo, "Telemóvel: %s%s", p.Phone, template.NewLine)
	}

	fmt.Fprintf(&authorInfo, "Data: %s%s", time.Now().Format("2006-01-02"), template.NewLine)
	fmt.Fprintf(&authorInfo, "Palavras: %s%s", count.Words, template.NewLine)
	fmt.Fprintf(&authorInfo, "Caracteres: %s%s", count.CharactersWithSpace, template.NewLine)

	master := template.Content(template.MasterTemplate, map[string]string{
		"title":       p.Title,
		"author_info": strings.TrimSpace(authorInfo.String()),
		"content":     content,
	})

	return strings.TrimSpace(master) + template.NewLine + content
}
