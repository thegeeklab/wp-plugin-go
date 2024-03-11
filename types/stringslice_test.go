package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitWithEscaping(t *testing.T) {
	tests := []struct {
		input  string
		output []string
	}{
		{input: "", output: []string{}},
		{input: "a,b", output: []string{"a", "b"}},
		{input: ",,,", output: []string{"", "", "", ""}},
		{input: ",a\\,", output: []string{"", "a,"}},
		{input: "a,b\\,c\\\\d,e", output: []string{"a", "b,c\\\\d", "e"}},
	}

	for _, tt := range tests {
		strings := splitWithEscaping(tt.input, ",", "\\")
		got, want := strings, tt.output

		assert.Equal(t, got, want)
	}
}
