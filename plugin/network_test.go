package plugin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

func TestNetworkFromContext(t *testing.T) {
	tests := []struct {
		name           string
		envs           map[string]string
		wantSkipVerify bool
	}{
		{
			name: "insecure skip verify disabled",
			envs: map[string]string{
				"PLUGIN_INSECURE_SKIP_VERIFY": "false",
			},
			wantSkipVerify: false,
		},
		{
			name: "insecure skip verify enabled",
			envs: map[string]string{
				"PLUGIN_INSECURE_SKIP_VERIFY": "true",
			},
			wantSkipVerify: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plugin := New(Options{
				Name:    "dummy",
				Execute: func(_ context.Context) error { return nil },
			})

			t.Setenv("SOCKS_PROXY", "")
			t.Setenv("SOCKS_PROXY_OFF", "")

			for key, value := range tt.envs {
				t.Setenv(key, value)
			}

			var got Network

			plugin.App.Action = func(_ context.Context, cmd *cli.Command) error {
				got = NetworkFromContext(cmd)

				return nil
			}

			assert.NoError(t, plugin.App.Run(t.Context(), []string{"dummy"}))

			assert.Equal(t, tt.wantSkipVerify, got.InsecureSkipVerify)
			assert.NotNil(t, got.Client)
			assert.NotNil(t, got.Context)
		})
	}
}
