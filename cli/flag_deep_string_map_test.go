package cli

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeepStringMapSet(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  map[string]map[string]string
	}{
		{
			name:  "empty string",
			input: "",
			want:  map[string]map[string]string{},
		},
		{
			name:  "valid JSON",
			input: `{"group1":{"key1":"value1","key2":"value2"},"group2":{"key3":"value3"}}`,
			want:  map[string]map[string]string{"group1": {"key1": "value1", "key2": "value2"}, "group2": {"key3": "value3"}},
		},
		{
			name:  "single group",
			input: `{"group1":{"key1":"value1"}}`,
			want:  map[string]map[string]string{"group1": {"key1": "value1"}},
		},
		{
			name:  "single-level map",
			input: `{"key1":"value1","key2":"value2"}`,
			want:  map[string]map[string]string{"*": {"key1": "value1", "key2": "value2"}},
		},
		{
			name:  "empty JSON object",
			input: "{}",
			want:  map[string]map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got map[string]map[string]string
			d := &DeepStringMap{
				destination: &got,
			}

			err := d.Set(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDeepStringMapString(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]map[string]string
		want  string
	}{
		{
			name:  "empty map",
			input: map[string]map[string]string{},
			want:  "",
		},
		{
			name:  "nil map",
			input: nil,
			want:  "",
		},
		{
			name:  "single group",
			input: map[string]map[string]string{"group1": {"key1": "value1"}},
			want:  `{"group1":{"key1":"value1"}}`,
		},
		{
			name:  "multiple groups",
			input: map[string]map[string]string{"group1": {"key1": "value1"}, "group2": {"key2": "value2"}},
			want:  `{"group1":{"key1":"value1"},"group2":{"key2":"value2"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DeepStringMap{
				destination: &tt.input,
			}

			got := d.String()

			if len(tt.input) > 1 {
				var expected, actual map[string]map[string]string
				_ = json.Unmarshal([]byte(tt.want), &expected)
				_ = json.Unmarshal([]byte(got), &actual)
				assert.EqualValues(t, expected, actual)

				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDeepStringMapGet(t *testing.T) {
	tests := []struct {
		name string
		want map[string]map[string]string
	}{
		{
			name: "empty map",
			want: map[string]map[string]string{},
		},
		{
			name: "single group",
			want: map[string]map[string]string{"group1": {"key1": "value1"}},
		},
		{
			name: "multiple groups",
			want: map[string]map[string]string{"group1": {"key1": "value1"}, "group2": {"key2": "value2"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DeepStringMap{
				destination: &tt.want,
			}

			got := d.Get()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDeepStringMapCreate(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]map[string]string
		want  map[string]map[string]string
	}{
		{
			name:  "empty map",
			input: nil,
			want:  map[string]map[string]string{},
		},
		{
			name:  "single group",
			input: map[string]map[string]string{"group1": {"key1": "value1"}},
			want:  map[string]map[string]string{"group1": {"key1": "value1"}},
		},
		{
			name:  "multiple groups",
			input: map[string]map[string]string{"group1": {"key1": "value1"}, "group2": {"key2": "value2"}},
			want:  map[string]map[string]string{"group1": {"key1": "value1"}, "group2": {"key2": "value2"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dest map[string]map[string]string

			d := DeepStringMap{}
			config := DeepStringMapConfig{}
			got := d.Create(tt.input, &dest, config)

			assert.Equal(t, tt.want, dest)
			assert.Equal(t, &dest, got.(*DeepStringMap).destination)
		})
	}
}

func TestDeepStringMapToString(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]map[string]string
		want  string
	}{
		{
			name:  "empty map",
			input: map[string]map[string]string{},
			want:  "",
		},
		{
			name:  "single group",
			input: map[string]map[string]string{"group1": {"key1": "value1"}},
			want:  `{"group1":{"key1":"value1"}}`,
		},
		{
			name:  "multiple groups",
			input: map[string]map[string]string{"group1": {"key1": "value1"}, "group2": {"key2": "value2"}},
			want:  `{"group1":{"key1":"value1"},"group2":{"key2":"value2"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DeepStringMap{}

			got := d.ToString(tt.input)

			if len(tt.input) > 1 {
				var expected, actual map[string]map[string]string
				_ = json.Unmarshal([]byte(tt.want), &expected)
				_ = json.Unmarshal([]byte(got), &actual)
				assert.EqualValues(t, expected, actual)

				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
