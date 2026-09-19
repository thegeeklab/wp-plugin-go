package plugin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
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
				Name:    "dummy",
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
