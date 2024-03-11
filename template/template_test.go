package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddPrefix(t *testing.T) {
	type args struct {
		prefix string
		input  string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty input",
			args: args{
				prefix: "pre",
				input:  "",
			},
			want: "",
		},
		{
			name: "input already has prefix",
			args: args{
				prefix: "pre",
				input:  "pre-existing",
			},
			want: "pre-existing",
		},
		{
			name: "add prefix",
			args: args{
				prefix: "pre",
				input:  "-existing",
			},
			want: "pre-existing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddPrefix(tt.args.prefix, tt.args.input)

			assert.Equal(t, tt.want, got)
		})
	}
}
