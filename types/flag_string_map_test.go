package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringMapSet(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want map[string]string
	}{
		{
			name: "empty string",
			got:  "",
			want: map[string]string{},
		},
		{
			name: "valid JSON",
			got:  `{"key1":"value1","key2":"value2"}`,
			want: map[string]string{"key1": "value1", "key2": "value2"},
		},
		{
			name: "single key-value",
			got:  `{"key":"value"}`,
			want: map[string]string{"key": "value"},
		},
		{
			name: "non-JSON string",
			got:  "not-json",
			want: map[string]string{"*": "not-json"},
		},
		{
			name: "empty JSON object",
			got:  "{}",
			want: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dest map[string]string
			s := &StringMap{
				destination: &dest,
			}

			err := s.Set(tt.got)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, dest)
		})
	}
}

func TestStringMapString(t *testing.T) {
	tests := []struct {
		name string
		got  map[string]string
		want string
	}{
		{
			name: "empty map",
			got:  map[string]string{},
			want: "",
		},
		{
			name: "nil map",
			got:  nil,
			want: "",
		},
		{
			name: "single key-value",
			got:  map[string]string{"key": "value"},
			want: `{"key":"value"}`,
		},
		{
			name: "multiple key-values",
			got:  map[string]string{"key1": "value1", "key2": "value2"},
			want: `{"key1":"value1","key2":"value2"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StringMap{
				destination: &tt.got,
			}

			result := s.String()

			if len(tt.got) > 1 {
				var expected, actual map[string]string
				_ = json.Unmarshal([]byte(tt.want), &expected)
				_ = json.Unmarshal([]byte(result), &actual)
				assert.EqualValues(t, expected, actual)

				return
			}

			assert.Equal(t, tt.want, result)
		})
	}
}

func TestStringMapGet(t *testing.T) {
	tests := []struct {
		name string
		got  map[string]string
	}{
		{
			name: "empty map",
			got:  map[string]string{},
		},
		{
			name: "single key-value",
			got:  map[string]string{"key": "value"},
		},
		{
			name: "multiple key-values",
			got:  map[string]string{"key1": "value1", "key2": "value2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StringMap{
				destination: &tt.got,
			}

			result := s.Get()
			assert.Equal(t, tt.got, result)
		})
	}
}

func TestStringMapCreate(t *testing.T) {
	tests := []struct {
		name string
		got  map[string]string
	}{
		{
			name: "empty map",
			got:  map[string]string{},
		},
		{
			name: "single key-value",
			got:  map[string]string{"key": "value"},
		},
		{
			name: "multiple key-values",
			got:  map[string]string{"key1": "value1", "key2": "value2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dest map[string]string

			s := StringMap{}
			config := StringMapConfig{}

			value := s.Create(tt.got, &dest, config)
			assert.Equal(t, tt.got, dest)
			assert.Equal(t, &dest, value.(*StringMap).destination)
		})
	}
}

//nolint:dupl
func TestStringMapToString(t *testing.T) {
	tests := []struct {
		name string
		got  map[string]string
		want string
	}{
		{
			name: "empty map",
			got:  map[string]string{},
			want: "",
		},
		{
			name: "single key-value",
			got:  map[string]string{"key": "value"},
			want: `{"key":"value"}`,
		},
		{
			name: "multiple key-values",
			got:  map[string]string{"key1": "value1", "key2": "value2"},
			want: `{"key1":"value1","key2":"value2"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StringMap{}

			result := s.ToString(tt.got)

			if len(tt.got) > 1 {
				var expected, actual map[string]string
				_ = json.Unmarshal([]byte(tt.want), &expected)
				_ = json.Unmarshal([]byte(result), &actual)
				assert.EqualValues(t, expected, actual)

				return
			}

			assert.Equal(t, tt.want, result)
		})
	}
}
