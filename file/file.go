package file

import (
	"fmt"
	"os"
	"path/filepath"
)

// The MSDN docs appear to say that a normal path that is 248 bytes long will work;
// empirically the path must be less then 248 bytes long.
// See https://learn.microsoft.com/en-us/windows/win32/fileio/naming-a-file#maximum-path-length-limitation
const maxPathLenght = 248

// ReadStringOrFile returns the content of a string or a file if the file exists.
// The returned boolean value indicates whether the specified `input` was a file path or not.
func ReadStringOrFile(input string) (string, bool, error) {
	if len(input) >= maxPathLenght {
		return input, false, nil
	}

	// Check if input is a file path
	if _, err := os.Stat(input); err != nil && os.IsNotExist(err) {
		// No file found => use input as result
		return input, false, nil
	} else if err != nil {
		return "", false, err
	}

	result, err := os.ReadFile(input)
	if err != nil {
		return "", true, err
	}

	return string(result), true, nil
}

// DeleteDir deletes the directory at the given path.
// It returns nil if the deletion succeeds, or the deletion error otherwise.
// If the directory does not exist, DeleteDir returns nil.
func DeleteDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(path)
}

// ExpandFileList takes a list of file globs and expands them into a list
// of matching file paths. It returns the expanded file list and any errors
// from glob matching. This allows safely passing user input globs through to
// glob matching.
func ExpandFileList(fileList []string) ([]string, error) {
	var result []string

	for _, glob := range fileList {
		globbed, err := filepath.Glob(glob)
		if err != nil {
			return result, fmt.Errorf("failed to match %s: %w", glob, err)
		}

		if globbed != nil {
			result = append(result, globbed...)
		}
	}

	return result, nil
}
