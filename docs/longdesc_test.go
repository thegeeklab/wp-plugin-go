package docs

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLongDescriptions(t *testing.T) {
	const baseSrc = `package flags

import "github.com/urfave/cli/v3"

func flags() []cli.Flag {
	return []cli.Flag{
		%s
	}
}
`

	tests := []struct {
		name    string
		body    string
		want    map[string]string
		wantErr bool
	}{
		{
			name: "single-line comments",
			body: `// Short help, short comment.
		&cli.StringFlag{
			Name:    "foo",
			Usage:   "foo usage",
			Sources: cli.EnvVars("PLUGIN_FOO"),
		},`,
			want: map[string]string{
				"foo": "Short help, short comment.",
			},
		},
		{
			name: "multi-line within a paragraph",
			body: `// Long help, split across source lines
		// within the same paragraph.
		&cli.BoolFlag{
			Name:    "bar",
			Usage:   "bar usage",
			Sources: cli.EnvVars("PLUGIN_BAR"),
		},`,
			want: map[string]string{
				"bar": "Long help, split across source lines within the same paragraph.",
			},
		},
		{
			name: "multiple paragraphs",
			body: `// Long help with multiple paragraphs.
		//
		// The second paragraph has multiple lines
		// spanning two source lines.
		&cli.IntFlag{
			Name:    "baz",
			Usage:   "baz usage",
			Sources: cli.EnvVars("PLUGIN_BAZ"),
		},`,
			want: map[string]string{
				"baz": "Long help with multiple paragraphs.\n\nThe second paragraph has multiple lines spanning two source lines.",
			},
		},
		{
			name: "no comment, flag is skipped",
			body: `&cli.StringSliceFlag{
			Name:    "no-comment",
			Usage:   "no comment usage",
			Sources: cli.EnvVars("PLUGIN_NO_COMMENT"),
		},`,
			want: map[string]string{},
		},
		{
			name: "blank line breaks association",
			body: `// Floating comment that is not the leading
// comment of any literal below.

		&cli.StringFlag{
			Name:    "no-comment",
			Usage:   "no leading comment",
			Sources: cli.EnvVars("PLUGIN_NO_COMMENT"),
		},`,
			want: map[string]string{},
		},
		{
			name:    "missing file returns error",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "flags.go")

			if tt.wantErr {
				path = filepath.Join(dir, "does-not-exist.go")
			} else if err := os.WriteFile(path, []byte(fmt.Sprintf(baseSrc, tt.body)), 0o600); err != nil {
				t.Fatal(err)
			}

			got, err := LongDescriptions(path)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, got)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
