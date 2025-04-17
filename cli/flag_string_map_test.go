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
