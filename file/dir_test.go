package file

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestIsDirEmpty(t *testing.T) {
	t.Run("empty directory", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "test")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(dir)

		isEmpty, err := IsDirEmpty(dir)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !isEmpty {
			t.Error("expected directory to be empty")
		}
	})

	t.Run("non-empty directory", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "test")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(dir)

		file, err := os.CreateTemp(dir, "test")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}

		file.Close()

		isEmpty, err := IsDirEmpty(dir)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if isEmpty {
			t.Error("expected directory to be non-empty")
		}
	})

	t.Run("non-existent directory", func(t *testing.T) {
		dir := filepath.Join(os.TempDir(), "non-existent")

		isEmpty, err := IsDirEmpty(dir)
		if err == nil {
			t.Error("expected an error for non-existent directory")
		}

		if isEmpty {
			t.Error("expected directory to be non-empty")
		}

		if !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
