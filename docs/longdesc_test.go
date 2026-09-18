package docs

import (
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"strings"
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
		want    map[string]*LongDescription
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
			want: map[string]*LongDescription{
				"foo": {Paragraphs: [][]string{{"Short help, short comment."}}},
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
			want: map[string]*LongDescription{
				"bar": {Paragraphs: [][]string{
					{"Long help, split across source lines", "within the same paragraph."},
				}},
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
			want: map[string]*LongDescription{
				"baz": {Paragraphs: [][]string{
					{"Long help with multiple paragraphs."},
					{"The second paragraph has multiple lines", "spanning two source lines."},
				}},
			},
		},
		{
			name: "bullet list within a paragraph keeps continuation indent",
			body: `// Supported actions:
		//
		// - **upload:** does the thing.
		//   More on upload.
		// - **download:** other thing.
		&cli.StringSliceFlag{
			Name:    "qux",
			Usage:   "qux usage",
			Sources: cli.EnvVars("PLUGIN_QUX"),
		},`,
			want: map[string]*LongDescription{
				"qux": {Paragraphs: [][]string{
					{"Supported actions:"},
					{
						"- **upload:** does the thing.",
						"  More on upload.",
						"- **download:** other thing.",
					},
				}},
			},
		},
		{
			name: "no comment, flag is skipped",
			body: `&cli.StringSliceFlag{
			Name:    "no-comment",
			Usage:   "no comment usage",
			Sources: cli.EnvVars("PLUGIN_NO_COMMENT"),
		},`,
			want: map[string]*LongDescription{},
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
			want: map[string]*LongDescription{},
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

			got, err := LongDescriptionsWith(path)
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

func TestLongDescriptionsWithCustomMatcher(t *testing.T) {
	const src = `package flags

import plugin_cli "github.com/thegeeklab/wp-plugin-go/v6/cli"
import "github.com/urfave/cli/v3"

// Custom map flag explanation, with two paragraphs.
//
// Second paragraph.
var custom = &plugin_cli.StringMapFlag{
	Name:    "custom",
	Usage:   "custom usage",
	Sources: cli.EnvVars("PLUGIN_CUSTOM"),
}

// Core flag explanation.
var core = &cli.StringFlag{
	Name:    "core",
	Usage:   "core usage",
	Sources: cli.EnvVars("PLUGIN_CORE"),
}
`

	dir := t.TempDir()
	path := filepath.Join(dir, "flags.go")

	if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		matchers []FlagTypeMatcher
		want     map[string]*LongDescription
	}{
		{
			name: "default matcher ignores custom flag",
			want: map[string]*LongDescription{
				"core": {Paragraphs: [][]string{{"Core flag explanation."}}},
			},
		},
		{
			name:     "custom matcher walks custom flag only",
			matchers: []FlagTypeMatcher{SelectorMatcher("plugin_cli", "StringMapFlag")},
			want: map[string]*LongDescription{
				"custom": {Paragraphs: [][]string{
					{"Custom map flag explanation, with two paragraphs."},
					{"Second paragraph."},
				}},
			},
		},
		{
			name: "matchers are OR'd; core and custom both returned",
			matchers: []FlagTypeMatcher{
				DefaultFlagTypeMatcher,
				SelectorMatcher("plugin_cli", "StringMapFlag"),
			},
			want: map[string]*LongDescription{
				"core": {Paragraphs: [][]string{{"Core flag explanation."}}},
				"custom": {Paragraphs: [][]string{
					{"Custom map flag explanation, with two paragraphs."},
					{"Second paragraph."},
				}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LongDescriptionsWith(path, tt.matchers...)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSelectorMatcher(t *testing.T) {
	named := SelectorMatcher("pkg", "Foo", "Bar")

	tests := []struct {
		name    string
		matcher FlagTypeMatcher
		expr    ast.Expr
		want    bool
	}{
		{
			name:    "non selector",
			matcher: named,
			expr:    &ast.Ident{Name: "Foo"},
			want:    false,
		},
		{
			name:    "wrong package",
			matcher: named,
			expr:    makeSelector("other", "Foo"),
			want:    false,
		},
		{
			name:    "unlisted selector",
			matcher: named,
			expr:    makeSelector("pkg", "Baz"),
			want:    false,
		},
		{
			name:    "Foo",
			matcher: named,
			expr:    makeSelector("pkg", "Foo"),
			want:    true,
		},
		{
			name:    "Bar",
			matcher: named,
			expr:    makeSelector("pkg", "Bar"),
			want:    true,
		},
		{
			name:    "empty name list matches nothing",
			matcher: SelectorMatcher("pkg"),
			expr:    makeSelector("pkg", "Foo"),
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.matcher(tt.expr))
		})
	}
}

func makeSelector(pkg, name string) ast.Expr {
	return &ast.SelectorExpr{
		X:   &ast.Ident{Name: pkg},
		Sel: &ast.Ident{Name: name},
	}
}

func TestLongDescriptionFlat(t *testing.T) {
	tests := []struct {
		name string
		desc *LongDescription
		want string
	}{
		{
			name: "nil",
			desc: nil,
			want: "",
		},
		{
			name: "zero value",
			desc: &LongDescription{},
			want: "",
		},
		{
			name: "single line",
			desc: &LongDescription{Paragraphs: [][]string{{"Hello world."}}},
			want: "Hello world.",
		},
		{
			name: "multi-line paragraph",
			desc: &LongDescription{Paragraphs: [][]string{
				{"Line one", "line two"},
			}},
			want: "Line one\nline two",
		},
		{
			name: "multi-paragraph",
			desc: &LongDescription{Paragraphs: [][]string{
				{"First."},
				{"Second line one", "second line two"},
			}},
			want: "First.\n\nSecond line one\nsecond line two",
		},
		{
			name: "preserves continuation indent",
			desc: &LongDescription{Paragraphs: [][]string{
				{"- item", "  continuation"},
			}},
			want: "- item\n  continuation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.desc.Flat())
			assert.Equal(t, tt.want, tt.desc.String())
		})
	}
}

func TestLongDescriptionMarkdown(t *testing.T) {
	tests := []struct {
		name string
		desc *LongDescription
		want string
	}{
		{
			name: "nil",
			desc: nil,
			want: "",
		},
		{
			name: "single line",
			desc: &LongDescription{Paragraphs: [][]string{{"Hello world."}}},
			want: "&emsp;Hello world.",
		},
		{
			name: "multi-line collapses to single line per paragraph",
			desc: &LongDescription{Paragraphs: [][]string{
				{"Line one", "line two"},
			}},
			want: "&emsp;Line one line two",
		},
		{
			name: "multi-paragraph separated by blank line",
			desc: &LongDescription{Paragraphs: [][]string{
				{"First."},
				{"Second."},
			}},
			want: "&emsp;First.\n\n&emsp;Second.",
		},
		{
			name: "preserves list indentation inside the paragraph",
			desc: &LongDescription{Paragraphs: [][]string{
				{"- item", "  continuation"},
			}},
			want: "&emsp;- item   continuation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, LongDescriptionMarkdown(tt.desc))
		})
	}
}

func TestLongDescriptionYAMLBlock(t *testing.T) {
	tests := []struct {
		name   string
		desc   *LongDescription
		indent string
		want   string
	}{
		{
			name:   "nil description",
			desc:   nil,
			indent: "  ",
			want:   "",
		},
		{
			name:   "zero value",
			desc:   &LongDescription{},
			indent: "  ",
			want:   "",
		},
		{
			name:   "single line",
			desc:   &LongDescription{Paragraphs: [][]string{{"Hello world."}}},
			indent: "  ",
			want:   "  Hello world.",
		},
		{
			name: "multi-line paragraph",
			desc: &LongDescription{Paragraphs: [][]string{
				{"Line one", "line two"},
			}},
			indent: "    ",
			want:   "    Line one\n    line two",
		},
		{
			name: "multi-paragraph with blank separator",
			desc: &LongDescription{Paragraphs: [][]string{
				{"First."},
				{"Second line.", "second line continued"},
			}},
			indent: "      ",
			want:   "      First.\n\n      Second line.\n      second line continued",
		},
		{
			name: "preserves list continuation indent",
			desc: &LongDescription{Paragraphs: [][]string{
				{"- item", "  continuation"},
			}},
			indent: "    ",
			want:   "    - item\n      continuation",
		},
		{
			name:   "empty indent",
			desc:   &LongDescription{Paragraphs: [][]string{{"x"}}},
			indent: "",
			want:   "x",
		},
		{
			name: "blank-line paragraph separator has no trailing whitespace",
			desc: &LongDescription{Paragraphs: [][]string{
				{"A"},
				{"B"},
			}},
			indent: "      ",
			want:   "      A\n\n      B",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, LongDescriptionYAMLBlock(tt.desc, tt.indent))

			if tt.desc == nil || len(tt.desc.Paragraphs) == 0 {
				return
			}

			out := LongDescriptionYAMLBlock(tt.desc, tt.indent)
			for i, line := range strings.Split(out, "\n") {
				if line == "" {
					assert.Empty(t, line, "line %d has trailing whitespace", i)

					continue
				}

				if strings.HasSuffix(line, " ") {
					t.Errorf("line %d has trailing whitespace: %q", i, line)
				}
			}
		})
	}
}
