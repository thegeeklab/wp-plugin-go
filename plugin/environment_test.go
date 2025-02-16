package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironment_Lookup(t *testing.T) {
	tests := []struct {
		name   string
		env    Environment
		key    string
		osEnv  map[string]string
		want   string
		wantOk bool
	}{
		{
			name: "value exists in Environment",
			env: Environment{
				"EXISTING_KEY": "env_value",
			},
			key:    "EXISTING_KEY",
			want:   "env_value",
			wantOk: true,
		},
		{
			name: "value exists in OS environment",
			env:  Environment{},
			key:  "OS_ENV_KEY",
			osEnv: map[string]string{
				"OS_ENV_KEY": "os_value",
			},
			want:   "os_value",
			wantOk: true,
		},
		{
			name:   "value does not exist in either environment",
			env:    Environment{},
			key:    "NONEXISTENT_KEY",
			want:   "",
			wantOk: false,
		},
		{
			name:   "empty key lookup",
			env:    Environment{},
			key:    "",
			want:   "",
			wantOk: false,
		},
		{
			name: "Environment value takes precedence over OS environment",
			env: Environment{
				"PRECEDENCE_KEY": "env_value",
			},
			key: "PRECEDENCE_KEY",
			osEnv: map[string]string{
				"PRECEDENCE_KEY": "os_value",
			},
			want:   "env_value",
			wantOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.osEnv {
				t.Setenv(k, v)
			}

			got, ok := tt.env.Lookup(tt.key)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantOk, ok)
		})
	}
}
