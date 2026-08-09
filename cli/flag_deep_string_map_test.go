package cli

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeepStringMapSet(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    map[string]map[string]string
		wantErr bool
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
		{
			name:  "nested with boolean value",
			input: `{"group1":{"enabled":true}}`,
			want:  map[string]map[string]string{"group1": {"enabled": "true"}},
		},
		{
			name:  "nested with numeric values",
			input: `{"group1":{"count":42,"ratio":3.14}}`,
			want:  map[string]map[string]string{"group1": {"count": "42", "ratio": "3.14"}},
		},
		{
			name:  "nested with null value",
			input: `{"group1":{"empty":null}}`,
			want:  map[string]map[string]string{"group1": {"empty": ""}},
		},
		{
			name:  "nested with mixed value types",
			input: `{"group1":{"str":"hello","flag":false,"num":7,"empty":null}}`,
			want:  map[string]map[string]string{"group1": {"str": "hello", "flag": "false", "num": "7", "empty": ""}},
		},
		{
			name:  "flat map with boolean value",
			input: `{"enabled":true}`,
			want:  map[string]map[string]string{"*": {"enabled": "true"}},
		},
		{
			name:  "flat map with numeric values",
			input: `{"count":42,"ratio":3.14}`,
			want:  map[string]map[string]string{"*": {"count": "42", "ratio": "3.14"}},
		},
		{
			name:  "nested with null-only group",
			input: `{"group1":null}`,
			want:  map[string]map[string]string{"group1": nil},
		},
		{
			name:  "flat map with mixed value types",
			input: `{"str":"hello","flag":false,"num":7,"empty":null}`,
			want:  map[string]map[string]string{"*": {"str": "hello", "flag": "false", "num": "7", "empty": ""}},
		},
		{
			name:    "not parseable input returns error",
			input:   `not-json`,
			want:    map[string]map[string]string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got map[string]map[string]string

			d := &DeepStringMap{
				destination: &got,
			}

			err := d.Set(tt.input)

			if tt.wantErr {
				assert.Error(t, err)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDeepStringMapSetRoundtrip(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "valid nested JSON",
			input: `{"group1":{"key1":"value1","key2":"value2"}}`,
			want:  `{"group1":{"key1":"value1","key2":"value2"}}`,
		},
		{
			name:  "nested with non-string values",
			input: `{"group1":{"enabled":true,"count":42}}`,
			want:  `{"group1":{"count":"42","enabled":"true"}}`,
		},
		{
			name:  "flat map fallback",
			input: `{"key1":"value1","key2":"value2"}`,
			want:  `{"*":{"key1":"value1","key2":"value2"}}`,
		},
		{
			name:  "flat map with non-string values",
			input: `{"flag":false,"num":7,"empty":null}`,
			want:  `{"*":{"empty":"","flag":"false","num":"7"}}`,
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

			result := d.String()

			var expected, actual map[string]map[string]string

			_ = json.Unmarshal([]byte(tt.want), &expected)
			_ = json.Unmarshal([]byte(result), &actual)
			assert.EqualValues(t, expected, actual)
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
