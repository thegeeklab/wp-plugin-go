package template

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToSentence(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "sentence without end period",
			input: "this is a sentence",
			want:  "This is a sentence.",
		},
		{
			name:  "sentence with end period",
			input: "this is a sentence.",
			want:  "This is a sentence.",
		},
		{
			name:  "single word",
			input: "word",
			want:  "Word.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToSentence(tt.input)

			assert.Equal(t, got, tt.want)
		})
	}
}

func TestLoadFuncMap(t *testing.T) {
	tests := []struct {
		name     string
		want     []string
		wantDiff int
	}{
		{
			name: "valid",
			want: []string{
				"toSentence",
				"sha256file",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LoadFuncMap()

			_, ok := got["toSentence"]
			assert.True(t, ok, "LoadFuncMap() missing toSentence func")

			_, ok = got["sha256file"]
			assert.True(t, ok, "LoadFuncMap() missing sha256file func")
		})
	}
}

func TestSha256File(t *testing.T) {
	dir := t.TempDir()

	file1 := filepath.Join(dir, "file1.txt")
	if err := os.WriteFile(file1, []byte("hello"), 0o600); err != nil {
		t.Fatal(err)
	}

	file2 := filepath.Join(dir, "file2.txt")
	if err := os.WriteFile(file2, []byte("world"), 0o600); err != nil {
		t.Fatal(err)
	}

	emptyFile := filepath.Join(dir, "empty.txt")
	if err := os.WriteFile(emptyFile, []byte{}, 0o600); err != nil {
		t.Fatal(err)
	}

	subDir := filepath.Join(dir, "subdir")
	if err := os.Mkdir(subDir, 0o700); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		paths   []string
		want    string
		wantErr bool
	}{
		{
			name:  "single file",
			paths: []string{file1},
			want:  "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
		{
			name:  "multiple files",
			paths: []string{file1, file2},
			want:  "936a185caaa266bb9cbe981e9e05cb78cd732b0b3280eb944412bb6f8f8f07af",
		},
		{
			name:  "empty file",
			paths: []string{emptyFile},
			want:  "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:    "missing file",
			paths:   []string{filepath.Join(dir, "missing.txt")},
			wantErr: true,
		},
		{
			name:    "directory as input",
			paths:   []string{subDir},
			wantErr: true,
		},
		{
			name:    "no paths",
			paths:   []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sha256File(tt.paths...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
