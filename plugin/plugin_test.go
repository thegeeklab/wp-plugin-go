package plugin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

func TestPluginAction(t *testing.T) {
	tests := []struct {
		name    string
		execute ExecuteFunc
		wantErr error
	}{
		{
			name:    "execute runs successfully",
			execute: func(_ context.Context) error { return nil },
		},
		{
			name:    "execute returns error",
			execute: func(_ context.Context) error { return assert.AnError },
			wantErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := LoggingFlags(FlagsPluginCategory)
			flags = append(flags, NetworkFlags(FlagsPluginCategory)...)
			flags = append(flags, EnvironmentFlags(FlagsPluginCategory)...)

			plugin := New(Options{
				Name:    "dummy",
				Flags:   flags,
				Execute: tt.execute,
			})

			err := plugin.App.Run(t.Context(), []string{"dummy"})
			if tt.wantErr != nil {
				assert.Error(t, err)

				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, plugin.Network.Client)
			assert.NotNil(t, plugin.Environment)
		})
	}
}

func TestPluginFlags(t *testing.T) {
	tests := []struct {
		name         string
		flags        []cli.Flag
		wantPresent  []string
		wantAbsent   []string
	}{
		{
			name:  "metadata flags are always registered",
			flags: nil,
			wantPresent: []string{
				"repo.slug",
				"pipeline.number",
				"commit.sha",
				"step.number",
				"system.name",
			},
			wantAbsent: []string{
				"log-level",
				"transport.insecure-skip-verify",
				"environment",
			},
		},
		{
			name: "plugin flags can be opted in",
			flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "my-flag",
					Usage:   "my flag",
					Sources: cli.EnvVars("PLUGIN_MY_FLAG"),
				},
			},
			wantPresent: []string{"my-flag", "repo.slug"},
			wantAbsent:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plugin := New(Options{
				Name:    "test",
				Flags:   tt.flags,
				Execute: func(_ context.Context) error { return nil },
			})

			var flagNames []string
			for _, f := range plugin.App.Flags {
				flagNames = append(flagNames, f.Names()[0])
			}

			for _, flag := range tt.wantPresent {
				assert.Contains(t, flagNames, flag)
			}

			for _, flag := range tt.wantAbsent {
				assert.NotContains(t, flagNames, flag)
			}
		})
	}
}
