package template

import (
	"strings"
	"unicode"
)

const NewLine = "\n\n"

const StandaloneTemplate = "# «title»"

const PartedTemplate = `# «title»

::: {custom-style="WT_PART_NAME"}
Mrs. Dalloway
:::

::: {custom-style="WT_PART_EPIGRAPH"}
She felt very young; at the same time unspeakably aged.
She sliced like a knife through everything; at the same time was outside, looking on.
:::

::: {custom-style="WT_PART_EPIGRAPH_AUTHOR"}
Virginia Woolf, *Mrs. Dalloway*
:::
`

const MasterTemplate = `# «title»

::: {custom-style="WT_AUTHOR_INFO"}
«author_info»
:::

«content»
`

func Content(template string, options map[string]string) string {
	content := template
	for option, value := range options {
		content = strings.ReplaceAll(content, "«"+option+"»", value)
	}

	return strings.TrimFunc(content, unicode.IsSpace) + NewLine
}
