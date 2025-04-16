package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringSliceSet(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want []string
	}{
		{
			name: "empty string",
			got:  "",
			want: []string{},
		},
		{
			name: "simple comma separated",
			got:  "a,b",
			want: []string{"a", "b"},
		},
		{
			name: "multiple commas",
			got:  ",,,",
			want: []string{"", "", "", ""},
		},
		{
			name: "escaped comma",
			got:  ",a\\,",
			want: []string{"", "a,"},
		},
		{
			name: "escaped backslash",
			got:  "a,b\\,c\\\\d,e",
			want: []string{"a", "b,c\\\\d", "e"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dest []string
			s := &StringSlice{
				destination:  &dest,
				separator:    ",",
				escapeString: "\\",
			}

			err := s.Set(tt.got)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, dest)
		})
	}
}

func TestStringSliceString(t *testing.T) {
	tests := []struct {
		name string
		got  []string
		want string
	}{
		{
			name: "empty slice",
			got:  []string{},
			want: "",
		},
		{
			name: "nil slice",
			got:  nil,
			want: "",
		},
		{
			name: "single item",
			got:  []string{"a"},
			want: "a",
		},
		{
			name: "multiple items",
			got:  []string{"a", "b", "c"},
			want: "a,b,c",
		},
		{
			name: "items with commas",
			got:  []string{"a,b", "c"},
			want: "a,b,c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StringSlice{
				destination:  &tt.got,
				separator:    ",",
				escapeString: "\\",
			}

			assert.Equal(t, tt.want, s.String())
		})
	}
}

func TestStringSliceGet(t *testing.T) {
	tests := []struct {
		name string
		got  []string
	}{
		{
			name: "empty slice",
			got:  []string{},
		},
		{
			name: "single item",
			got:  []string{"a"},
		},
		{
			name: "multiple items",
			got:  []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StringSlice{
				destination:  &tt.got,
				separator:    ",",
				escapeString: "\\",
			}

			result := s.Get()
			assert.Equal(t, tt.got, result)
		})
	}
}

func TestStringSliceCreate(t *testing.T) {
	tests := []struct {
		name   string
		got    []string
		config StringSliceConfig
	}{
		{
			name: "default config",
			got:  []string{"a", "b"},
			config: StringSliceConfig{
				Separator:    ",",
				EscapeString: "\\",
			},
		},
		{
			name: "custom config",
			got:  []string{"a", "b"},
			config: StringSliceConfig{
				Separator:    ";",
				EscapeString: "#",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dest []string

			s := StringSlice{}

			value := s.Create(tt.got, &dest, tt.config)

			// Check that destination was set
			assert.Equal(t, tt.got, dest)

			// Check that returned value has correct properties
			stringSlice, ok := value.(*StringSlice)
			assert.True(t, ok)
			assert.Equal(t, &dest, stringSlice.destination)
			assert.Equal(t, tt.config.Separator, stringSlice.separator)
			assert.Equal(t, tt.config.EscapeString, stringSlice.escapeString)
		})
	}
}

func TestStringSliceToString(t *testing.T) {
	tests := []struct {
		name      string
		got       []string
		separator string
		want      string
	}{
		{
			name:      "empty slice",
			got:       []string{},
			separator: ",",
			want:      "",
		},
		{
			name:      "single item",
			got:       []string{"a"},
			separator: ",",
			want:      `"a"`,
		},
		{
			name:      "multiple items",
			got:       []string{"a", "b", "c"},
			separator: ",",
			want:      `"a,b,c"`,
		},
		{
			name:      "custom separator",
			got:       []string{"a", "b", "c"},
			separator: ";",
			want:      `"a;b;c"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StringSlice{separator: tt.separator}

			result := s.ToString(tt.got)
			assert.Equal(t, tt.want, result)
		})
	}
}
