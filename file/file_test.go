package file

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const helloWorld = "Hello, World!"

func TestWriteTmpFile(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		content  string
		wantErr  error
	}{
		{
			name:     "write to temp file",
			fileName: "test.txt",
			content:  helloWorld,
		},
		{
			name:     "empty file name",
			fileName: "",
			content:  helloWorld,
		},
		{
			name:     "empty file content",
			fileName: "test.txt",
			content:  "",
		},
		{
			name:     "create temp file error",
			fileName: filepath.Join(os.TempDir(), "non-existent", "test.txt"),
			content:  helloWorld,
			wantErr:  assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := WriteTmpFile(tt.fileName, tt.content)
			if tt.wantErr != nil {
				assert.Error(t, err)

				return
			}

			assert.NoError(t, err)

			defer os.Remove(tmpFile)

			data, err := os.ReadFile(tmpFile)
			assert.NoError(t, err)
			assert.Equal(t, tt.content, string(data))
		})
	}
}

func TestReadStringOrFile(t *testing.T) {
	longInput := strings.Repeat("a", maxPathLenght)

	tests := []struct {
		name        string
		input       string
		fileContent string
		useTempDir  bool
		want        string
		wantIsFile  bool
		wantErr     error
	}{
		{
			name:  "plain string",
			input: "hello",
			want:  "hello",
		},
		{
			name:  "string at max path length",
			input: longInput,
			want:  longInput,
		},
		{
			name:  "non-existent path treated as string",
			input: filepath.Join(os.TempDir(), "no-such-file"),
			want:  filepath.Join(os.TempDir(), "no-such-file"),
		},
		{
			name:        "existing file",
			fileContent: "file content",
			want:        "file content",
			wantIsFile:  true,
		},
		{
			name:       "directory path",
			useTempDir: true,
			want:       "",
			wantIsFile: true,
			wantErr:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.input

			if tt.fileContent != "" {
				input = filepath.Join(t.TempDir(), "content.txt")
				if err := os.WriteFile(input, []byte(tt.fileContent), 0o600); err != nil {
					t.Fatal(err)
				}
			}

			if tt.useTempDir {
				input = t.TempDir()
			}

			got, isFile, err := ReadStringOrFile(input)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.want, got)
				assert.Equal(t, tt.wantIsFile, isFile)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantIsFile, isFile)
		})
	}
}

func TestExpandFileList(t *testing.T) {
	dir := t.TempDir()

	fileA := filepath.Join(dir, "a.txt")
	fileB := filepath.Join(dir, "b.txt")
	fileC := filepath.Join(dir, "c.md")

	for _, file := range []string{fileA, fileB, fileC} {
		if err := os.WriteFile(file, []byte("content"), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name    string
		files   []string
		want    []string
		wantErr error
	}{
		{
			name:  "empty list",
			files: []string{},
		},
		{
			name:  "no matches",
			files: []string{filepath.Join(dir, "no-such-*.txt")},
		},
		{
			name:  "single glob match",
			files: []string{filepath.Join(dir, "*.txt")},
			want:  []string{fileA, fileB},
		},
		{
			name:  "multiple globs",
			files: []string{filepath.Join(dir, "*.txt"), filepath.Join(dir, "*.md")},
			want:  []string{fileA, fileB, fileC},
		},
		{
			name:    "invalid glob pattern",
			files:   []string{"["},
			wantErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExpandFileList(tt.files)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Nil(t, got)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
