package file

import "os"

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

func DeleteDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(path)
}
