package tag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_stripTagPrefix(t *testing.T) {
	tests := []struct {
		name   string
		before string
		after  string
	}{
		{name: "strip ref", before: "refs/tags/1.0.0", after: "1.0.0"},
		{name: "strip ref and version prefix", before: "refs/tags/v1.0.0", after: "1.0.0"},
		{name: "strip version prefix", before: "v1.0.0", after: "1.0.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripTagPrefix(tt.before)
			assert.Equal(t, got, tt.after)
		})
	}
}

func TestSemverTagsStrict(t *testing.T) {
	tests := []struct {
		name   string
		before string
		after  []string
	}{
		{name: "empty", before: "", after: []string{"latest"}},
		{name: "main branch", before: "refs/heads/main", after: []string{"latest"}},
		{name: "zero version", before: "refs/tags/0.9.0", after: []string{"0.9", "0.9.0"}},
		{name: "regular version", before: "refs/tags/1.0.0", after: []string{"1", "1.0", "1.0.0"}},
		{name: "regular version with version prefix", before: "refs/tags/v1.0.0", after: []string{"1", "1.0", "1.0.0"}},
		{name: "regular version with meta", before: "refs/tags/v1.0.0+1", after: []string{"1", "1.0", "1.0.0"}},
		{name: "regular pre-release version count", before: "refs/tags/v1.0.0-alpha.1", after: []string{"1.0.0-alpha.1"}},
		{name: "regular pre-release version", before: "refs/tags/v1.0.0-alpha", after: []string{"1.0.0-alpha"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tags, err := SemverTags(tt.before, true)
			assert.NoError(t, err)

			got, want := tags, tt.after
			assert.Equal(t, got, want)
		})
	}
}

func TestSemverTags(t *testing.T) {
	tests := []struct {
		name    string
		before  string
		after   []string
		strict  bool
		wantErr error
	}{
		{
			name:    "empty",
			before:  "",
			after:   []string{"latest"},
			strict:  false,
			wantErr: nil,
		},
		{
			name:    "main branch",
			before:  "refs/heads/main",
			after:   []string{"latest"},
			strict:  false,
			wantErr: nil,
		},
		{
			name:    "zero version",
			before:  "refs/tags/0.9.0",
			after:   []string{"0.9", "0.9.0"},
			strict:  false,
			wantErr: nil,
		},
		{
			name:    "regular version",
			before:  "refs/tags/1.0.0",
			after:   []string{"1", "1.0", "1.0.0"},
			strict:  false,
			wantErr: nil,
		},
		{
			name:    "regular version with prefix",
			before:  "refs/tags/v1.0.0",
			after:   []string{"1", "1.0", "1.0.0"},
			strict:  false,
			wantErr: nil,
		},
		{
			name:    "regular version with meta",
			before:  "refs/tags/v1.0.0+1",
			after:   []string{"1", "1.0", "1.0.0"},
			strict:  false,
			wantErr: nil,
		},
		{
			name:    "regular prerelease version count",
			before:  "refs/tags/v1.0.0-alpha.1",
			after:   []string{"1.0.0-alpha.1"},
			strict:  false,
			wantErr: nil,
		},
		{
			name:    "regular prerelease version",
			before:  "refs/tags/v1.0.0-alpha",
			after:   []string{"1.0.0-alpha"},
			strict:  false,
			wantErr: nil,
		},
		{
			name:    "prerelease version",
			before:  "refs/tags/v1.0-alpha",
			after:   []string{"1.0.0-alpha"},
			strict:  false,
			wantErr: nil,
		},
		{
			name:    "invalid semver",
			before:  "refs/tags/x1.0.0",
			strict:  true,
			wantErr: assert.AnError,
		},
		{
			name:    "date tag",
			before:  "refs/tags/20190203",
			strict:  true,
			wantErr: assert.AnError,
		},
		{
			name:    "regular version",
			before:  "refs/tags/22.04.0",
			strict:  true,
			wantErr: assert.AnError,
		},
		{
			name:    "regular version shorthand",
			before:  "refs/tags/22.4",
			strict:  true,
			wantErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		tags, err := SemverTags(tt.before, tt.strict)
		if tt.wantErr != nil {
			assert.Error(t, err)

			continue
		}

		assert.NoError(t, err)
		assert.Equal(t, tags, tt.after)
	}
}

func TestSemverTagSuffix(t *testing.T) {
	tests := []struct {
		name   string
		before string
		suffix string
		after  []string
	}{
		// without suffix
		{
			name:  "empty ref",
			after: []string{"latest"},
		},
		{
			name:   "ref without suffix",
			before: "refs/tags/v1.0.0",
			after: []string{
				"1",
				"1.0",
				"1.0.0",
			},
		},
		// with suffix
		{
			name:   "empty ref with suffix",
			suffix: "linux-amd64",
			after:  []string{"linux-amd64"},
		},
		{
			name:   "ref with suffix",
			before: "refs/tags/v1.0.0",
			suffix: "linux-amd64",
			after: []string{
				"1-linux-amd64",
				"1.0-linux-amd64",
				"1.0.0-linux-amd64",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tag, err := SemverTagSuffix(tt.before, tt.suffix, true)
			assert.NoError(t, err)

			got, want := tag, tt.after
			assert.Equal(t, got, want)
		})
	}
}

func Test_stripHeadPrefix(t *testing.T) {
	type args struct {
		ref string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "main branch",
			args: args{
				ref: "refs/heads/main",
			},
			want: "main",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripHeadPrefix(tt.args.ref)
			assert.Equal(t, got, tt.want)
		})
	}
}
