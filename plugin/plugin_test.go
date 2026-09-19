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
			plugin := New(Options{
				Name: "dummy",
				Flags: append(
					[]cli.Flag{},
					append(
						LoggingFlags(FlagsPluginCategory),
						append(
							NetworkFlags(FlagsPluginCategory),
							EnvironmentFlags(FlagsPluginCategory)...,
						)...,
					)...,
				),
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
	t.Run("metadata flags are always registered", func(t *testing.T) {
		plugin := New(Options{
			Name:    "test",
			Execute: func(_ context.Context) error { return nil },
		})

		var flagNames []string
		for _, f := range plugin.App.Flags {
			flagNames = append(flagNames, f.Names()[0])
		}

		// Metadata flags should be present
		assert.Contains(t, flagNames, "repo.slug")
		assert.Contains(t, flagNames, "pipeline.number")
		assert.Contains(t, flagNames, "commit.sha")
		assert.Contains(t, flagNames, "step.number")
		assert.Contains(t, flagNames, "system.name")

		// Plugin flags should NOT be present (opt-in)
		assert.NotContains(t, flagNames, "log-level")
		assert.NotContains(t, flagNames, "transport.insecure-skip-verify")
		assert.NotContains(t, flagNames, "environment")
	})

	t.Run("plugin flags can be opted in", func(t *testing.T) {
		plugin := New(Options{
			Name: "test",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "my-flag",
					Usage:   "my flag",
					Sources: cli.EnvVars("PLUGIN_MY_FLAG"),
				},
			},
			Execute: func(_ context.Context) error { return nil },
		})

		var flagNames []string
		for _, f := range plugin.App.Flags {
			flagNames = append(flagNames, f.Names()[0])
		}

		// User flag should be present
		assert.Contains(t, flagNames, "my-flag")

		// Metadata flags should still be present
		assert.Contains(t, flagNames, "repo.slug")
	})
}
