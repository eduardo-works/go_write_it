package assemble

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/eduardo-works/writeit/pkg/parser"
	"github.com/eduardo-works/writeit/pkg/project"
	"github.com/eduardo-works/writeit/pkg/template"
	"github.com/eduardo-works/writeit/pkg/terminal"
)

func generateMasterContent(p project.Project, files []string) string {
	content := ""

	for _, file := range files {
		page := readAndFixPage(p, file)
		if page == "" {
			continue
		}

		page = template.ReplaceAll(page)

		content = strings.TrimSpace(content)

		if content != "" {
			content += template.NewLine
		}

		content += page
	}

	return content
}

func readAndFixPage(p project.Project, file string) string {
	parsed, err := parser.FileParse(p, file)
	if err != nil {
		terminal.Warning(fmt.Sprintf("Something failed while parsing %s: %v", file, err))
		return ""
	}

	data, err := os.ReadFile(file)
	if err != nil {
		terminal.Warning(fmt.Sprintf("Something failed while reading %s: %v", file, err))
		return ""
	}

	page := strings.TrimSpace(string(data))

	re := regexp.MustCompile(`(?m)^(#{1,6}) +(\p{L}.*)$`)

	firstLine, _, _ := strings.Cut(page, "\n")

	if !re.MatchString(firstLine) {
		terminal.Warning(fmt.Sprintf("File does not start with a valid header: %s", file))
		return ""
	}

	levels := 1
	if p.IsParted() && !parsed.IsPart {
		levels = 2
	}

	return re.ReplaceAllStringFunc(page, func(line string) string {
		matches := re.FindStringSubmatch(line)

		level := len(matches[1]) + levels
		if level > 6 {
			level = 6
		}

		return strings.Repeat("#", level) + " " + matches[2]
	})
}

func formatNumber(n int) string {
	s := strconv.Itoa(n)

	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + " " + s[i:]
	}

	return s
}

func countContent(content string) ContentInfo {
	words := len(strings.Fields(content))
	chars := utf8.RuneCountInString(content)

	charsNoSpaces := 0
	for _, r := range content {
		if !unicode.IsSpace(r) {
			charsNoSpaces++
		}
	}

	pages := func(size int) string {
		return formatNumber(int(math.Ceil(float64(chars) / float64(size))))
	}

	return ContentInfo{
		Words:                  formatNumber(words),
		CharactersWithSpace:    formatNumber(chars),
		CharactersWithoutSpace: formatNumber(charsNoSpaces),
		Pages1000:              pages(1000),
		Pages1400:              pages(1400),
		Pages1800:              pages(1800),
		Pages2100:              pages(2100),
	}
}
