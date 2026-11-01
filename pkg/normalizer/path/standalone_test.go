package path

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/eduardo-works/writeit/pkg/project"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		name        string
		projectType project.ProjectType
		projectSlug string
		parts       int
		filename    string
		wantName    string
		wantChanged bool
		wantErr     bool
	}{
		{
			name:        "standalone already normalized",
			projectType: project.Standalone,
			projectSlug: "my-book",
			filename:    "v01_my-book.md",
			wantName:    "v01_my-book.md",
		},
		{
			name:        "standalone requires normalization",
			projectType: project.Standalone,
			projectSlug: "my-book",
			filename:    "v03_old-title.md",
			wantName:    "v03_my-book.md",
			wantChanged: true,
		},
		{
			name:        "chaptered already normalized",
			projectType: project.Chaptered,
			projectSlug: "my-book",
			filename:    "i0010_v02_first-chapter.md",
			wantName:    "i0010_v02_first-chapter.md",
		},
		{
			name:        "chaptered requires normalization",
			projectType: project.Chaptered,
			projectSlug: "my-book",
			filename:    "i10_v2_first-chapter.md",
			wantName:    "i0010_v02_first-chapter.md",
			wantChanged: true,
		},
		{
			name:        "parted chapter already normalized",
			projectType: project.Parted,
			projectSlug: "my-book",
			parts:       3,
			filename:    "p01_i0010_v02_first-chapter.md",
			wantName:    "p01_i0010_v02_first-chapter.md",
		},
		{
			name:        "parted chapter requires normalization",
			projectType: project.Parted,
			projectSlug: "my-book",
			parts:       3,
			filename:    "p1_i10_v2_first-chapter.md",
			wantName:    "p01_i0010_v02_first-chapter.md",
			wantChanged: true,
		},
		{
			name:        "part file already normalized",
			projectType: project.Parted,
			projectSlug: "my-book",
			parts:       3,
			filename:    "p01_i0000_parte-01.md",
			wantName:    "p01_i0000_parte-01.md",
		},
		{
			name:        "part file requires normalization",
			projectType: project.Parted,
			projectSlug: "my-book",
			parts:       3,
			filename:    "i0000_parte-01.md",
			wantName:    "p01_i0000_parte-01.md",
			wantChanged: true,
		},
		{
			name:        "standalone missing version",
			projectType: project.Standalone,
			projectSlug: "my-book",
			filename:    "my-book.md",
			wantErr:     true,
		},
		{
			name:        "chaptered missing index",
			projectType: project.Chaptered,
			projectSlug: "my-book",
			filename:    "v01_first-chapter.md",
			wantErr:     true,
		},
		{
			name:        "chaptered zero index",
			projectType: project.Chaptered,
			projectSlug: "my-book",
			filename:    "i0000_v01_first-chapter.md",
			wantErr:     true,
		},
		{
			name:        "parted chapter missing version",
			projectType: project.Parted,
			projectSlug: "my-book",
			parts:       3,
			filename:    "p01_i0010_first-chapter.md",
			wantErr:     true,
		},
		{
			name:        "parted chapter outside part range",
			projectType: project.Parted,
			projectSlug: "my-book",
			parts:       3,
			filename:    "p04_i0010_v01_first-chapter.md",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()

			p := project.Project{
				RootFolder:  root,
				ProjectType: tt.projectType,
				Slug:        tt.projectSlug,
				Parts:       tt.parts,
			}

			context := ""

			switch tt.projectType {
			case project.Standalone:
				context = "texto"
			case project.Chaptered:
				context = "capitulos"
			case project.Parted:
				context = "parte_01"
			}

			dir := filepath.Join(root, context)

			if err := os.Mkdir(dir, 0755); err != nil {
				t.Fatal(err)
			}

			source := filepath.Join(dir, tt.filename)

			if err := os.WriteFile(source, []byte("test content"), 0644); err != nil {
				t.Fatal(err)
			}

			got := Normalize(&p, source)

			if got.Project != &p {
				t.Error("Normalize() project pointer does not match")
			}

			if got.Source != source {
				t.Errorf("Normalize() source = %q, want %q", got.Source, source)
			}

			if tt.wantErr {
				if !errors.Is(got.Issue, project.ErrInvalidValue) {
					t.Fatalf("expected ErrInvalidValue, got %v", got.Issue)
				}

				if got.Target != source {
					t.Errorf("Normalize() target = %q, want original %q", got.Target, source)
				}

				if got.IsChanged() {
					t.Error("Normalize() changed should be false on error")
				}

				return
			}

			if got.Issue != nil {
				t.Fatalf("unexpected error: %v", got.Issue)
			}

			want := filepath.Join(dir, tt.wantName)

			if got.Target != want {
				t.Errorf("Normalize() target = %q, want %q", got.Target, want)
			}

			if got.IsChanged() != tt.wantChanged {
				t.Errorf(
					"Normalize() changed = %t, want %t",
					got.IsChanged(),
					tt.wantChanged,
				)
			}

			wantSourceRel := filepath.Join(context, tt.filename)
			wantTargetRel := filepath.Join(context, tt.wantName)

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
