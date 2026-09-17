package trace

import (
	"context"
	"net/http/httptrace"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTP(t *testing.T) {
	tests := []struct {
		name string
		ctx  func() context.Context
	}{
		{
			name: "background context",
			ctx:  context.Background,
		},
		{
			name: "todo context",
			ctx:  context.TODO,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HTTP(tt.ctx())

			assert.NotNil(t, httptrace.ContextClientTrace(got))
			assert.Nil(t, got.Err())
		})
	}
}
