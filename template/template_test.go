package template

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddPrefix(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		input  string
		want   string
	}{
		{
			name:   "empty input",
			prefix: "pre",
			input:  "",
			want:   "",
		},
		{
			name:   "input already has prefix",
			prefix: "pre",
			input:  "pre-existing",
			want:   "pre-existing",
		},
		{
			name:   "add prefix",
			prefix: "pre",
			input:  "-existing",
			want:   "pre-existing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddPrefix(tt.prefix, tt.input)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRender(t *testing.T) {
	dir := t.TempDir()

	tmplFile := filepath.Join(dir, "template.tpl")
	if err := os.WriteFile(tmplFile, []byte("Hello {{.Name}}"), 0o600); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		tmpl    string
		payload any
		want    string
		wantErr error
	}{
		{
			name:    "inline template",
			tmpl:    "Hello {{.Name}}",
			payload: map[string]string{"Name": "World"},
			want:    "Hello World",
		},
		{
			name:    "template with function",
			tmpl:    `{{.Word | ToSentence}}`,
			payload: map[string]string{"Word": "hello"},
			want:    "Hello.",
		},
		{
			name:    "template from file",
			tmpl:    "file://" + tmplFile,
			payload: map[string]string{"Name": "World"},
			want:    "Hello World",
		},
		{
			name:    "missing file",
			tmpl:    "file:///no/such/file.tpl",
			wantErr: assert.AnError,
		},
		{
			name:    "invalid template",
			tmpl:    "{{.Name",
			wantErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Render(context.Background(), http.Client{}, tt.tmpl, tt.payload)
			if tt.wantErr != nil {
				assert.Error(t, err)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRenderTrim(t *testing.T) {
	tests := []struct {
		name    string
		tmpl    string
		payload any
		want    string
	}{
		{
			name:    "trims leading and trailing whitespace",
			tmpl:    "\n  hello {{.Word}}  \n",
			payload: map[string]string{"Word": "world"},
			want:    "hello world",
		},
		{
			name:    "nothing to trim",
			tmpl:    "hello {{.Word}}",
			payload: map[string]string{"Word": "world"},
			want:    "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RenderTrim(context.Background(), http.Client{}, tt.tmpl, tt.payload)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRenderHTTP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "Hello {{.Name}}")
	}))
	defer server.Close()

	closedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "unused")
	}))
	closedServer.Close()

	tests := []struct {
		name    string
		tmpl    string
		payload any
		want    string
		wantErr error
	}{
		{
			name:    "template from http",
			tmpl:    server.URL,
			payload: map[string]string{"Name": "World"},
			want:    "Hello World",
		},
		{
			name:    "unreachable host",
			tmpl:    closedServer.URL,
			payload: map[string]string{"Name": "World"},
			wantErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Render(context.Background(), http.Client{}, tt.tmpl, tt.payload)
			if tt.wantErr != nil {
				assert.Error(t, err)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
