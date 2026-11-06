package template

import (
	"fmt"
	"regexp"
	"strings"
)

var AuthorInfo = Replacement{
	Identifier: "%%%AuthorInfo",
	Style:      "WT_AUTHOR_INFO",
	Content:    "",
	IsBlock:    true,
}

var Divider = Replacement{
	Identifier: "***",
	Style:      "WT_DIVIDER",
	Content:    "\\* \\* \\*",
	IsBlock:    false,
}

var PartName = Replacement{
	Identifier: "%%%PartName",
	Style:      "WT_PART_NAME",
	Content:    "",
	IsBlock:    true,
}

var PartEpigraph = Replacement{
	Identifier: "%%%PartEpigraph",
	Style:      "WT_PART_EPIGRAPH",
	Content:    "",
	IsBlock:    true,
}

var PartEpigraphAuthor = Replacement{
	Identifier: "%%%PartEpigraphAuthor",
	Style:      "WT_PART_EPIGRAPH_AUTHOR",
	Content:    "",
	IsBlock:    true,
}

var Poem = Replacement{
	Identifier: "%%%Poem",
	Style:      "WT_POEM",
	Content:    "",
	IsBlock:    true,
}

var PoemEven = Replacement{
	Identifier: "%%%PoemEven",
	Style:      "WT_POEM_EVEN",
	Content:    "",
	IsBlock:    true,
}

var TheEnd = Replacement{
	Identifier: "%%%TheEnd",
	Style:      "WT_THE_END",
	Content:    "The End",
	IsBlock:    false,
}

func Replace(page string, replacement Replacement) string {
	identifier := regexp.QuoteMeta(replacement.Identifier)

	if replacement.IsBlock {
		re := regexp.MustCompile(`(?im)^` + identifier + `\r?\n([^\r\n]*(?:\r?\n[^\r\n]+)*)`)

		return re.ReplaceAllStringFunc(page, func(match string) string {
			content := strings.SplitN(match, "\n", 2)[1]

			block := replacement
			block.Content = strings.TrimSpace(content)

			return strings.TrimRight(Style(block), "\n")
		})
	}

	re := regexp.MustCompile(`(?im)^` + identifier + `\r?$`)
	return re.ReplaceAllString(page, strings.TrimRight(Style(replacement), "\n"))
}

func ReplaceAll(page string) string {
	page = Replace(page, Poem)
	page = Replace(page, PoemEven)
	page = Replace(page, PartName)
	page = Replace(page, PartEpigraph)
	page = Replace(page, PartEpigraphAuthor)
	page = Replace(page, TheEnd)

	return Replace(page, Divider)
}

func Style(replacement Replacement) string {
	return fmt.Sprintf("::: {custom-style=\"%s\"}\n%s\n:::\n\n", replacement.Style, replacement.Content)
}
