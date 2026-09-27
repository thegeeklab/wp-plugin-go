package template

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"text/template"
	"unicode"
	"unicode/utf8"

	"github.com/Masterminds/sprig/v3"
)

var errNoPaths = errors.New("sha256file: at least one path is required")

// LoadFuncMap merges the sprig template functions with any custom functions
// provided, giving priority to the custom functions in case of collisions.
func LoadFuncMap() template.FuncMap {
	sprigFuncs := sprig.GenericFuncMap()
	customFuncs := template.FuncMap{
		"toSentence": ToSentence,
		"sha256file": Sha256File,
	}

	for name, f := range customFuncs {
		if _, ok := sprigFuncs[name]; ok {
			continue
		}

		sprigFuncs[name] = f
	}

	return sprigFuncs
}

// ToSentence capitalizes the first letter of the input string,
// adds a period at the end if needed, and returns the resulting sentence.
func ToSentence(s string) string {
	if s == "" {
		return ""
	}

	r, n := utf8.DecodeRuneInString(s)

	closer := ""
	if getLastRune(s, 1) != "." {
		closer = "."
	}

	return fmt.Sprintf("%s%s%s", string(unicode.ToUpper(r)), s[n:], closer)
}

// getLastRune returns the last n runes in the string s.
// It decodes s from the end, counting n runes.
func getLastRune(s string, c int) string {
	j := len(s)
	for i := 0; i < c && j > 0; i++ {
		_, size := utf8.DecodeLastRuneInString(s[:j])
		j -= size
	}

	return s[j:]
}

// Sha256File computes the SHA-256 hash of a file's contents.
// It supports multiple file paths, hashing their contents in order.
func Sha256File(paths ...string) (string, error) {
	if len(paths) == 0 {
		return "", errNoPaths
	}

	h := sha256.New()

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("failed to read %s: %w", path, err)
		}

		h.Write(data)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
