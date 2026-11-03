package resolver

import (
	"path/filepath"
	"regexp"
	"testing"
)

func TestVersionRegex(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected string
	}{
		{
			name:     "standalone",
			filename: "v06_velho-novo-mundo.md",
			expected: `^v([0-9]+)_.*\.md$`,
		},
		{
			name:     "chaptered",
			filename: "i0010_v06_velho-novo-mundo.md",
			expected: `^i0010_v([0-9]+)_.*\.md$`,
		},
		{
			name:     "parted",
			filename: "p01_i0010_v06_velho-novo-mundo.md",
			expected: `^p01_i0010_v([0-9]+)_.*\.md$`,
		},
		{
			name:     "full path",
			filename: "/books/novel/parte_01/p01_i0020_v04_dedicatoria.md",
			expected: `^p01_i0020_v([0-9]+)_.*\.md$`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := versionRegex(tt.filename)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.expected {
				t.Errorf("got %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestVersionRegexInvalid(t *testing.T) {
	filenames := []string{
		"",
		"invalid.md",
		"i0010_velho-novo-mundo.md",
		"v_conto.md",
		"i0010_v_conto.md",
		"p01_i0010_v06.md",
	}

	for _, filename := range filenames {
		t.Run(filename, func(t *testing.T) {
			_, err := versionRegex(filename)
			if err == nil {
				t.Errorf("expected error for %q", filename)
			}
		})
	}
}

func TestVersionRegexMatches(t *testing.T) {
	tests := []struct {
		name      string
		reference string
		candidate string
		wantMatch bool
	}{
		{"standalone same manuscript", "v01_conto.md", "v10_outro-titulo.md", true},
		{"chaptered same index", "i0030_v07_reflexos.md", "i0030_v10_vazio-lindo.md", true},
		{"chaptered different index", "i0030_v07_reflexos.md", "i0040_v10_vazio-lindo.md", false},
		{"parted same part and index", "p01_i0020_v04_texto.md", "p01_i0020_v12_novo.md", true},
		{"parted different part", "p01_i0020_v04_texto.md", "p02_i0020_v12_novo.md", false},
		{"parted different index", "p01_i0020_v04_texto.md", "p01_i0030_v12_novo.md", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern, err := versionRegex(tt.reference)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			re := regexp.MustCompile(pattern)
			got := re.MatchString(filepath.Base(tt.candidate))

			if got != tt.wantMatch {
				t.Errorf("match = %t, want %t", got, tt.wantMatch)
			}
		})
	}
}

func TestIterateFilesForMaxValue(t *testing.T) {
	tests := []struct {
		name     string
		files    []string
		pattern  string
		expected int
	}{
		{
			name:     "empty directory",
			files:    []string{},
			pattern:  `(?:^|_)i([0-9]+)(?:_|$)`,
			expected: 0,
		},
		{
			name: "highest index",
			files: []string{
				"i0010_v01_inicio.md",
				"i0017_v02_meio.md",
				"i0030_v01_fim.md",
			},
			pattern:  `(?:^|_)i([0-9]+)(?:_|$)`,
			expected: 30,
		},
		{
			name: "ignore files without indexes",
			files: []string{
				"master.md",
				"v20_conto.md",
				"i0020_v01_capitulo.md",
			},
			pattern:  `(?:^|_)i([0-9]+)(?:_|$)`,
			expected: 20,
		},
		{
			name: "highest version for same chapter",
			files: []string{
				"i0030_v07_reflexos.md",
				"i0030_v10_vazio-lindo.md",
				"i0040_v99_outro.md",
			},
			pattern:  `^i0030_v([0-9]+)_.*\.md$`,
			expected: 10,
		},
		{
			name: "no matching versions",
			files: []string{
				"master.md",
				"i0040_v10_outro.md",
			},
			pattern:  `^i0030_v([0-9]+)_.*\.md$`,
			expected: 0,
		},
		{
			name: "standalone versions",
			files: []string{
				"v01_conto.md",
				"v08_novo-titulo.md",
				"master.md",
			},
			pattern:  `^v([0-9]+)_.*\.md$`,
			expected: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := iterateFilesForMaxValue(tt.files, tt.pattern)

			if got != tt.expected {
				t.Errorf("got %d, want %d", got, tt.expected)
			}
		})
	}
}
