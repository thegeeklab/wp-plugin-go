package cli

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeepStringMapSet(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want map[string]map[string]string
	}{
		{
			name: "empty string",
			got:  "",
			want: map[string]map[string]string{},
		},
		{
			name: "valid JSON",
			got:  `{"group1":{"key1":"value1","key2":"value2"},"group2":{"key3":"value3"}}`,
			want: map[string]map[string]string{"group1": {"key1": "value1", "key2": "value2"}, "group2": {"key3": "value3"}},
		},
		{
			name: "single group",
			got:  `{"group1":{"key1":"value1"}}`,
			want: map[string]map[string]string{"group1": {"key1": "value1"}},
		},
		{
			name: "single-level map",
			got:  `{"key1":"value1","key2":"value2"}`,
			want: map[string]map[string]string{"*": {"key1": "value1", "key2": "value2"}},
		},
		{
			name: "empty JSON object",
			got:  "{}",
			want: map[string]map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dest map[string]map[string]string
			d := &DeepStringMap{
				destination: &dest,
			}

			err := d.Set(tt.got)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, dest)
		})
	}
}

func TestDeepStringMapString(t *testing.T) {
	tests := []struct {
		name string
		got  map[string]map[string]string
		want string
	}{
		{
			name: "empty map",
			got:  map[string]map[string]string{},
			want: "",
		},
		{
			name: "nil map",
			got:  nil,
			want: "",
		},
		{
			name: "single group",
			got:  map[string]map[string]string{"group1": {"key1": "value1"}},
			want: `{"group1":{"key1":"value1"}}`,
		},
		{
			name: "multiple groups",
			got:  map[string]map[string]string{"group1": {"key1": "value1"}, "group2": {"key2": "value2"}},
			want: `{"group1":{"key1":"value1"},"group2":{"key2":"value2"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DeepStringMap{
				destination: &tt.got,
			}

			result := d.String()

			if len(tt.got) > 1 {
				var expected, actual map[string]map[string]string
				_ = json.Unmarshal([]byte(tt.want), &expected)
				_ = json.Unmarshal([]byte(result), &actual)
				assert.EqualValues(t, expected, actual)

				return
			}

			assert.Equal(t, tt.want, result)
		})
	}
}

func TestDeepStringMapGet(t *testing.T) {
	tests := []struct {
		name string
		got  map[string]map[string]string
	}{
		{
			name: "empty map",
			got:  map[string]map[string]string{},
		},
		{
			name: "single group",
			got:  map[string]map[string]string{"group1": {"key1": "value1"}},
		},
		{
			name: "multiple groups",
			got:  map[string]map[string]string{"group1": {"key1": "value1"}, "group2": {"key2": "value2"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DeepStringMap{
				destination: &tt.got,
			}

			result := d.Get()
			assert.Equal(t, tt.got, result)
		})
	}
}

func TestDeepStringMapCreate(t *testing.T) {
	tests := []struct {
		name string
		got  map[string]map[string]string
	}{
		{
			name: "empty map",
			got:  map[string]map[string]string{},
		},
		{
			name: "single group",
			got:  map[string]map[string]string{"group1": {"key1": "value1"}},
		},
		{
			name: "multiple groups",
			got:  map[string]map[string]string{"group1": {"key1": "value1"}, "group2": {"key2": "value2"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dest map[string]map[string]string

			d := DeepStringMap{}
			config := DeepStringMapConfig{}
			value := d.Create(tt.got, &dest, config)

			assert.Equal(t, tt.got, dest)
			assert.Equal(t, &dest, value.(*DeepStringMap).destination)
		})
	}
}

func TestDeepStringMapToString(t *testing.T) {
	tests := []struct {
		name string
		got  map[string]map[string]string
		want string
	}{
		{
			name: "empty map",
			got:  map[string]map[string]string{},
			want: "",
		},
		{
			name: "single group",
			got:  map[string]map[string]string{"group1": {"key1": "value1"}},
			want: `{"group1":{"key1":"value1"}}`,
		},
		{
			name: "multiple groups",
			got:  map[string]map[string]string{"group1": {"key1": "value1"}, "group2": {"key2": "value2"}},
			want: `{"group1":{"key1":"value1"},"group2":{"key2":"value2"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DeepStringMap{}

			result := d.ToString(tt.got)

			if len(tt.got) > 1 {
				var expected, actual map[string]map[string]string
				_ = json.Unmarshal([]byte(tt.want), &expected)
				_ = json.Unmarshal([]byte(result), &actual)
				assert.EqualValues(t, expected, actual)

				return
			}

			assert.Equal(t, tt.want, result)
		})
	}
}
