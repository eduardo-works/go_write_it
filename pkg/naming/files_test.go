package naming

import (
	"errors"
	"testing"

	"github.com/eduardo-works/writeit/pkg/project"
)

func TestManuscript(t *testing.T) {
	standalone := &project.Project{
		ProjectType: project.ProjectType("standalone"),
	}

	chaptered := &project.Project{
		ProjectType: project.ProjectType("chaptered"),
	}

	parted := &project.Project{
		ProjectType: project.ProjectType("parted"),
		Parts:       100,
	}

	tests := []struct {
		name     string
		params   NamingParameters
		expected string
		wantErr  bool
	}{
		{
			name:     "standalone normal",
			params:   NamingParameters{Project: standalone, Version: 1, Title: "O Gato Preto"},
			expected: "v01_gato-preto.md",
		},
		{
			name:     "standalone version zero",
			params:   NamingParameters{Project: standalone, Version: 0, Title: "O Gato Preto"},
			expected: "v01_gato-preto.md",
		},
		{
			name:     "standalone negative version",
			params:   NamingParameters{Project: standalone, Version: -2, Title: "O Gato Preto"},
			expected: "v01_gato-preto.md",
		},
		{
			name:     "standalone version 100",
			params:   NamingParameters{Project: standalone, Version: 100, Title: "O Gato Preto"},
			expected: "v100_gato-preto.md",
		},
		{
			name:     "standalone accents",
			params:   NamingParameters{Project: standalone, Version: 2, Title: "A Ausência"},
			expected: "v02_ausencia.md",
		},
		{
			name:    "standalone empty title",
			params:  NamingParameters{Project: standalone, Version: 1},
			wantErr: true,
		},
		{
			name:     "chaptered normal",
			params:   NamingParameters{Project: chaptered, Index: 1, Version: 1, Title: "O Gato Preto"},
			expected: "i0001_v01_gato-preto.md",
		},
		{
			name:     "chaptered version zero",
			params:   NamingParameters{Project: chaptered, Index: 10, Version: 0, Title: "O Gato Preto"},
			expected: "i0010_v01_gato-preto.md",
		},
		{
			name:     "chaptered negative version",
			params:   NamingParameters{Project: chaptered, Index: 10, Version: -3, Title: "O Gato Preto"},
			expected: "i0010_v01_gato-preto.md",
		},
		{
			name:     "chaptered index 9999",
			params:   NamingParameters{Project: chaptered, Index: 9999, Version: 9, Title: "O Gato Preto"},
			expected: "i9999_v09_gato-preto.md",
		},
		{
			name:     "chaptered index 10000",
			params:   NamingParameters{Project: chaptered, Index: 10000, Version: 10, Title: "O Gato Preto"},
			expected: "i10000_v10_gato-preto.md",
		},
		{
			name:     "chaptered accents",
			params:   NamingParameters{Project: chaptered, Index: 12, Version: 3, Title: "A Ausência"},
			expected: "i0012_v03_ausencia.md",
		},
		{
			name:    "chaptered index zero",
			params:  NamingParameters{Project: chaptered, Index: 0, Version: 1, Title: "O Gato Preto"},
			wantErr: true,
		},
		{
			name:    "chaptered negative index",
			params:  NamingParameters{Project: chaptered, Index: -1, Version: 1, Title: "O Gato Preto"},
			wantErr: true,
		},
		{
			name:    "chaptered empty title",
			params:  NamingParameters{Project: chaptered, Index: 10, Version: 1},
			wantErr: true,
		},
		{
			name:     "parted normal",
			params:   NamingParameters{Project: parted, Part: 1, Index: 1, Version: 1, Title: "O Gato Preto"},
			expected: "p01_i0001_v01_gato-preto.md",
		},
		{
			name:     "parted version zero",
			params:   NamingParameters{Project: parted, Part: 1, Index: 1, Version: 0, Title: "O Gato Preto"},
			expected: "p01_i0001_v01_gato-preto.md",
		},
		{
			name:     "parted negative version",
			params:   NamingParameters{Project: parted, Part: 1, Index: 1, Version: -2, Title: "O Gato Preto"},
			expected: "p01_i0001_v01_gato-preto.md",
		},
		{
			name:     "parted multiple digits",
			params:   NamingParameters{Project: parted, Part: 12, Index: 123, Version: 15, Title: "O Gato Preto"},
			expected: "p12_i0123_v15_gato-preto.md",
		},
		{
			name:     "parted large values",
			params:   NamingParameters{Project: parted, Part: 100, Index: 10000, Version: 100, Title: "O Gato Preto"},
			expected: "p100_i10000_v100_gato-preto.md",
		},
		{
			name:     "parted accents",
			params:   NamingParameters{Project: parted, Part: 2, Index: 3, Version: 4, Title: "A Ausência"},
			expected: "p02_i0003_v04_ausencia.md",
		},
		{
			name:    "parted part zero",
			params:  NamingParameters{Project: parted, Part: 0, Index: 1, Version: 1, Title: "O Gato Preto"},
			wantErr: true,
		},
		{
			name:    "parted negative part",
			params:  NamingParameters{Project: parted, Part: -1, Index: 1, Version: 1, Title: "O Gato Preto"},
			wantErr: true,
		},
		{
			name:    "parted part outside range",
			params:  NamingParameters{Project: parted, Part: 101, Index: 1, Version: 1, Title: "O Gato Preto"},
			wantErr: true,
		},
		{
			name:    "parted index zero",
			params:  NamingParameters{Project: parted, Part: 1, Index: 0, Version: 1, Title: "O Gato Preto"},
			wantErr: true,
		},
		{
			name:    "parted negative index",
			params:  NamingParameters{Project: parted, Part: 1, Index: -1, Version: 1, Title: "O Gato Preto"},
			wantErr: true,
		},
		{
			name:    "parted empty title",
			params:  NamingParameters{Project: parted, Part: 1, Index: 1, Version: 1},
			wantErr: true,
		},
		{
			name:    "nil project",
			params:  NamingParameters{Version: 1, Title: "O Gato Preto"},
			wantErr: true,
		},
		{
			name: "invalid project type",
			params: NamingParameters{
				Project: &project.Project{
					ProjectType: project.ProjectType("invalid"),
				},
				Version: 1,
				Title:   "O Gato Preto",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Manuscript(tt.params)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got filename %q", got)
				}
				if got != "" {
					t.Errorf("expected empty filename on error, got %q", got)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.expected {
				t.Errorf("Manuscript() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestManuscriptNilProject(t *testing.T) {
	_, err := Manuscript(NamingParameters{
		Version: 1,
		Title:   "O Gato Preto",
	})

	if !errors.Is(err, project.ErrNoProject) {
		t.Errorf("expected ErrNoProject, got %v", err)
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
			got := PartedPart(tt.part)

			if got != tt.expected {
				t.Errorf("PartedPart(%d) = %q, want %q",
					tt.part, got, tt.expected)
			}
		})
	}
}
