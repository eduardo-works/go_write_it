package info

import (
	"fmt"

	"github.com/eduardo-works/writeit/pkg/project"
	"github.com/eduardo-works/writeit/pkg/terminal"
	"github.com/eduardo-works/writeit/pkg/tree"
)

func ProjectInfo(p project.Project) {
	terminal.WriteLine("{BrightCyan}{Bold}PROJECT INFORMATION")
	terminal.WriteLine("{BrightCyan}----------------------------------------")
	fmt.Println("")

	tree.PrintLabel("Title", p.Title)
	tree.PrintLabel("Slug", p.Slug)
	tree.PrintLabel("Type", p.ProjectType)
	tree.PrintLabel("Parts", p.Parts)
	tree.PrintLabel("Author", p.Author)
	tree.PrintLabel("Email", p.Email)
	tree.PrintLabel("Phone", p.Phone)
	tree.PrintLabel("Locale", p.Locale)
	tree.PrintLabel("Created", p.CreatedAt)
	tree.PrintLabel("Updated", p.UpdatedAt)
	tree.PrintLabel("Path", p.RootFolder)

	fmt.Println()
	terminal.WriteLine("{BrightCyan}{Bold}PROJECT STRUCTURE")
	terminal.WriteLine("{BrightCyan}----------------------------------------")
	fmt.Println("")

	tree.Print(p)
}
