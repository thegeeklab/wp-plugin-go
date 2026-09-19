package plugin

import (
	"context"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
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

func TestEnvironment_Value(t *testing.T) {
	tests := []struct {
		name string
		env  Environment
		want []string
	}{
		{
			name: "empty environment",
			env:  Environment{},
			want: []string{},
		},
		{
			name: "single key-value pair",
			env: Environment{
				"KEY1": "value1",
			},
			want: []string{"KEY1=value1"},
		},
		{
			name: "multiple key-value pairs",
			env: Environment{
				"KEY1": "value1",
				"KEY2": "value2",
				"KEY3": "value3",
			},
			want: []string{"KEY1=value1", "KEY2=value2", "KEY3=value3"},
		},
		{
			name: "keys with empty values",
			env: Environment{
				"EMPTY1": "",
				"EMPTY2": "",
			},
			want: []string{"EMPTY1=", "EMPTY2="},
		},
		{
			name: "values with special characters",
			env: Environment{
				"SPECIAL": "value=with=equals",
				"SPACES":  "value with spaces",
			},
			want: []string{"SPECIAL=value=with=equals", "SPACES=value with spaces"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.env.Value()
			// Sort both slices to ensure consistent comparison
			sort.Strings(got)
			sort.Strings(tt.want)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEnvironmentFromContext(t *testing.T) {
	tests := []struct {
		name    string
		env     string
		bare    bool
		want    Environment
		wantErr error
	}{
		{
			name: "environment map from context",
			env:  `{"KEY1":"value1","KEY2":"value2"}`,
			want: Environment{
				"KEY1": "value1",
				"KEY2": "value2",
			},
		},
		{
			name:    "missing environment flag",
			bare:    true,
			wantErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.bare {
				got, err := EnvironmentFromContext(&cli.Command{})
				assert.Error(t, err)
				assert.Nil(t, got)

				return
			}

			plugin := New(Options{
				Name:    "dummy",
				Execute: func(_ context.Context) error { return nil },
			})

			t.Setenv("PLUGIN_ENVIRONMENT", tt.env)

			var got Environment

			plugin.App.Action = func(_ context.Context, cmd *cli.Command) error {
				var err error

				got, err = EnvironmentFromContext(cmd)

				return err
			}

			assert.NoError(t, plugin.App.Run(t.Context(), []string{"dummy"}))
			assert.Equal(t, tt.want, got)
		})
	}
}
