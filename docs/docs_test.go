package docs

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

func testApp() *cli.Command {
	app := &cli.Command{
		Name:        "test",
		Description: "test description",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:     "dummy-flag-int",
				Usage:    "dummy int flag desc",
				Sources:  cli.EnvVars("PLUGIN_DUMMY_FLAG_INT"),
				Required: true,
			},
			&cli.StringFlag{
				Name:     "dummy-flag",
				Usage:    "Dummy flag desc.",
				Sources:  cli.EnvVars("PLUGIN_DUMMY_FLAG"),
				Required: true,
			},
			&cli.StringFlag{
				Name:    "simpe-flag",
				Value:   "simple",
				Sources: cli.EnvVars("PLUGIN_X_SIMPLE_FLAG"),
			},
			&cli.StringFlag{
				Name:    "other.flag",
				Usage:   "other flag with desc",
				Sources: cli.EnvVars("PLUGIN_Z_OTHER_FLAG"),
			},
			&cli.StringSliceFlag{
				Name:    "slice.flag",
				Usage:   "slice flag",
				Sources: cli.EnvVars("PLUGIN_SLICE_FLAG"),
			},
			&cli.StringFlag{
				Name:    "hidden.flag",
				Usage:   "hidden flag",
				Sources: cli.EnvVars("HIDDEN_FLAG", "PLUGIN_HIDDEN_FLAG"),
			},
		},
	}

	return app
}

func testFileContent(t *testing.T, file string) string {
	t.Helper()

	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))

	return string(data)
}

func TestToMarkdown(t *testing.T) {
	want := testFileContent(t, "testdata/expected-doc-no-long.md")

	got, err := ToMarkdown(testApp())
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestToMarkdownWithSource(t *testing.T) {
	tests := []struct {
		name       string
		app        *cli.Command
		sourcePath string
		want       string
	}{
		{
			name:       "normal branch",
			app:        testApp(),
			sourcePath: "testdata/flags.go",
			want:       "testdata/expected-doc-full.md",
		},
		{
			name:       "empty source path skips long descriptions",
			app:        testApp(),
			sourcePath: "",
			want:       "testdata/expected-doc-no-long.md",
		},
		{
			name:       "missing source path skips long descriptions gracefully",
			app:        testApp(),
			sourcePath: "testdata/does-not-exist.go",
			want:       "testdata/expected-doc-no-long.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := testFileContent(t, tt.want)
			got, _ := ToMarkdownWithSource(tt.app, tt.sourcePath)
			assert.Equal(t, want, got)
		})
	}
}

func TestGetTemplateData(t *testing.T) {
	assert.Equal(t, GetTemplateDataWithSource(testApp(), ""), GetTemplateData(testApp()))
}

