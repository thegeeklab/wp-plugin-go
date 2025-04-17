package cli

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapSet(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    map[string]string
		wantErr error
	}{
		{
			name:    "empty string",
			input:   "",
			want:    map[string]string{},
			wantErr: nil,
		},
		{
			name:    "valid JSON",
			input:   `{"key1":"value1","key2":"value2"}`,
			want:    map[string]string{"key1": "value1", "key2": "value2"},
			wantErr: nil,
		},
		{
			name:    "single key-value",
			input:   `{"key":"value"}`,
			want:    map[string]string{"key": "value"},
			wantErr: nil,
		},
		{
			name:    "empty JSON object",
			input:   "{}",
			want:    map[string]string{},
			wantErr: nil,
		},
		{
			name:    "invalid JSON",
			input:   "not-json",
			want:    map[string]string{},
			wantErr: &json.SyntaxError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got map[string]string
			m := &Map{
				destination: &got,
			}

			err := m.Set(tt.input)

			if tt.wantErr != nil {
				assert.ErrorAs(t, err, &tt.wantErr)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMapString(t *testing.T) {
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
			m := &Map{
				destination: &tt.input,
			}

			got := m.String()

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

func TestMapGet(t *testing.T) {
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
			m := &Map{
				destination: &tt.want,
			}

			got := m.Get()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMapCreate(t *testing.T) {
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

			m := Map{}
			config := MapConfig{}

			got := m.Create(tt.input, &dest, config)
			assert.Equal(t, tt.want, dest)
			assert.Equal(t, &dest, got.(*Map).destination)
		})
	}
}

func TestMapToString(t *testing.T) {
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
			m := Map{}

			got := m.ToString(tt.input)

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
