package plugin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

func TestPluginAction(t *testing.T) {
	allFlags := func() []cli.Flag {
		flags := LoggingFlags(FlagsPluginCategory)
		flags = append(flags, NetworkFlags(FlagsPluginCategory)...)
		flags = append(flags, EnvironmentFlags(FlagsPluginCategory)...)

		return flags
	}

	tests := []struct {
		name               string
		flags              []cli.Flag
		execute            ExecuteFunc
		wantErr            error
		wantNetworkClient  bool
		wantEnvironmentNil bool
	}{
		{
			name:               "execute runs successfully with all flags",
			flags:              allFlags(),
			execute:            func(_ context.Context) error { return nil },
			wantNetworkClient:  true,
			wantEnvironmentNil: false,
		},
		{
			name:               "execute returns error",
			flags:              allFlags(),
			execute:            func(_ context.Context) error { return assert.AnError },
			wantErr:            assert.AnError,
			wantNetworkClient:  true,
			wantEnvironmentNil: false,
		},
		{
			name:               "succeeds without any opt-in flags",
			flags:              nil,
			execute:            func(_ context.Context) error { return nil },
			wantNetworkClient:  false,
			wantEnvironmentNil: true,
		},
		{
			name:               "succeeds with only logging flags",
			flags:              LoggingFlags(FlagsPluginCategory),
			execute:            func(_ context.Context) error { return nil },
			wantNetworkClient:  false,
			wantEnvironmentNil: true,
		},
		{
			name:               "succeeds with only network flags",
			flags:              NetworkFlags(FlagsPluginCategory),
			execute:            func(_ context.Context) error { return nil },
			wantNetworkClient:  true,
			wantEnvironmentNil: true,
		},
		{
			name:               "succeeds with only environment flags",
			flags:              EnvironmentFlags(FlagsPluginCategory),
			execute:            func(_ context.Context) error { return nil },
			wantNetworkClient:  false,
			wantEnvironmentNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plugin := New(Options{
				Name:    "dummy",
				Flags:   tt.flags,
				Execute: tt.execute,
			})

			err := plugin.App.Run(t.Context(), []string{"dummy"})
			if tt.wantErr != nil {
				assert.Error(t, err)

				return
			}

			assert.NoError(t, err)

			network, err := plugin.GetNetwork()
			if tt.wantNetworkClient {
				assert.NoError(t, err)
				assert.NotNil(t, network.Client)
			} else {
				assert.Error(t, err)
			}

			environment, err := plugin.GetEnvironment()
			if tt.wantEnvironmentNil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, environment)
			}
		})
	}
}

func TestPluginFlags(t *testing.T) {
	tests := []struct {
		name        string
		flags       []cli.Flag
		wantPresent []string
		wantAbsent  []string
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
