package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetDifference(t *testing.T) {
	tests := []struct {
		name   string
		a      []string
		b      []string
		want   []string
		unique bool
	}{
		{
			name:   "both empty",
			a:      []string{},
			b:      []string{},
			want:   []string{},
			unique: false,
		},
		{
			name:   "remove common element",
			a:      []string{"a", "b", "c"},
			b:      []string{"b"},
			want:   []string{"a", "c"},
			unique: false,
		},
		{
			name:   "remove a and c",
			a:      []string{"a", "b", "c"},
			b:      []string{"a", "c"},
			want:   []string{"b"},
			unique: false,
		},
		{
			name:   "no common elements",
			a:      []string{"a", "b", "c"},
			b:      []string{"d", "e"},
			want:   []string{"a", "b", "c"},
			unique: false,
		},
		{
			name:   "remove duplicates",
			a:      []string{"a", "a", "b"},
			b:      []string{"b"},
			want:   []string{"a"},
			unique: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SetDifference(tt.a, tt.b, tt.unique)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestUnique(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{
			name:  "empty slice",
			input: []int{},
			want:  []int{},
		},
		{
			name:  "single value",
			input: []int{1},
			want:  []int{1},
		},
		{
			name:  "duplicates",
			input: []int{1, 2, 3, 2},
			want:  []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Unique(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
