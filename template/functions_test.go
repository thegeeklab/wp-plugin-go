package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToSentence(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "sentence without end period",
			input: "this is a sentence",
			want:  "This is a sentence.",
		},
		{
			name:  "sentence with end period",
			input: "this is a sentence.",
			want:  "This is a sentence.",
		},
		{
			name:  "single word",
			input: "word",
			want:  "Word.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToSentence(tt.input)

			assert.Equal(t, got, tt.want)
		})
	}
}

func TestLoadFuncMap(t *testing.T) {
	tests := []struct {
		name     string
		want     []string
		wantDiff int
	}{
		{
			name: "valid",
			want: []string{
				"ToSentence",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LoadFuncMap()

			_, ok := got["ToSentence"]
			assert.True(t, ok, "LoadFuncMap() missing ToSentence func")
		})
	}
}
