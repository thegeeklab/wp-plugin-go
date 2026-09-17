package tag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTaggable(t *testing.T) {
	tests := []struct {
		name          string
		ref           string
		defaultBranch string
		want          bool
	}{
		{
			name:          "latest tag for default branch",
			ref:           "refs/heads/main",
			defaultBranch: "main",
			want:          true,
		},
		{
			name:          "build from tags",
			ref:           "refs/tags/v1.0.0",
			defaultBranch: "main",
			want:          true,
		},
		{
			name:          "skip build for not default branch",
			ref:           "refs/heads/develop",
			defaultBranch: "main",
			want:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsTaggable(tt.ref, tt.defaultBranch)
			assert.Equal(t, tt.want, got, "%q. IsTaggable() = %v, want %v", tt.name, got, tt.want)
		})
	}
}
