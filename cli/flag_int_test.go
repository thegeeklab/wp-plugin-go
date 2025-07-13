package cli

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntSet(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr error
	}{
		{
			name:    "empty string",
			input:   "",
			want:    0,
			wantErr: nil,
		},
		{
			name:    "zero",
			input:   "0",
			want:    0,
			wantErr: nil,
		},
		{
			name:    "positive integer",
			input:   "42",
			want:    42,
			wantErr: nil,
		},
		{
			name:    "negative integer",
			input:   "-42",
			want:    -42,
			wantErr: nil,
		},
		{
			name:    "invalid integer",
			input:   "not-an-int",
			want:    0,
			wantErr: &strconv.NumError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got int

			i := &Int{
				destination: &got,
			}

			err := i.Set(tt.input)
			if tt.wantErr != nil {
				assert.ErrorAs(t, err, &tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestIntString(t *testing.T) {
	tests := []struct {
		name  string
		input *int
		want  string
	}{
		{
			name:  "nil pointer",
			input: nil,
			want:  "0",
		},
		{
			name:  "zero",
			input: intPtr(0),
			want:  "0",
		},
		{
			name:  "positive integer",
			input: intPtr(42),
			want:  "42",
		},
		{
			name:  "negative integer",
			input: intPtr(-42),
			want:  "-42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Int{
				destination: tt.input,
			}

			got := i.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIntGet(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "zero",
			want: 0,
		},
		{
			name: "positive integer",
			want: 42,
		},
		{
			name: "negative integer",
			want: -42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Int{
				destination: &tt.want,
			}

			result := i.Get()
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestIntCreate(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "zero",
			want: 0,
		},
		{
			name: "positive integer",
			want: 42,
		},
		{
			name: "negative integer",
			want: -42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dest int

			i := Int{}
			config := IntConfig{}

			got := i.Create(tt.want, &dest, config)
			assert.Equal(t, tt.want, dest)
			assert.Equal(t, &dest, got.(*Int).destination)
		})
	}
}

func TestIntToString(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  string
	}{
		{
			name:  "zero",
			input: 0,
			want:  "0",
		},
		{
			name:  "positive integer",
			input: 42,
			want:  "42",
		},
		{
			name:  "negative integer",
			input: -42,
			want:  "-42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Int{}

			result := i.ToString(tt.input)
			assert.Equal(t, tt.want, result)
		})
	}
}

// Helper function to create an int pointer.
func intPtr(i int) *int {
	return &i
}
