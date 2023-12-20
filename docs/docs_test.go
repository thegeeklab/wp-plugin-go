package docs

import (
	"bytes"
	"os"
	"testing"

	"github.com/urfave/cli/v2"
)

func testApp() *cli.App {
	app := &cli.App{
		Name:        "test",
		Description: "test description",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "dummy-flag",
				Usage:   "dummy flag desc",
				EnvVars: []string{"PLUGIN_DUMMY_FLAG"},
				Value:   "test",
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
		{"normal branch", testApp(), "testdata/expected-doc-full.md"},
	}

	for _, tt := range tests {
		want := testFileContent(t, tt.want)
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := ToMarkdown(tt.app); got != want {
				t.Errorf("got = %v, want %v", got, want)
			}
		})
	}
}
