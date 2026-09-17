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
		t.Error(err)
	}

	data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))

	return string(data)
}

func TestToMarkdownFull(t *testing.T) {
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
			got, _ := ToMarkdown(tt.app, tt.sourcePath)
			assert.Equal(t, want, got)
		})
	}
}

func TestToData(t *testing.T) {
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
			got := GetTemplateData(tt.app, tt.sourcePath)
			assert.Equal(t, tt.want, got)
		})
	}
}
