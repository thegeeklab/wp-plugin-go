package file

import (
	"errors"
	"io"
	"os"
)

// DeleteDir deletes the directory at the given path.
// It returns nil if the deletion succeeds, or the deletion error otherwise.
// If the directory does not exist, DeleteDir returns nil.
func DeleteDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(path)
}

// IsDir returns whether the given path is a directory. If the path does not exist, it returns (false, nil).
// If there is an error checking the path, it returns (false, err).
func IsDir(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// IsDirEmpty returns whether the given directory path is empty.
// If the path does not exist or is not a directory, it returns (false, err).
func IsDirEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdir(1)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, io.EOF) {
		return true, nil
	}

	return false, err
}
