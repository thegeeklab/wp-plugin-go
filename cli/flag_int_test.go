package cli

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntSet(t *testing.T) {
	tests := []struct {
		name    string
		got     string
		want    int
		wantErr error
	}{
		{
			name:    "empty string",
			got:     "",
			want:    0,
			wantErr: nil,
		},
		{
			name:    "zero",
			got:     "0",
			want:    0,
			wantErr: nil,
		},
		{
			name:    "positive integer",
			got:     "42",
			want:    42,
			wantErr: nil,
		},
		{
			name:    "negative integer",
			got:     "-42",
			want:    -42,
			wantErr: nil,
		},
		{
			name:    "invalid integer",
			got:     "not-an-int",
			want:    0,
			wantErr: &strconv.NumError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dest int
			i := &Int{
				destination: &dest,
			}

			err := i.Set(tt.got)

			if tt.wantErr != nil {
				assert.ErrorAs(t, err, &tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, dest)
			}
		})
	}
}

func TestIntString(t *testing.T) {
	tests := []struct {
		name string
		got  *int
		want string
	}{
		{
			name: "nil pointer",
			got:  nil,
			want: "0",
		},
		{
			name: "zero",
			got:  intPtr(0),
			want: "0",
		},
		{
			name: "positive integer",
			got:  intPtr(42),
			want: "42",
		},
		{
			name: "negative integer",
			got:  intPtr(-42),
			want: "-42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Int{
				destination: tt.got,
			}

			result := i.String()
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestIntGet(t *testing.T) {
	tests := []struct {
		name string
		got  int
	}{
		{
			name: "zero",
			got:  0,
		},
		{
			name: "positive integer",
			got:  42,
		},
		{
			name: "negative integer",
			got:  -42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Int{
				destination: &tt.got,
			}

			result := i.Get()
			assert.Equal(t, tt.got, result)
		})
	}
}

func TestIntCreate(t *testing.T) {
	tests := []struct {
		name string
		got  int
	}{
		{
			name: "zero",
			got:  0,
		},
		{
			name: "positive integer",
			got:  42,
		},
		{
			name: "negative integer",
			got:  -42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dest int

			i := Int{}
			config := IntConfig{}

			value := i.Create(tt.got, &dest, config)
			assert.Equal(t, tt.got, dest)
			assert.Equal(t, &dest, value.(*Int).destination)
		})
	}
}

func TestIntToString(t *testing.T) {
	tests := []struct {
		name string
		got  int
		want string
	}{
		{
			name: "zero",
			got:  0,
			want: "0",
		},
		{
			name: "positive integer",
			got:  42,
			want: "42",
		},
		{
			name: "negative integer",
			got:  -42,
			want: "-42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Int{}

			result := i.ToString(tt.got)
			assert.Equal(t, tt.want, result)
		})
	}
}

// Helper function to create an int pointer.
func intPtr(i int) *int {
	return &i
}
