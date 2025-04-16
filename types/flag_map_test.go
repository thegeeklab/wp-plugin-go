package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapSet(t *testing.T) {
	tests := []struct {
		name    string
		got     string
		want    map[string]string
		wantErr error
	}{
		{
			name:    "empty string",
			got:     "",
			want:    map[string]string{},
			wantErr: nil,
		},
		{
			name:    "valid JSON",
			got:     `{"key1":"value1","key2":"value2"}`,
			want:    map[string]string{"key1": "value1", "key2": "value2"},
			wantErr: nil,
		},
		{
			name:    "single key-value",
			got:     `{"key":"value"}`,
			want:    map[string]string{"key": "value"},
			wantErr: nil,
		},
		{
			name:    "empty JSON object",
			got:     "{}",
			want:    map[string]string{},
			wantErr: nil,
		},
		{
			name:    "invalid JSON",
			got:     "not-json",
			want:    map[string]string{},
			wantErr: &json.SyntaxError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dest map[string]string
			m := &Map{
				destination: &dest,
			}

			err := m.Set(tt.got)

			if tt.wantErr != nil {
				assert.ErrorAs(t, err, &tt.wantErr)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, dest)
		})
	}
}

func TestMapString(t *testing.T) {
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
			m := &Map{
				destination: &tt.got,
			}

			result := m.String()

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

func TestMapGet(t *testing.T) {
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
			m := &Map{
				destination: &tt.got,
			}

			result := m.Get()
			assert.Equal(t, tt.got, result)
		})
	}
}

func TestMapCreate(t *testing.T) {
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

			m := Map{}
			config := MapConfig{}

			value := m.Create(tt.got, &dest, config)
			assert.Equal(t, tt.got, dest)
			assert.Equal(t, &dest, value.(*Map).destination)
		})
	}
}

//nolint:dupl
func TestMapToString(t *testing.T) {
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
			m := Map{}

			result := m.ToString(tt.got)

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
