package file

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDirEmpty(t *testing.T) {
	tests := []struct {
		name       string
		dir        string
		createFile bool
		want       bool
		wantErr    error
	}{
		{
			name: "empty directory",
			want: true,
		},
		{
			name:       "non-empty directory",
			createFile: true,
			want:       false,
		},
		{
			name:    "non-existent directory",
			dir:     filepath.Join(os.TempDir(), "non-existent"),
			want:    false,
			wantErr: fs.ErrNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.dir
			if dir == "" {
				dir = t.TempDir()
			}

			if tt.createFile {
				file, err := os.CreateTemp(dir, "test")
				if err != nil {
					t.Fatal(err)
				}

				file.Close()
			}

			isEmpty, err := IsDirEmpty(dir)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.False(t, isEmpty)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, isEmpty)
		})
	}
}

func TestIsDir(t *testing.T) {
	dir := t.TempDir()

	file, err := os.CreateTemp(dir, "test")
	if err != nil {
		t.Fatal(err)
	}

	file.Close()

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "directory",
			path: dir,
			want: true,
		},
		{
			name: "file",
			path: file.Name(),
			want: false,
		},
		{
			name: "non-existent path",
			path: filepath.Join(dir, "non-existent"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsDir(tt.path)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDeleteDir(t *testing.T) {
	dir := t.TempDir()

	existingDir := filepath.Join(dir, "existing")
	if err := os.Mkdir(existingDir, 0o750); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		path string
	}{
		{
			name: "delete existing directory",
			path: existingDir,
		},
		{
			name: "delete non-existent directory",
			path: filepath.Join(dir, "non-existent"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DeleteDir(tt.path)
			assert.NoError(t, err)

			_, statErr := os.Stat(tt.path)
			assert.ErrorIs(t, statErr, fs.ErrNotExist)
		})
	}
}
