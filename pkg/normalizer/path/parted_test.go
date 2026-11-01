package path

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/eduardo-works/writeit/pkg/project"
)

func TestNormalizeParted(t *testing.T) {
	tests := []struct {
		name        string
		projectType project.ProjectType
		parts       int
		folder      string
		filename    string
		wantName    string
		wantChanged bool
		wantErr     bool
	}{
		{
			name:        "part file already normalized",
			projectType: project.Parted,
			parts:       3,
			folder:      "parte_01",
			filename:    "p01_i0000_parte-01.md",
			wantName:    "p01_i0000_parte-01.md",
		},
		{
			name:        "part file missing prefix",
			projectType: project.Parted,
			parts:       3,
			folder:      "parte_01",
			filename:    "i0000_parte-01.md",
			wantName:    "p01_i0000_parte-01.md",
			wantChanged: true,
		},
		{
			name:        "part file missing index",
			projectType: project.Parted,
			parts:       3,
			folder:      "parte_02",
			filename:    "p02_parte-02.md",
			wantName:    "p02_i0000_parte-02.md",
			wantChanged: true,
		},
		{
			name:        "part file with nonzero index",
			projectType: project.Parted,
			parts:       3,
			folder:      "parte_03",
			filename:    "p03_i0010_parte-03.md",
			wantName:    "p03_i0000_parte-03.md",
			wantChanged: true,
		},
		{
			name:        "part file without folder prefix",
			projectType: project.Parted,
			parts:       3,
			folder:      "parte_02",
			filename:    "i0000_parte-02.md",
			wantName:    "p02_i0000_parte-02.md",
			wantChanged: true,
		},
		{
			name:        "chapter already normalized",
			projectType: project.Parted,
			parts:       3,
			folder:      "parte_01",
			filename:    "p01_i0010_v02_first-chapter.md",
			wantName:    "p01_i0010_v02_first-chapter.md",
		},
		{
			name:        "chapter missing part prefix",
			projectType: project.Parted,
			parts:       3,
			folder:      "parte_02",
			filename:    "i0010_v03_second-chapter.md",
			wantName:    "p02_i0010_v03_second-chapter.md",
			wantChanged: true,
		},
		{
			name:        "chapter with short numeric fields",
			projectType: project.Parted,
			parts:       3,
			folder:      "parte_01",
			filename:    "p1_i10_v2_first-chapter.md",
			wantName:    "p01_i0010_v02_first-chapter.md",
			wantChanged: true,
		},
		{
			name:        "chapter with reordered fields",
			projectType: project.Parted,
			parts:       3,
			folder:      "parte_01",
			filename:    "v02_i0010_p01_first-chapter.md",
			wantName:    "p01_i0010_v02_first-chapter.md",
			wantChanged: true,
		},
		{
			name:        "chapter with zero index",
			projectType: project.Parted,
			parts:       3,
			folder:      "parte_01",
			filename:    "p01_i0000_v01_first-chapter.md",
			wantErr:     true,
		},
		{
			name:        "chapter without index",
			projectType: project.Parted,
			parts:       3,
			folder:      "parte_01",
			filename:    "p01_v01_first-chapter.md",
			wantErr:     true,
		},
		{
			name:        "chapter without version",
			projectType: project.Parted,
			parts:       3,
			folder:      "parte_01",
			filename:    "p01_i0010_first-chapter.md",
			wantErr:     true,
		},
		{
			name:        "chapter with zero version",
			projectType: project.Parted,
			parts:       3,
			folder:      "parte_01",
			filename:    "p01_i0010_v00_first-chapter.md",
			wantErr:     true,
		},
		{
			name:        "chapter with part above limit",
			projectType: project.Parted,
			parts:       3,
			folder:      "parte_04",
			filename:    "p04_i0010_v01_first-chapter.md",
			wantErr:     true,
		},
		{
			name:        "part file with version",
			projectType: project.Parted,
			parts:       3,
			folder:      "parte_01",
			filename:    "p01_i0000_v01_parte-01.md",
			wantErr:     true,
		},
		{
			name:        "part file with part above limit",
			projectType: project.Parted,
			parts:       3,
			folder:      "parte_04",
			filename:    "p04_i0000_parte-04.md",
			wantErr:     true,
		},
		{
			name:        "non parted project",
			projectType: project.Chaptered,
			folder:      "capitulos",
			filename:    "i0010_v01_first-chapter.md",
			wantName:    "i0010_v01_first-chapter.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()

			p := project.Project{
				RootFolder:  root,
				ProjectType: tt.projectType,
				Parts:       tt.parts,
				Slug:        "my-book",
			}

			dir := filepath.Join(root, tt.folder)

			if err := os.Mkdir(dir, 0755); err != nil {
				t.Fatal(err)
			}

			source := filepath.Join(dir, tt.filename)

			if err := os.WriteFile(source, []byte("test content"), 0644); err != nil {
				t.Fatal(err)
			}

			got := NormalizeParted(&p, source)

			if got.Project != &p {
				t.Error("NormalizeParted() project pointer does not match")
			}

			if got.Source != source {
				t.Errorf("NormalizeParted() source = %q, want %q", got.Source, source)
			}

			if tt.wantErr {
				if !errors.Is(got.Issue, project.ErrInvalidValue) {
					t.Fatalf("expected ErrInvalidValue, got %v", got.Issue)
				}

				if got.Target != source {
					t.Errorf("target = %q, want original %q", got.Target, source)
				}

				if got.IsChanged() {
					t.Error("changed should be false on error")
				}

				return
			}

			if got.Issue != nil {
				t.Fatalf("unexpected error: %v", got.Issue)
			}

			want := filepath.Join(dir, tt.wantName)

			if got.Target != want {
				t.Errorf("NormalizeParted() target = %q, want %q", got.Target, want)
			}

			if got.IsChanged() != tt.wantChanged {
				t.Errorf(
					"NormalizeParted() changed = %t, want %t",
					got.IsChanged(),
					tt.wantChanged,
				)
			}

			wantSourceRel := filepath.Join(tt.folder, tt.filename)
			wantTargetRel := filepath.Join(tt.folder, tt.wantName)

			if got.SourceRel != wantSourceRel {
				t.Errorf("source relative = %q, want %q", got.SourceRel, wantSourceRel)
			}

			if got.TargetRel != wantTargetRel {
				t.Errorf("target relative = %q, want %q", got.TargetRel, wantTargetRel)
			}

			wantDate := time.Now().Format("2006-01-02")

			if got.Date != wantDate {
				t.Errorf("date = %q, want %q", got.Date, wantDate)
			}

			if _, err := os.Stat(source); err != nil {
				t.Errorf("original file should still exist: %v", err)
			}
		})
	}
}
