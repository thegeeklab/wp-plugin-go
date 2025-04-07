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
				Value:    10,
				Required: true,
			},
			&cli.StringFlag{
				Name:     "dummy-flag",
				Usage:    "Dummy flag desc.",
				Sources:  cli.EnvVars("PLUGIN_DUMMY_FLAG"),
				Value:    "test",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "simpe-flag",
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
		name string
		app  *cli.Command
		want string
	}{
		{
			"normal branch",
			testApp(),
			"testdata/expected-doc-full.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := testFileContent(t, tt.want)
			got, _ := ToMarkdown(tt.app)
			assert.Equal(t, want, got)
		})
	}
}

func TestToData(t *testing.T) {
	tests := []struct {
		name string
		app  *cli.Command
		want *CliTemplate
	}{
		{
			name: "normal branch",
			app:  testApp(),
			want: &CliTemplate{
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
		t.Run(tt.name, func(t *testing.T) {
			got := GetTemplateData(tt.app)
			assert.Equal(t, tt.want, got)
		})
	}
}
