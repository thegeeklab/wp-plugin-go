package plugin

import (
	"context"
	"testing"

	"github.com/urfave/cli/v2"
)

func Test_currFromContext(t *testing.T) {
	tests := []struct {
		envs map[string]string
		want map[string]string
	}{
		{
			envs: map[string]string{
				"CI_COMMIT_MESSAGE": "",
			},
			want: map[string]string{
				"title":   "",
				"desc":    "",
				"message": "",
			},
		},
		{
			envs: map[string]string{
				"CI_COMMIT_MESSAGE": "test_title\ntest_desc",
			},
			want: map[string]string{
				"title":   "test_title",
				"desc":    "test_desc",
				"message": "test_title\ntest_desc",
			},
		},
		{
			envs: map[string]string{
				"CI_COMMIT_MESSAGE": "test_title\ntest_desc\nadditional",
			},
			want: map[string]string{
				"title":   "test_title",
				"desc":    "test_desc\nadditional",
				"message": "test_title\ntest_desc\nadditional",
			},
		},
	}

	for _, tt := range tests {
		for key, value := range tt.envs {
			t.Setenv(key, value)
		}

		options := Options{
			Name:    "dummy",
			Execute: func(ctx context.Context) error { return nil },
		}

		got := New(options)
		got.App.Action = func(ctx *cli.Context) error {
			got.Metadata = MetadataFromContext(ctx)

			return nil
		}

		_ = got.App.Run([]string{"dummy"})

		if got.Metadata.Curr.Message != tt.want["message"] {
			t.Errorf("got = %q, want = %q", got.Metadata.Curr.Message, tt.want["message"])
		}

		if got.Metadata.Curr.Title != tt.want["title"] {
			t.Errorf("got = %q, want = %q", got.Metadata.Curr.Title, tt.want["title"])
		}

		if got.Metadata.Curr.Description != tt.want["desc"] {
			t.Errorf("got = %q, want = %q", got.Metadata.Curr.Description, tt.want["desc"])
		}
	}
}
