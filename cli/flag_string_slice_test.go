package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringSliceSet(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "empty string",
			input: "",
			want:  []string{},
		},
		{
			name:  "simple comma separated",
			input: "a,b",
			want:  []string{"a", "b"},
		},
		{
			name:  "multiple commas",
			input: ",,,",
			want:  []string{"", "", "", ""},
		},
		{
			name:  "escaped comma",
			input: ",a\\,",
			want:  []string{"", "a,"},
		},
		{
			name:  "escaped backslash",
			input: "a,b\\,c\\\\d,e",
			want:  []string{"a", "b,c\\\\d", "e"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []string
			s := &StringSlice{
				destination:  &got,
				delimiter:    ",",
				escapeString: "\\",
			}

			err := s.Set(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStringSliceString(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{
			name:  "empty slice",
			input: []string{},
			want:  "",
		},
		{
			name:  "nil slice",
			input: nil,
			want:  "",
		},
		{
			name:  "single item",
			input: []string{"a"},
			want:  "a",
		},
		{
			name:  "multiple items",
			input: []string{"a", "b", "c"},
			want:  "a,b,c",
		},
		{
			name:  "items with commas",
			input: []string{"a,b", "c"},
			want:  "a,b,c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StringSlice{
				destination:  &tt.input,
				delimiter:    ",",
				escapeString: "\\",
			}

			assert.Equal(t, tt.want, s.String())
		})
	}
}

func TestStringSliceGet(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "empty slice",
			want: []string{},
		},
		{
			name: "single item",
			want: []string{"a"},
		},
		{
			name: "multiple items",
			want: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StringSlice{
				destination:  &tt.want,
				delimiter:    ",",
				escapeString: "\\",
			}

			result := s.Get()
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestStringSliceCreate(t *testing.T) {
	tests := []struct {
		name   string
		input  []string
		want   []string
		config StringSliceConfig
	}{
		{
			name:  "empty slice",
			input: nil,
			want:  []string{},
		},
		{
			name:  "default config",
			input: []string{"a", "b"},
			want:  []string{"a", "b"},
			config: StringSliceConfig{
				Delimiter:    ",",
				EscapeString: "\\",
			},
		},
		{
			name:  "custom config",
			input: []string{"a", "b"},
			want:  []string{"a", "b"},
			config: StringSliceConfig{
				Delimiter:    ";",
				EscapeString: "#",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dest []string

			s := StringSlice{}
			got := s.Create(tt.input, &dest, tt.config)

			assert.Equal(t, tt.input, dest)
			assert.Equal(t, &dest, got.(*StringSlice).destination)
			assert.Equal(t, tt.config.Delimiter, got.(*StringSlice).delimiter)
			assert.Equal(t, tt.config.EscapeString, got.(*StringSlice).escapeString)
		})
	}
}

func TestStringSliceToString(t *testing.T) {
	tests := []struct {
		name      string
		input     []string
		delimiter string
		want      string
	}{
		{
			name:      "empty slice",
			input:     []string{},
			delimiter: ",",
			want:      "",
		},
		{
			name:      "single item",
			input:     []string{"a"},
			delimiter: ",",
			want:      `"a"`,
		},
		{
			name:      "multiple items",
			input:     []string{"a", "b", "c"},
			delimiter: ",",
			want:      `"a,b,c"`,
		},
		{
			name:      "custom delimiter",
			input:     []string{"a", "b", "c"},
			delimiter: ";",
			want:      `"a;b;c"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StringSlice{delimiter: tt.delimiter}

			got := s.ToString(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
