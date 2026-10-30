package naming

import (
	"slices"
	"strings"

	"github.com/gosimple/slug"
)

var punctuation = [...]string{
	".", ",", ";", ":", "!", "?",
	"'", "\"", "‘", "’", "“", "”",
	"…", "/", "\\",
}

var stopWordsEN = []string{
	"a", "an", "the",
	"of", "in", "on", "at",
	"by", "for", "from",
	"with", "to",
	"and", "or",
	"that",
	"if",
}

var stopWordsPT = []string{
	"a", "as", "o", "os",
	"um", "uma", "uns", "umas",
	"de", "da", "das", "do", "dos",
	"em", "na", "nas", "no", "nos",
	"por", "pela", "pelas", "pelo", "pelos",
	"para",
	"com",
	"e", "ou",
	"que",
	"se",
}

func Slugify(text string) string {
	text = strings.TrimSpace(text)
	text = removePunctuationMarks(text)
	text = removeStopWords(text)

	return slug.Make(text)
}

func PortugueseNumbersToWords(n int) string {
	words := []string{
		"zero", "um", "dois", "três", "quatro",
		"cinco", "seis", "sete", "oito", "nove", "dez",
	}

	if n < 0 || n >= len(words) {
		return ""
	}

	return words[n]
}

func removePunctuationMarks(text string) string {
	text = strings.TrimSpace(text)

	for _, item := range punctuation {
		text = strings.ReplaceAll(text, item, " ")
	}

	return text
}

func removeStopWords(text string) string {
	applied := []string{}

	for _, word := range strings.Fields(text) {
		lowerCased := strings.ToLower(word)

		if slices.Contains(stopWordsEN, lowerCased) ||
			slices.Contains(stopWordsPT, lowerCased) {
			continue
		}

		applied = append(applied, word)
	}

	return strings.Join(applied, " ")
}
