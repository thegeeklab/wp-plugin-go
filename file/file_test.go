package file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const helloWorld = "Hello, World!"

func TestWriteTmpFile(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		content  string
		wantErr  bool
	}{
		{
			name:     "write to temp file",
			fileName: "test.txt",
			content:  helloWorld,
			wantErr:  false,
		},
		{
			name:     "empty file name",
			fileName: "",
			content:  helloWorld,
			wantErr:  false,
		},
		{
			name:     "empty file content",
			fileName: "test.txt",
			content:  "",
			wantErr:  false,
		},
		{
			name:     "create temp file error",
			fileName: filepath.Join(os.TempDir(), "non-existent", "test.txt"),
			content:  helloWorld,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := WriteTmpFile(tt.fileName, tt.content)
			if tt.wantErr {
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
