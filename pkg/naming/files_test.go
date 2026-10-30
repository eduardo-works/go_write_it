package naming

import "testing"

func TestStandaloneManuscript(t *testing.T) {
	tests := []struct {
		name     string
		version  int
		title    string
		expected string
	}{
		{"normal", 1, "O Gato Preto", "v01_gato-preto.md"},
		{"version zero", 0, "O Gato Preto", "v01_gato-preto.md"},
		{"negative version", -2, "O Gato Preto", "v01_gato-preto.md"},
		{"version 9", 9, "O Gato Preto", "v09_gato-preto.md"},
		{"version 10", 10, "O Gato Preto", "v10_gato-preto.md"},
		{"version 100", 100, "O Gato Preto", "v100_gato-preto.md"},
		{"accents", 2, "A Ausência", "v02_ausencia.md"},
		{"empty title", 1, "", "v01_.md"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StandaloneManuscript(tt.version, tt.title)

			if result != tt.expected {
				t.Errorf("StandaloneManuscript(%d, %q) = %q, want %q",
					tt.version, tt.title, result, tt.expected)
			}
		})
	}
}

func TestChapteredManuscript(t *testing.T) {
	tests := []struct {
		name     string
		index    int
		version  int
		title    string
		expected string
	}{
		{"normal", 1, 1, "O Gato Preto", "i0001_v01_gato-preto.md"},
		{"index zero", 0, 1, "O Gato Preto", "i0000_v01_gato-preto.md"},
		{"negative index", -1, 1, "O Gato Preto", ""},
		{"version zero", 1, 0, "O Gato Preto", "i0001_v01_gato-preto.md"},
		{"negative version", 1, -3, "O Gato Preto", "i0001_v01_gato-preto.md"},
		{"index 9999", 9999, 9, "O Gato Preto", "i9999_v09_gato-preto.md"},
		{"index 10000", 10000, 10, "O Gato Preto", "i10000_v10_gato-preto.md"},
		{"accents", 12, 3, "A Ausência", "i0012_v03_ausencia.md"},
		{"empty title", 1, 1, "", "i0001_v01_.md"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ChapteredManuscript(tt.index, tt.version, tt.title)

			if result != tt.expected {
				t.Errorf("ChapteredManuscript(%d, %d, %q) = %q, want %q",
					tt.index, tt.version, tt.title, result, tt.expected)
			}
		})
	}
}

func TestPartedPart(t *testing.T) {
	tests := []struct {
		name     string
		part     int
		expected string
	}{
		{"part 1", 1, "p01_i0000_parte-01.md"},
		{"part 9", 9, "p09_i0000_parte-09.md"},
		{"part 10", 10, "p10_i0000_parte-10.md"},
		{"part 99", 99, "p99_i0000_parte-99.md"},
		{"part 100", 100, "p100_i0000_parte-100.md"},
		{"part zero", 0, ""},
		{"negative part", -1, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PartedPart(tt.part)

			if result != tt.expected {
				t.Errorf("PartedPart(%d) = %q, want %q",
					tt.part, result, tt.expected)
			}
		})
	}
}

func TestPartedManuscript(t *testing.T) {
	tests := []struct {
		name     string
		part     int
		index    int
		version  int
		title    string
		expected string
	}{
		{"normal", 1, 1, 1, "O Gato Preto", "p01_i0001_v01_gato-preto.md"},
		{"part zero", 0, 1, 1, "O Gato Preto", ""},
		{"negative part", -1, 1, 1, "O Gato Preto", ""},
		{"index zero", 1, 0, 1, "O Gato Preto", "p01_i0000_v01_gato-preto.md"},
		{"negative index", 1, -1, 1, "O Gato Preto", ""},
		{"version zero", 1, 1, 0, "O Gato Preto", "p01_i0001_v01_gato-preto.md"},
		{"negative version", 1, 1, -2, "O Gato Preto", "p01_i0001_v01_gato-preto.md"},
		{"multiple digits", 12, 123, 15, "O Gato Preto", "p12_i0123_v15_gato-preto.md"},
		{"large values", 100, 10000, 100, "O Gato Preto", "p100_i10000_v100_gato-preto.md"},
		{"accents", 2, 3, 4, "A Ausência", "p02_i0003_v04_ausencia.md"},
		{"empty title", 1, 1, 1, "", "p01_i0001_v01_.md"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PartedManuscript(tt.part, tt.index, tt.version, tt.title)

			if result != tt.expected {
				t.Errorf("PartedManuscript(%d, %d, %d, %q) = %q, want %q",
					tt.part, tt.index, tt.version, tt.title, result, tt.expected)
			}
		})
	}
}
