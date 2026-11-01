package path

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/eduardo-works/writeit/pkg/project"
)

func TestNormalizeChaptered(t *testing.T) {
	tests := []struct {
		name        string
		projectType project.ProjectType
		filename    string
		wantName    string
		wantChanged bool
		wantErr     bool
	}{
		{
			name:        "already normalized",
			projectType: project.Chaptered,
			filename:    "i0010_v01_first-chapter.md",
			wantName:    "i0010_v01_first-chapter.md",
		},
		{
			name:        "short index",
			projectType: project.Chaptered,
			filename:    "i10_v01_first-chapter.md",
			wantName:    "i0010_v01_first-chapter.md",
			wantChanged: true,
		},
		{
			name:        "short version",
			projectType: project.Chaptered,
			filename:    "i0010_v2_first-chapter.md",
			wantName:    "i0010_v02_first-chapter.md",
			wantChanged: true,
		},
		{
			name:        "short index and version",
			projectType: project.Chaptered,
			filename:    "i10_v2_first-chapter.md",
			wantName:    "i0010_v02_first-chapter.md",
			wantChanged: true,
		},
		{
			name:        "reordered fields",
			projectType: project.Chaptered,
			filename:    "v02_i0010_first-chapter.md",
			wantName:    "i0010_v02_first-chapter.md",
			wantChanged: true,
		},
		{
			name:        "higher index and version",
			projectType: project.Chaptered,
			filename:    "i0125_v15_final-chapter.md",
			wantName:    "i0125_v15_final-chapter.md",
		},
		{
			name:        "missing index",
			projectType: project.Chaptered,
			filename:    "v01_first-chapter.md",
			wantErr:     true,
		},
		{
			name:        "zero index",
			projectType: project.Chaptered,
			filename:    "i0000_v01_first-chapter.md",
			wantErr:     true,
		},
		{
			name:        "missing version",
			projectType: project.Chaptered,
			filename:    "i0010_first-chapter.md",
			wantErr:     true,
		},
		{
			name:        "zero version",
			projectType: project.Chaptered,
			filename:    "i0010_v00_first-chapter.md",
			wantErr:     true,
		},
		{
			name:        "part prefix not allowed",
			projectType: project.Chaptered,
			filename:    "p01_i0010_v01_first-chapter.md",
			wantErr:     true,
		},
		{
			name:        "non chaptered project",
			projectType: project.Standalone,
			filename:    "v01_my-book.md",
			wantName:    "v01_my-book.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()

			p := project.Project{
				RootFolder:  root,
				ProjectType: tt.projectType,
				Slug:        "my-book",
			}

			context := "capitulos"
			if tt.projectType == project.Standalone {
				context = "texto"
			}

			dir := filepath.Join(root, context)

			if err := os.Mkdir(dir, 0755); err != nil {
				t.Fatal(err)
			}

			source := filepath.Join(dir, tt.filename)

			if err := os.WriteFile(source, []byte("test content"), 0644); err != nil {
				t.Fatal(err)
			}

			got := NormalizeChaptered(&p, source)

			if got.Project != &p {
				t.Error("NormalizeChaptered() project pointer does not match")
			}

			if got.Source != source {
				t.Errorf("NormalizeChaptered() source = %q, want %q", got.Source, source)
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
				t.Errorf("NormalizeChaptered() target = %q, want %q", got.Target, want)
			}

			if got.IsChanged() != tt.wantChanged {
				t.Errorf(
					"NormalizeChaptered() changed = %t, want %t",
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
