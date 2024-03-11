package plugin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func Test_currFromContext(t *testing.T) {
	tests := []struct {
		name string
		envs map[string]string
		want map[string]string
	}{
		{
			name: "empty commit message",
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
			name: "commit message with title and desc",
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
			name: "commit message with title, desc and additional",
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
			Execute: func(_ context.Context) error { return nil },
		}

		t.Run(tt.name, func(t *testing.T) {
			got := New(options)
			got.App.Action = func(ctx *cli.Context) error {
				got.Metadata = MetadataFromContext(ctx)

				return nil
			}

			_ = got.App.Run([]string{"dummy"})

			assert.Equal(t, got.Metadata.Curr.Message, tt.want["message"])
			assert.Equal(t, got.Metadata.Curr.Title, tt.want["title"])
			assert.Equal(t, got.Metadata.Curr.Description, tt.want["desc"])
		})
	}
}

func TestSplitMessage(t *testing.T) {
	tests := []struct {
		name            string
		message         string
		wantTitle       string
		wantDescription string
	}{
		{
			name:            "empty message",
			message:         "",
			wantTitle:       "",
			wantDescription: "",
		},
		{
			name:            "only title",
			message:         "Title",
			wantTitle:       "Title",
			wantDescription: "",
		},
		{
			name:            "title and description",
			message:         "Title\nDescription",
			wantTitle:       "Title",
			wantDescription: "Description",
		},
		{
			name:            "title and description with blank line",
			message:         "Title\n\nDescription with blank line",
			wantTitle:       "Title",
			wantDescription: "\nDescription with blank line",
		},
		{
			name:            "title and description with multiple blank lines",
			message:         "Title\n\n\nMultiple blank lines\nDescription",
			wantTitle:       "Title",
			wantDescription: "\n\nMultiple blank lines\nDescription",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTitle, gotDescription := splitMessage(tt.message)

			assert.Equal(t, tt.wantTitle, gotTitle)
			assert.Equal(t, tt.wantDescription, gotDescription)
		})
	}
}