func TestGetTemplateDataWithSource(t *testing.T) {
	tests := []struct {
		name       string
		app        *cli.Command
		sourcePath string
		want       *CliTemplate
	}{
		{
			name:       "normal branch",
			app:        testApp(),
			sourcePath: "testdata/flags.go",
			want: &CliTemplate{
				Name:        "test",
				Description: "test description",
				GlobalArgs: []*PluginArg{
					{
						Name:        "dummy_flag",
						Description: "Dummy flag desc.",
						LongDescription: "&emsp;Dummy flag long description spanning two source lines " +
							"in the same paragraph.\n\n&emsp;Second paragraph of the dummy flag long description.",
						Type:     "string",
						Required: true,
					},
					{
						Name:        "dummy_flag_int",
						Description: "dummy int flag desc",
						Type:        "integer",
						Required:    true,
					},
					{
						Name:        "slice_flag",
						Description: "slice flag",
						LongDescription: "&emsp;Long description for the slice flag with multiple paragraphs." +
							"\n\n&emsp;Second paragraph for slice flag.",
						Default:  "",
						Type:     "list",
						Required: false,
					},
					{
						Name:     "x_simple_flag",
						Type:     "string",
						Default:  "\"simple\"",
						Required: false,
					},
					{
						Name:        "z_other_flag",
						Description: "other flag with desc",
						Type:        "string",
						Required:    false,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTemplateDataWithSource(tt.app, tt.sourcePath)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLongDescriptionsFor(t *testing.T) {
	got := LongDescriptionsFor("testdata/flags.go")

	want := map[string]*LongDescription{
		"dummy_flag": {Paragraphs: [][]string{
			{"Dummy flag long description spanning", "two source lines in the same paragraph."},
			{"Second paragraph of the dummy flag long description."},
		}},
		"slice_flag": {Paragraphs: [][]string{
			{"Long description for the slice flag with", "multiple paragraphs."},
			{"Second paragraph for slice flag."},
		}},
	}

	assert.Equal(t, want, got)
}

func TestLongDescriptionFunc(t *testing.T) {
	extracted := &LongDescription{Paragraphs: [][]string{{"extracted long desc"}}}
	fallback := func(arg *PluginArg) *LongDescription {
		return &LongDescription{Paragraphs: [][]string{{"fallback for " + arg.Name}}}
	}

	tests := []struct {
		name  string
		descs map[string]*LongDescription
		fb    LongDescriptionFallback
		arg   *PluginArg
		want  *LongDescription
	}{
		{
			name:  "extracted wins over fallback",
			descs: map[string]*LongDescription{"foo": extracted},
			fb:    fallback,
			arg:   &PluginArg{Name: "foo", Description: "short"},
			want:  extracted,
		},
		{
			name:  "zero extracted triggers fallback",
			descs: map[string]*LongDescription{"foo": {}},
			fb:    fallback,
			arg:   &PluginArg{Name: "foo"},
			want:  &LongDescription{Paragraphs: [][]string{{"fallback for foo"}}},
		},
		{
			name:  "missing entry triggers fallback",
			descs: map[string]*LongDescription{},
			fb:    fallback,
			arg:   &PluginArg{Name: "foo"},
			want:  &LongDescription{Paragraphs: [][]string{{"fallback for foo"}}},
		},
		{
			name:  "fallback returning nil propagates nil",
			descs: map[string]*LongDescription{},
			fb:    func(*PluginArg) *LongDescription { return nil },
			arg:   &PluginArg{Name: "foo"},
			want:  nil,
		},
		{
			name:  "missing entry and no fallback returns nil",
			descs: map[string]*LongDescription{},
			fb:    nil,
			arg:   &PluginArg{Name: "foo"},
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LongDescriptionFunc(tt.descs, tt.fb)(tt.arg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestShortDescriptionFallback(t *testing.T) {
	tests := []struct {
		name string
		arg  *PluginArg
		want *LongDescription
	}{
		{
			name: "synthesizes sentence-formatted paragraph",
			arg:  &PluginArg{Name: "foo", Description: "short desc"},
			want: &LongDescription{Paragraphs: [][]string{{"Short desc."}}},
		},
		{
			name: "leaves already-sentence text untouched",
			arg:  &PluginArg{Name: "foo", Description: "Already a sentence."},
			want: &LongDescription{Paragraphs: [][]string{{"Already a sentence."}}},
		},
		{
			name: "empty description returns nil",
			arg:  &PluginArg{Name: "foo"},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ShortDescriptionFallback(tt.arg))
		})
	}
}

func TestBoolFlagDefault(t *testing.T) {
	app := &cli.Command{
		Name: "test",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "bool-flag",
				Usage:   "bool flag desc",
				Sources: cli.EnvVars("PLUGIN_BOOL_FLAG"),
			},
		},
	}

	got := GetTemplateData(app)

	assert.Len(t, got.GlobalArgs, 1)
	assert.Equal(t, "bool_flag", got.GlobalArgs[0].Name)
	assert.Equal(t, "false", got.GlobalArgs[0].Default)
	assert.Equal(t, "bool", got.GlobalArgs[0].Type)
}
