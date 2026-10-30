package naming

import "testing"

func TestRemovePunctuationMarks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple punctuation", "Hello, world!", "Hello  world "},
		{"multiple marks", "Hello... World?!", "Hello    World  "},
		{"quotes", `"Hello" ‘World’`, " Hello   World "},
		{"slashes", `one/two\three`, "one two three"},
		{"no punctuation", "Hello World", "Hello World"},
		{"empty string", "", ""},
		{"trim spaces", "  Hello!  ", "Hello "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removePunctuationMarks(tt.input)

			if result != tt.expected {
				t.Errorf("removePunctuationMarks(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRemoveStopWords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"portuguese", "O gato e o vaso", "gato vaso"},
		{"english", "The cat and the dog", "cat dog"},
		{"case insensitive", "O Gato E O Vaso", "Gato Vaso"},
		{"preserve other words", "Uma casa bonita", "casa bonita"},
		{"multiple spaces", "O   gato   e   o vaso", "gato vaso"},
		{"all stop words", "o a e de", ""},
		{"no stop words", "gato bonito", "gato bonito"},
		{"empty string", "", ""},
		{"tabs and newlines", "O\tgato\ne o vaso", "gato vaso"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeStopWords(tt.input)

			if result != tt.expected {
				t.Errorf("removeStopWords(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"portuguese title", "O Coração da Montanha", "coracao-montanha"},
		{"english title", "The Cat and the Dog", "cat-dog"},
		{"accents", "Árvores e Flores", "arvores-flores"},
		{"punctuation", "Olá, Mundo!", "ola-mundo"},
		{"multiple spaces", "  O   Gato   e   o Vaso  ", "gato-vaso"},
		{"mixed case", "UM GATO E UM VASO", "gato-vaso"},
		{"all stop words", "o a e de", ""},
		{"empty string", "", ""},
		{"real book title", "Um gato, um vaso, uma ausência esquecida", "gato-vaso-ausencia-esquecida"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Slugify(tt.input)

			if result != tt.expected {
				t.Errorf("Slugify(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
