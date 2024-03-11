package tag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_stripTagPrefix(t *testing.T) {
	tests := []struct {
		before string
		after  string
	}{
		{before: "refs/tags/1.0.0", after: "1.0.0"},
		{before: "refs/tags/v1.0.0", after: "1.0.0"},
		{before: "v1.0.0", after: "1.0.0"},
	}

	for _, tt := range tests {
		got, want := stripTagPrefix(tt.before), tt.after
		assert.Equal(t, got, want)
	}
}

func TestSemverTagsStrict(t *testing.T) {
	tests := []struct {
		before string
		after  []string
	}{
		{before: "", after: []string{"latest"}},
		{before: "refs/heads/main", after: []string{"latest"}},
		{before: "refs/tags/0.9.0", after: []string{"0.9", "0.9.0"}},
		{before: "refs/tags/1.0.0", after: []string{"1", "1.0", "1.0.0"}},
		{before: "refs/tags/v1.0.0", after: []string{"1", "1.0", "1.0.0"}},
		{before: "refs/tags/v1.0.0+1", after: []string{"1", "1.0", "1.0.0"}},
		{before: "refs/tags/v1.0.0-alpha.1", after: []string{"1.0.0-alpha.1"}},
		{before: "refs/tags/v1.0.0-alpha", after: []string{"1.0.0-alpha"}},
	}

	for _, tt := range tests {
		tags, err := SemverTags(tt.before, true)
		assert.NoError(t, err)

		got, want := tags, tt.after
		assert.Equal(t, got, want)
	}
}

func TestSemverTags(t *testing.T) {
	tests := []struct {
		Before string
		After  []string
	}{
		{"", []string{"latest"}},
		{"refs/heads/main", []string{"latest"}},
		{"refs/tags/0.9.0", []string{"0.9", "0.9.0"}},
		{"refs/tags/1.0.0", []string{"1", "1.0", "1.0.0"}},
		{"refs/tags/v1.0.0", []string{"1", "1.0", "1.0.0"}},
		{"refs/tags/v1.0.0+1", []string{"1", "1.0", "1.0.0"}},
		{"refs/tags/v1.0.0-alpha.1", []string{"1.0.0-alpha.1"}},
		{"refs/tags/v1.0.0-alpha", []string{"1.0.0-alpha"}},
		{"refs/tags/v1.0-alpha", []string{"1.0.0-alpha"}},
		{"refs/tags/22.04.0", []string{"22", "22.4", "22.4.0"}},
		{"refs/tags/22.04", []string{"22", "22.4", "22.4.0"}},
	}
	for _, tt := range tests {
		tags, err := SemverTags(tt.Before, false)
		assert.NoError(t, err)

		got, want := tags, tt.After
		assert.Equal(t, got, want)
	}
}

func TestSemverTagsSrtictError(t *testing.T) {
	tests := []string{
		"refs/tags/x1.0.0",
		"refs/tags/20190203",
		"refs/tags/22.04.0",
		"refs/tags/22.04",
	}

	for _, tt := range tests {
		_, err := SemverTags(tt, true)
		assert.Error(t, err, "Expect tag error for %s", tt)
	}
}

func TestSemverTagSuffix(t *testing.T) {
	tests := []struct {
		before string
		suffix string
		after  []string
	}{
		// without suffix
		{
			after: []string{"latest"},
		},
		{
			before: "refs/tags/v1.0.0",
			after: []string{
				"1",
				"1.0",
				"1.0.0",
			},
		},
		// with suffix
		{
			suffix: "linux-amd64",
			after:  []string{"linux-amd64"},
		},
		{
			before: "refs/tags/v1.0.0",
			suffix: "linux-amd64",
			after: []string{
				"1-linux-amd64",
				"1.0-linux-amd64",
				"1.0.0-linux-amd64",
			},
		},
		{
			suffix: "nanoserver",
			after:  []string{"nanoserver"},
		},
		{
			before: "refs/tags/v1.9.2",
			suffix: "nanoserver",
			after: []string{
				"1-nanoserver",
				"1.9-nanoserver",
				"1.9.2-nanoserver",
			},
		},
	}

	for _, tt := range tests {
		tag, err := SemverTagSuffix(tt.before, tt.suffix, true)
		assert.NoError(t, err)

		got, want := tag, tt.after
		assert.Equal(t, got, want)
	}
}

func Test_stripHeadPrefix(t *testing.T) {
	type args struct {
		ref string
	}

	tests := []struct {
		args args
		want string
	}{
		{
			args: args{
				ref: "refs/heads/main",
			},
			want: "main",
		},
	}

	for _, tt := range tests {
		got := stripHeadPrefix(tt.args.ref)
		assert.Equal(t, got, tt.want)
	}
}
