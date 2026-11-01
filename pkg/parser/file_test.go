package parser

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/eduardo-works/writeit/pkg/project"
)

func TestFileParse(t *testing.T) {
	root := filepath.Join(string(filepath.Separator), "tmp", "writeit-test")

	standalone := project.Project{
		RootFolder:  root,
		ProjectType: project.Standalone,
	}

	chaptered := project.Project{
		RootFolder:  root,
		ProjectType: project.Chaptered,
	}

	parted := project.Project{
		RootFolder:  root,
		ProjectType: project.Parted,
		Parts:       3,
	}

	tests := []struct {
		name    string
		project project.Project
		path    string
		want    ParsedFiled
		wantErr bool
	}{
		{
			name:    "standalone file",
			project: standalone,
			path:    "texto/v01_my-book.md",
			want: ParsedFiled{
				Original: "texto/v01_my-book.md",
				Context:  "texto",
				Part:     -1,
				Index:    -1,
				Version:  1,
				Slug:     "my-book",
				IsPart:   false,
			},
		},
		{
			name:    "standalone file without context",
			project: standalone,
			path:    "v02_my-book.md",
			want: ParsedFiled{
				Original: "v02_my-book.md",
				Context:  "texto",
				Part:     -1,
				Index:    -1,
				Version:  2,
				Slug:     "my-book",
				IsPart:   false,
			},
		},
		{
			name:    "chaptered file",
			project: chaptered,
			path:    "capitulos/i0010_v03_first-chapter.md",
			want: ParsedFiled{
				Original: "capitulos/i0010_v03_first-chapter.md",
				Context:  "capitulos",
				Part:     -1,
				Index:    10,
				Version:  3,
				Slug:     "first-chapter",
				IsPart:   false,
			},
		},
		{
			name:    "chaptered file without context",
			project: chaptered,
			path:    "i0020_v01_second-chapter.md",
			want: ParsedFiled{
				Original: "i0020_v01_second-chapter.md",
				Context:  "capitulos",
				Part:     -1,
				Index:    20,
				Version:  1,
				Slug:     "second-chapter",
				IsPart:   false,
			},
		},
		{
			name:    "parted file with part prefix",
			project: parted,
			path:    "parte_01/p01_i0010_v07_new-world.md",
			want: ParsedFiled{
				Original: "parte_01/p01_i0010_v07_new-world.md",
				Context:  "parte_01",
				Part:     1,
				Index:    10,
				Version:  7,
				Slug:     "new-world",
				IsPart:   false,
			},
		},
		{
			name:    "parted file without part prefix",
			project: parted,
			path:    "parte_02/i0010_v03_new-world.md",
			want: ParsedFiled{
				Original: "parte_02/i0010_v03_new-world.md",
				Context:  "parte_02",
				Part:     2,
				Index:    10,
				Version:  3,
				Slug:     "new-world",
				IsPart:   false,
			},
		},
		{
			name:    "parted file with part in slug",
			project: parted,
			path:    "i0000_parte-03.md",
			want: ParsedFiled{
				Original: "i0000_parte-03.md",
				Context:  "parte_03",
				Part:     3,
				Index:    0,
				Version:  -1,
				Slug:     "parte-03",
				IsPart:   true,
			},
		},
		{
			name:    "parted file without version",
			project: parted,
			path:    "parte_01/p01_i0000_parte-01.md",
			want: ParsedFiled{
				Original: "parte_01/p01_i0000_parte-01.md",
				Context:  "parte_01",
				Part:     1,
				Index:    0,
				Version:  -1,
				Slug:     "parte-01",
				IsPart:   true,
			},
		},
		{
			name:    "absolute path",
			project: parted,
			path:    filepath.Join(root, "parte_01", "p01_i0010_v07_new-world.md"),
			want: ParsedFiled{
				Original: filepath.Join(root, "parte_01", "p01_i0010_v07_new-world.md"),
				Context:  "parte_01",
				Part:     1,
				Index:    10,
				Version:  7,
				Slug:     "new-world",
				IsPart:   false,
			},
		},
		{
			name:    "standalone file with zero version",
			project: standalone,
			path:    "texto/v00_my-book.md",
			wantErr: true,
		},
		{
			name:    "standalone file without version",
			project: standalone,
			path:    "texto/my-book.md",
			wantErr: true,
		},
		{
			name:    "chaptered file with zero version",
			project: chaptered,
			path:    "capitulos/i0010_v00_first-chapter.md",
			wantErr: true,
		},
		{
			name:    "chaptered file without version",
			project: chaptered,
			path:    "capitulos/i0010_first-chapter.md",
			wantErr: true,
		},
		{
			name:    "parted file with zero version",
			project: parted,
			path:    "parte_01/p01_i0010_v00_new-world.md",
			wantErr: true,
		},
		{
			name:    "parted file without version",
			project: parted,
			path:    "parte_01/p01_i0010_new-world.md",
			wantErr: true,
		},
		{
			name:    "part file with version",
			project: parted,
			path:    "parte_01/p01_i0000_v01_parte-01.md",
			wantErr: true,
		},
		{
			name:    "part file with zero version",
			project: parted,
			path:    "parte_01/p01_i0000_v00_parte-01.md",
			wantErr: true,
		},
		{
			name:    "parted file without identifiable part",
			project: parted,
			path:    "i0010_v01_new-world.md",
			wantErr: true,
		},
		{
			name:    "parted file with part above project limit",
			project: parted,
			path:    "parte_04/p04_i0010_v01_new-world.md",
			wantErr: true,
		},
		{
			name:    "parted file with context above project limit",
			project: parted,
			path:    "parte_04/i0010_v01_new-world.md",
			wantErr: true,
		},
		{
			name:    "parted file with slug above project limit",
			project: parted,
			path:    "i0000_parte-04.md",
			wantErr: true,
		},
		{
			name:    "parted file with zero part",
			project: parted,
			path:    "p00_i0010_v01_new-world.md",
			wantErr: true,
		},
		{
			name:    "standalone file with part prefix",
			project: standalone,
			path:    "p01_v01_my-book.md",
			wantErr: true,
		},
		{
			name:    "standalone file with zero part prefix",
			project: standalone,
			path:    "p00_v01_my-book.md",
			wantErr: true,
		},
		{
			name:    "chaptered file with part prefix",
			project: chaptered,
			path:    "p02_i0010_v01_first-chapter.md",
			wantErr: true,
		},
		{
			name:    "standalone file with part slug",
			project: standalone,
			path:    "v01_parte-01.md",
			want: ParsedFiled{
				Original: "v01_parte-01.md",
				Context:  "texto",
				Part:     -1,
				Index:    -1,
				Version:  1,
				Slug:     "parte-01",
				IsPart:   false,
			},
		},

		{
			name:    "standalone file with index",
			project: standalone,
			path:    "texto/i0010_v01_my-book.md",
			want: ParsedFiled{
				Original: "texto/i0010_v01_my-book.md",
				Context:  "texto",
				Part:     -1,
				Index:    -1,
				Version:  1,
				Slug:     "my-book",
				IsPart:   false,
			},
		},
		{
			name:    "chaptered file with zero index",
			project: chaptered,
			path:    "capitulos/i0000_v01_first-chapter.md",
			wantErr: true,
		},
		{
			name:    "parted file with zero index",
			project: parted,
			path:    "parte_01/p01_i0000_v01_new-world.md",
			wantErr: true,
		},
		{
			name:    "parted file without index",
			project: parted,
			path:    "parte_01/p01_v01_new-world.md",
			wantErr: true,
		},

		{
			name:    "part file without index",
			project: parted,
			path:    "parte_01/p01_parte-01.md",
			want: ParsedFiled{
				Original: "parte_01/p01_parte-01.md",
				Context:  "parte_01",
				Part:     1,
				Index:    0,
				Version:  -1,
				Slug:     "parte-01",
				IsPart:   true,
			},
		},
		{
			name:    "part file with nonzero index",
			project: parted,
			path:    "parte_01/p01_i0010_parte-01.md",
			want: ParsedFiled{
				Original: "parte_01/p01_i0010_parte-01.md",
				Context:  "parte_01",
				Part:     1,
				Index:    0,
				Version:  -1,
				Slug:     "parte-01",
				IsPart:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FileParse(tt.project, tt.path)

			if tt.wantErr {
				if !errors.Is(err, project.ErrInvalidValue) {
					t.Fatalf("expected ErrInvalidValue, got %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("FileParse() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
