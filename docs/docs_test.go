package docs

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func testApp() *cli.App {
	app := &cli.App{
		Name:        "test",
		Description: "test description",
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:     "dummy-flag-int",
				Usage:    "dummy int flag desc",
				EnvVars:  []string{"PLUGIN_DUMMY_FLAG_INT"},
				Value:    10,
				Required: true,
			},
			&cli.StringFlag{
				Name:     "dummy-flag",
				Usage:    "Dummy flag desc.",
				EnvVars:  []string{"PLUGIN_DUMMY_FLAG"},
				Value:    "test",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "simpe-flag",
				EnvVars: []string{"PLUGIN_X_SIMPLE_FLAG"},
			},
			&cli.StringFlag{
				Name:    "other.flag",
				Usage:   "other flag with desc",
				EnvVars: []string{"PLUGIN_Z_OTHER_FLAG"},
			},
			&cli.StringSliceFlag{
				Name:    "slice.flag",
				Usage:   "slice flag",
				EnvVars: []string{"PLUGIN_SLICE_FLAG"},
			},
			&cli.StringFlag{
				Name:    "hidden.flag",
				Usage:   "hidden flag",
				EnvVars: []string{"HIDDEN_FLAG", "PLUGIN_HIDDEN_FLAG"},
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
		name string
		app  *cli.App
		want string
	}{
		{
			"normal branch",
			testApp(),
			"testdata/expected-doc-full.md",
		},
	}

	for _, tt := range tests {
		want := testFileContent(t, tt.want)

		t.Run(tt.name, func(t *testing.T) {
			got, _ := ToMarkdown(tt.app)
			assert.Equal(t, want, got)
		})
	}
}

func TestToData(t *testing.T) {
	tests := []struct {
		name string
		app  *cli.App
		want *CliTemplate
	}{
		{
			"normal branch",
			testApp(),
			&CliTemplate{
				Name:        "test",
				Description: "test description",
				GlobalArgs: []*PluginArg{
					{
						Name:        "dummy_flag",
						Description: "Dummy flag desc.",
						Default:     "\"test\"",
						Type:        "string",
						Required:    true,
					},
					{
						Name:        "dummy_flag_int",
						Description: "dummy int flag desc",
						Default:     "10",
						Type:        "integer",
						Required:    true,
					},
					{
						Name:        "slice_flag",
						Description: "slice flag",
						Default:     "",
						Type:        "list",
						Required:    false,
					},
					{
						Name:     "x_simple_flag",
						Type:     "string",
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
		got := GetTemplateData(tt.app)

		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, got)
		})
	}
}
