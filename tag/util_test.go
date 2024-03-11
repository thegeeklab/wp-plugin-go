package tag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		t.Run(tt.name, func(t *testing.T) {
			got := IsTaggable(tt.args.ref, tt.args.defaultBranch)
			assert.Equal(t, got, tt.want, "%q. IsTaggable() = %v, want %v", tt.name, got, tt.want)
		})
	}
}
