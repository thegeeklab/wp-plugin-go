package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitWithEscaping(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output []string
	}{
		{name: "empty string", input: "", output: []string{}},
		{name: "simple comma separated", input: "a,b", output: []string{"a", "b"}},
		{name: "multiple commas", input: ",,,", output: []string{"", "", "", ""}},
		{name: "escaped comma", input: ",a\\,", output: []string{"", "a,"}},
		{name: "escaped backslash", input: "a,b\\,c\\\\d,e", output: []string{"a", "b,c\\\\d", "e"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strings := splitWithEscaping(tt.input, ",", "\\")
			got, want := strings, tt.output

			assert.Equal(t, got, want)
		})
	}
}
