package tag

import (
	"reflect"
	"testing"
)

func Test_stripTagPrefix(t *testing.T) {
	tests := []struct {
		Before string
		After  string
	}{
		{"refs/tags/1.0.0", "1.0.0"},
		{"refs/tags/v1.0.0", "1.0.0"},
		{"v1.0.0", "1.0.0"},
	}

	for _, test := range tests {
		got, want := stripTagPrefix(test.Before), test.After
		if got != want {
			t.Errorf("Got tag %s, want %s", got, want)
		}
	}
}

func TestSemverTagsStrict(t *testing.T) {
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
	}

	for _, test := range tests {
		tags, err := SemverTags(test.Before, true)
		if err != nil {
			t.Error(err)

			continue
		}

		got, want := tags, test.After
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Got tag %v, want %v", got, want)
		}
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

	for _, test := range tests {
		tags, err := SemverTags(test.Before, false)
		if err != nil {
			t.Error(err)

			continue
		}

		got, want := tags, test.After
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Got tag %v, want %v", got, want)
		}
	}
}

func TestSemverTagsSrtictError(t *testing.T) {
	tests := []string{
		"refs/tags/x1.0.0",
		"refs/tags/20190203",
		"refs/tags/22.04.0",
		"refs/tags/22.04",
	}

	for _, test := range tests {
		_, err := SemverTags(test, true)
		if err == nil {
			t.Errorf("Expect tag error for %s", test)
		}
	}
}

func TestSemverTagSuffix(t *testing.T) {
	tests := []struct {
		Before string
		Suffix string
		After  []string
	}{
		// without suffix
		{
			After: []string{"latest"},
		},
		{
			Before: "refs/tags/v1.0.0",
			After: []string{
				"1",
				"1.0",
				"1.0.0",
			},
		},
		// with suffix
		{
			Suffix: "linux-amd64",
			After:  []string{"linux-amd64"},
		},
		{
			Before: "refs/tags/v1.0.0",
			Suffix: "linux-amd64",
			After: []string{
				"1-linux-amd64",
				"1.0-linux-amd64",
				"1.0.0-linux-amd64",
			},
		},
		{
			Suffix: "nanoserver",
			After:  []string{"nanoserver"},
		},
		{
			Before: "refs/tags/v1.9.2",
			Suffix: "nanoserver",
			After: []string{
				"1-nanoserver",
				"1.9-nanoserver",
				"1.9.2-nanoserver",
			},
		},
	}

	for _, test := range tests {
		tag, err := SemverTagSuffix(test.Before, test.Suffix, true)
		if err != nil {
			t.Error(err)

			continue
		}

		got, want := tag, test.After
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Got tag %v, want %v", got, want)
		}
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
		if got := stripHeadPrefix(tt.args.ref); got != tt.want {
			t.Errorf("stripHeadPrefix() = %v, want %v", got, tt.want)
		}
	}
}

func TestIsTaggable(t *testing.T) {
	type args struct {
		ref           string
		defaultBranch string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "latest tag for default branch",
			args: args{
				ref:           "refs/heads/main",
				defaultBranch: "main",
			},
			want: true,
		},
		{
			name: "build from tags",
			args: args{
				ref:           "refs/tags/v1.0.0",
				defaultBranch: "main",
			},
			want: true,
		},
		{
			name: "skip build for not default branch",
			args: args{
				ref:           "refs/heads/develop",
				defaultBranch: "main",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		if got := IsTaggable(tt.args.ref, tt.args.defaultBranch); got != tt.want {
			t.Errorf("%q. IsTaggable() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
