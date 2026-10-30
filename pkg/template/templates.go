package template

import "strings"

const StandaloneTemplate = "# «title»"

const PartedTemplate = `# «title»"

::: {custom-style="PART_NAME"}
Mrs. Dalloway
:::

::: {custom-style="PART_EPIGRAPH"}
She felt very young; at the same time unspeakably aged.
She sliced like a knife through everything; at the same time was outside, looking on.
:::

::: {custom-style="PART_EPIGRAPH_AUTHOR"}
Virginia Woolf, *Mrs. Dalloway*
:::
`

func Content(template string, options map[string]string) string {
	content := template
	for option, value := range options {
		content = strings.ReplaceAll(content, "«"+option+"»", value)
	}

	return content
}
