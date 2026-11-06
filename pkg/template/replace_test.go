package template

import (
	"strings"
	"testing"
)

func TestReplace(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		replacement Replacement
		expected    string
	}{
		{
			name:        "simple marker",
			input:       "Before\n%%%TheEnd\nAfter",
			replacement: TheEnd,
			expected:    "Before\n" + strings.TrimRight(Style(TheEnd), "\n") + "\nAfter",
		},
		{
			name:        "case insensitive uppercase",
			input:       "%%%THEEND",
			replacement: TheEnd,
			expected:    strings.TrimRight(Style(TheEnd), "\n"),
		},
		{
			name:        "case insensitive mixed",
			input:       "%%%tHeEnD",
			replacement: TheEnd,
			expected:    strings.TrimRight(Style(TheEnd), "\n"),
		},
		{
			name:        "simple marker must occupy full line",
			input:       "%%%TheEnd extra",
			replacement: TheEnd,
			expected:    "%%%TheEnd extra",
		},
		{
			name:        "simple marker cannot have prefix",
			input:       "text %%%TheEnd",
			replacement: TheEnd,
			expected:    "text %%%TheEnd",
		},
		{
			name:        "divider",
			input:       "Before\n***\nAfter",
			replacement: Divider,
			expected:    "Before\n" + strings.TrimRight(Style(Divider), "\n") + "\nAfter",
		},
		{
			name:        "block poem",
			input:       "%%%Poem\nFirst line\nSecond line\n\nAfter",
			replacement: Poem,
			expected:    strings.TrimRight(Style(Replacement{Style: "WT_POEM", Content: "First line\nSecond line"}), "\n") + "\n\nAfter",
		},
		{
			name:        "block case insensitive",
			input:       "%%%pOeM\nFirst line\n",
			replacement: Poem,
			expected:    strings.TrimRight(Style(Replacement{Style: "WT_POEM", Content: "First line"}), "\n") + "\n",
		},
		{
			name:        "block marker must occupy full line",
			input:       "%%%Poem extra\nFirst line",
			replacement: Poem,
			expected:    "%%%Poem extra\nFirst line",
		},
		{
			name:  "block content overrides predefined content",
			input: "%%%Poem\nActual content\n\nAfter",
			replacement: Replacement{
				Identifier: "%%%Poem",
				Style:      "WT_POEM",
				Content:    "Ignored content",
				IsBlock:    true,
			},
			expected: strings.TrimRight(Style(Replacement{Style: "WT_POEM", Content: "Actual content"}), "\n") + "\n\nAfter",
		},
		{
			name:        "different block identifiers",
			input:       "%%%PoemEven\nEven content\n\nAfter",
			replacement: Poem,
			expected:    "%%%PoemEven\nEven content\n\nAfter",
		},
		{
			name:        "empty input",
			input:       "",
			replacement: Poem,
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Replace(tt.input, tt.replacement)

			if got != tt.expected {
				t.Errorf("Replace() mismatch:\ngot:  %q\nwant: %q", got, tt.expected)
			}
		})
	}
}
