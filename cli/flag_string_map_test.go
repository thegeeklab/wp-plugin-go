package cli

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringMapSet(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  map[string]string
	}{
		{
			name:  "empty string",
			input: "",
			want:  map[string]string{},
		},
		{
			name:  "valid JSON",
			input: `{"key1":"value1","key2":"value2"}`,
			want:  map[string]string{"key1": "value1", "key2": "value2"},
		},
		{
			name:  "single key-value",
			input: `{"key":"value"}`,
			want:  map[string]string{"key": "value"},
		},
		{
			name:  "non-JSON string",
			input: "not-json",
			want:  map[string]string{"*": "not-json"},
		},
		{
			name:  "empty JSON object",
			input: "{}",
			want:  map[string]string{},
		},
		{
			name:  "boolean value",
			input: `{"enabled":true}`,
			want:  map[string]string{"enabled": "true"},
		},
		{
			name:  "numeric values",
			input: `{"count":42,"ratio":3.14}`,
			want:  map[string]string{"count": "42", "ratio": "3.14"},
		},
		{
			name:  "null value",
			input: `{"empty":null}`,
			want:  map[string]string{"empty": ""},
		},
		{
			name:  "null with non-string values",
			input: `{"count":42,"empty":null}`,
			want:  map[string]string{"count": "42", "empty": ""},
		},
		{
			name:  "mixed string, bool, and numeric",
			input: `{"str":"hello","flag":false,"num":7}`,
			want:  map[string]string{"str": "hello", "flag": "false", "num": "7"},
		},
		{
			name:  "top-level JSON array",
			input: `["foo","bar"]`,
			want:  map[string]string{"*": `["foo","bar"]`},
		},
		{
			name:  "top-level JSON string primitive",
			input: `"just a string"`,
			want:  map[string]string{"*": `"just a string"`},
		},
		{
			name:  "top-level JSON number primitive",
			input: `42`,
			want:  map[string]string{"*": "42"},
		},
		{
			name:  "top-level JSON boolean primitive",
			input: `true`,
			want:  map[string]string{"*": "true"},
		},
		{
			name:  "top-level JSON null primitive",
			input: `null`,
			want:  nil,
		},
		{
			name:  "nested object value",
			input: `{"a":{"b":1}}`,
			want:  map[string]string{"a": `{"b":1}`},
		},
		{
			name:  "nested array value",
			input: `{"a":[1,2,3]}`,
			want:  map[string]string{"a": "[1,2,3]"},
		},
		{
			name:  "float edge case",
			input: `{"val":0.30000000000000004}`,
			want:  map[string]string{"val": "0.30000000000000004"},
		},
		{
			name:  "malformed JSON that looks like a map",
			input: `{"foo": bar}`,
			want:  map[string]string{"*": `{"foo": bar}`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dest map[string]string

			s := &StringMap{
				destination: &dest,
			}

			err := s.Set(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, dest)
		})
	}
}

func TestStringMapSetRoundtrip(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "boolean value",
			input: `{"enabled":true}`,
			want:  `{"enabled":"true"}`,
		},
		{
			name:  "numeric values",
			input: `{"count":42,"ratio":3.14}`,
			want:  `{"count":"42","ratio":"3.14"}`,
		},
		{
			name:  "null value",
			input: `{"empty":null}`,
			want:  `{"empty":""}`,
		},
		{
			name:  "mixed types",
			input: `{"str":"hello","flag":false,"num":7}`,
			want:  `{"flag":"false","num":"7","str":"hello"}`,
		},
		{
			name:  "non-JSON fallback",
			input: `not-json`,
			want:  `{"*":"not-json"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dest map[string]string

			s := &StringMap{
				destination: &dest,
			}

			err := s.Set(tt.input)
			assert.NoError(t, err)

			got := s.String()

			var expected, actual map[string]string

			_ = json.Unmarshal([]byte(tt.want), &expected)
			_ = json.Unmarshal([]byte(got), &actual)
			assert.EqualValues(t, expected, actual)
		})
	}
}

func TestStringMapString(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]string
		want  string
	}{
		{
			name:  "empty map",
			input: map[string]string{},
			want:  "",
		},
		{
			name:  "nil map",
			input: nil,
			want:  "",
		},
		{
			name:  "single key-value",
			input: map[string]string{"key": "value"},
			want:  `{"key":"value"}`,
		},
		{
			name:  "multiple key-values",
			input: map[string]string{"key1": "value1", "key2": "value2"},
			want:  `{"key1":"value1","key2":"value2"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StringMap{
				destination: &tt.input,
			}

			got := s.String()

			if len(tt.input) > 1 {
				var expected, actual map[string]string

				_ = json.Unmarshal([]byte(tt.want), &expected)
				_ = json.Unmarshal([]byte(got), &actual)
				assert.EqualValues(t, expected, actual)

				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStringMapGet(t *testing.T) {
	tests := []struct {
		name string
		want map[string]string
	}{
		{
			name: "empty map",
			want: map[string]string{},
		},
		{
			name: "single key-value",
			want: map[string]string{"key": "value"},
		},
		{
			name: "multiple key-values",
			want: map[string]string{"key1": "value1", "key2": "value2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StringMap{
				destination: &tt.want,
			}

			result := s.Get()
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestStringMapCreate(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]string
		want  map[string]string
	}{
		{
			name:  "empty map",
			input: nil,
			want:  map[string]string{},
		},
		{
			name:  "empty map",
			input: map[string]string{},
			want:  map[string]string{},
		},
		{
			name:  "single key-value",
			input: map[string]string{"key": "value"},
			want:  map[string]string{"key": "value"},
		},
		{
			name:  "multiple key-values",
			input: map[string]string{"key1": "value1", "key2": "value2"},
			want:  map[string]string{"key1": "value1", "key2": "value2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dest map[string]string

			s := StringMap{}
			config := StringMapConfig{}

			got := s.Create(tt.input, &dest, config)
			assert.Equal(t, tt.want, dest)
			assert.Equal(t, &dest, got.(*StringMap).destination)
		})
	}
}

func TestStringMapToString(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]string
		want  string
	}{
		{
			name:  "empty map",
			input: map[string]string{},
			want:  "",
		},
		{
			name:  "single key-value",
			input: map[string]string{"key": "value"},
			want:  `{"key":"value"}`,
		},
		{
			name:  "multiple key-values",
			input: map[string]string{"key1": "value1", "key2": "value2"},
			want:  `{"key1":"value1","key2":"value2"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StringMap{}

			got := s.ToString(tt.input)

			if len(tt.input) > 1 {
				var expected, actual map[string]string

				_ = json.Unmarshal([]byte(tt.want), &expected)
				_ = json.Unmarshal([]byte(got), &actual)
				assert.EqualValues(t, expected, actual)

				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
