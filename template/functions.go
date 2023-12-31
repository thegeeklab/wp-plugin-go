package template

import (
	"fmt"
	"text/template"
	"unicode"
	"unicode/utf8"

	"github.com/Masterminds/sprig/v3"
)

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

func getLastRune(s string, c int) string {
	j := len(s)
	for i := 0; i < c && j > 0; i++ {
		_, size := utf8.DecodeLastRuneInString(s[:j])
		j -= size
	}

	return s[j:]
}
