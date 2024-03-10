package template

import (
	"fmt"
	"text/template"
	"unicode"
	"unicode/utf8"

	"github.com/Masterminds/sprig/v3"
)

// LoadFuncMap merges the sprig template functions with any custom functions
// provided, giving priority to the custom functions in case of collisions.
func LoadFuncMap() template.FuncMap {
	sprigFuncs := sprig.GenericFuncMap()
	customFuncs := template.FuncMap{
		"ToSentence": ToSentence,
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
