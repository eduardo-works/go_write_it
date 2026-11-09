package main

import (
	"github.com/eduardo-works/writeit/pkg/terminal"
	"github.com/spf13/cobra"
)

var stylesCmd = &cobra.Command{
	Use:   "styles",
	Short: "Show manuscript formatting markers and examples.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		showStyles()
	},
}

func configureStyles() {
	rootCmd.AddCommand(stylesCmd)
}

func showStyles() {
	terminal.WriteLine("")
	terminal.WriteLine("     {Bold}{BrightCyan}WRITEIT / MANUSCRIPT STYLES")
	terminal.WriteLine("     {BrightBlack}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	terminal.WriteLine("")
	terminal.WriteLine("     {White}Writeit supports special Markdown markers to format manuscript content automatically.")
	terminal.WriteLine("     {Dim}Markers must be written on their own line. Block markers apply to the consecutive non-empty lines immediately below them.")
	terminal.WriteLine("")

	showStyle(
		"01", "AUTHOR INFORMATION", "%%%AuthorInfo",
		"Formats an author information block, such as a short biography or contact details.",
		[]string{
			"%%%AuthorInfo",
			"Jane Doe",
			"Writer and editor based in Lisbon.",
			"jane@example.com",
		},
	)

	showStyle(
		"02", "SCENE DIVIDER", "***",
		"Separates scenes or narrative sections without starting a new chapter. Place the marker on its own line, with blank lines around it.",
		[]string{
			"The house was silent. Nobody had returned.",
			"",
			"***",
			"",
			"Three days later, everything had changed.",
		},
	)

	showStyle(
		"03", "PART NAME", "%%%PartName",
		"Formats the title of a new part of the manuscript. Write the title immediately below the marker.",
		[]string{
			"%%%PartName",
			"Part One — The Beginning",
		},
	)

	showStyle(
		"04", "PART EPIGRAPH", "%%%PartEpigraph",
		"Formats a quotation or short passage placed at the beginning of a part or major section.",
		[]string{
			"%%%PartEpigraph",
			"Every journey begins with a question.",
			"Sometimes the answer is not what we expect.",
		},
	)

	showStyle(
		"05", "PART EPIGRAPH AUTHOR", "%%%PartEpigraphAuthor",
		"Identifies the author or source of the preceding epigraph. Supports Markdown emphasis for book titles.",
		[]string{
			"%%%PartEpigraphAuthor",
			"Jane Doe, *The Imaginary Book*",
		},
	)

	showStyle(
		"06", "POEM", "%%%Poem",
		"Formats a poem or stanza using the primary poetry layout. Use Markdown hard line breaks when necessary.",
		[]string{
			"%%%Poem",
			"The morning light returns,\\",
			"Across the hills and silent trees,\\",
			"And every shadow slowly fades.",
		},
	)

	showStyle(
		"07", "POEM — ALTERNATE", "%%%PoemEven",
		"Formats a poem or stanza using the alternate poetry layout. Useful when alternating the presentation of stanzas.",
		[]string{
			"%%%PoemEven",
			"Beneath the stars I found my way,\\",
			"Through distant fields and falling rain,\\",
			"Until the night became the day.",
		},
	)

	showStyle(
		"08", "THE END", "%%%TheEnd",
		"Places a specially formatted closing message at the end of the manuscript. No additional text is required.",
		[]string{
			"The story had finally come to an end.",
			"",
			"%%%TheEnd",
		},
	)

	terminal.WriteLine("     {BrightBlack}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	terminal.WriteLine("")
	terminal.WriteLine("     {Bold}{BrightCyan}NOTES")
	terminal.WriteLine("     {Dim}• Markers are replaced when assembling the master document.")
	terminal.WriteLine("     {Dim}• Original Markdown manuscripts remain unchanged.")
	terminal.WriteLine("     {Dim}• For block markers, a blank line marks the end of the block.")
	terminal.WriteLine("     {Dim}• Leave a blank line between separate formatting blocks.")
	terminal.WriteLine("")
}

func showStyle(number, title, syntax, description string, example []string) {
	terminal.WriteLine("     {BrightBlack}────────────────────────────────────────────────────────────────────────────────────────────────────")
	terminal.WriteLine("     {Bold}{BrightYellow}" + number + "  " + title)
	terminal.WriteLine("")
	terminal.WriteLine("     {White}" + description)
	terminal.WriteLine("")
	terminal.WriteLine("     {BrightGreen}SYNTAX   {Bold}{BrightMagenta}" + syntax)
	terminal.WriteLine("")
	terminal.WriteLine("     {BrightGreen}EXAMPLE")
	terminal.WriteLine("")

	for _, line := range example {
		if line == "" {
			terminal.WriteLine("")
			continue
		}

		if line == syntax {
			terminal.WriteLine("          {Bold}{BrightMagenta}" + line)
			continue
		}

		terminal.WriteLine("          {BrightWhite}" + line)
	}

	terminal.WriteLine("")
}
