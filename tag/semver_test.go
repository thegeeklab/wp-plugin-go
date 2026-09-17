package tag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_stripTagPrefix(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "strip ref",
			input: "refs/tags/1.0.0",
			want:  "1.0.0",
		},
		{
			name:  "strip ref and version prefix",
			input: "refs/tags/v1.0.0",
			want:  "1.0.0",
		},
		{
			name:  "strip version prefix",
			input: "v1.0.0",
			want:  "1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripTagPrefix(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSemverTagsStrict(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "empty",
			input: "",
			want:  []string{"latest"},
		},
		{
			name:  "main branch",
			input: "refs/heads/main",
			want:  []string{"latest"},
		},
		{
			name:  "zero version",
			input: "refs/tags/0.9.0",
			want:  []string{"0.9", "0.9.0"},
		},
		{
			name:  "regular version",
			input: "refs/tags/1.0.0",
			want:  []string{"1", "1.0", "1.0.0"},
		},
		{
			name:  "regular version with version prefix",
			input: "refs/tags/v1.0.0",
			want:  []string{"1", "1.0", "1.0.0"},
		},
		{
			name:  "regular version with meta",
			input: "refs/tags/v1.0.0+1",
			want:  []string{"1", "1.0", "1.0.0"},
		},
		{
			name:  "regular pre-release version count",
			input: "refs/tags/v1.0.0-alpha.1",
			want:  []string{"1.0.0-alpha.1"},
		},
		{
			name:  "regular pre-release version",
			input: "refs/tags/v1.0.0-alpha",
			want:  []string{"1.0.0-alpha"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tags, err := SemverTags(tt.input, true)
			assert.NoError(t, err)

			assert.Equal(t, tt.want, tags)
		})
	}
}

func TestSemverTags(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		strict  bool
		wantErr error
	}{
		{
			name:    "empty",
			input:   "",
			want:    []string{"latest"},
			strict:  false,
			wantErr: nil,
		},
		{
			name:    "main branch",
			input:   "refs/heads/main",
			want:    []string{"latest"},
			strict:  false,
			wantErr: nil,
		},
		{
			name:    "zero version",
			input:   "refs/tags/0.9.0",
			want:    []string{"0.9", "0.9.0"},
			strict:  false,
			wantErr: nil,
		},
		{
			name:    "regular version",
			input:   "refs/tags/1.0.0",
			want:    []string{"1", "1.0", "1.0.0"},
			strict:  false,
			wantErr: nil,
		},
		{
			name:    "regular version with prefix",
			input:   "refs/tags/v1.0.0",
			want:    []string{"1", "1.0", "1.0.0"},
			strict:  false,
			wantErr: nil,
		},
		{
			name:    "regular version with meta",
			input:   "refs/tags/v1.0.0+1",
			want:    []string{"1", "1.0", "1.0.0"},
			strict:  false,
			wantErr: nil,
		},
		{
			name:    "regular prerelease version count",
			input:   "refs/tags/v1.0.0-alpha.1",
			want:    []string{"1.0.0-alpha.1"},
			strict:  false,
			wantErr: nil,
		},
		{
			name:    "regular prerelease version",
			input:   "refs/tags/v1.0.0-alpha",
			want:    []string{"1.0.0-alpha"},
			strict:  false,
			wantErr: nil,
		},
		{
			name:    "prerelease version",
			input:   "refs/tags/v1.0-alpha",
			want:    []string{"1.0.0-alpha"},
			strict:  false,
			wantErr: nil,
		},
		{
			name:    "invalid semver",
			input:   "refs/tags/x1.0.0",
			strict:  true,
			wantErr: assert.AnError,
		},
		{
			name:    "date tag",
			input:   "refs/tags/20190203",
			strict:  true,
			wantErr: assert.AnError,
		},
		{
			name:    "version with leading zero",
			input:   "refs/tags/22.04.0",
			strict:  true,
			wantErr: assert.AnError,
		},
		{
			name:    "version shorthand",
			input:   "refs/tags/22.4",
			strict:  true,
			wantErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tags, err := SemverTags(tt.input, tt.strict)
			if tt.wantErr != nil {
				assert.Error(t, err)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, tags)
		})
	}
}

func TestSemverTagSuffix(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		suffix string
		want   []string
	}{
		// without suffix
		{
			name: "empty ref",
			want: []string{"latest"},
		},
		{
			name:  "ref without suffix",
			input: "refs/tags/v1.0.0",
			want: []string{
				"1",
				"1.0",
				"1.0.0",
			},
		},
		// with suffix
		{
			name:   "empty ref with suffix",
			suffix: "linux-amd64",
			want:   []string{"linux-amd64"},
		},
		{
			name:   "ref with suffix",
			input:  "refs/tags/v1.0.0",
			suffix: "linux-amd64",
			want: []string{
				"1-linux-amd64",
				"1.0-linux-amd64",
				"1.0.0-linux-amd64",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tags, err := SemverTagSuffix(tt.input, tt.suffix, true)
			assert.NoError(t, err)

			assert.Equal(t, tt.want, tags)
		})
	}
}

func Test_stripHeadPrefix(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "main branch",
			input: "refs/heads/main",
			want:  "main",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripHeadPrefix(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
