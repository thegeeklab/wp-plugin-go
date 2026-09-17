package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserHomeDir(t *testing.T) {
	assert.NotEmpty(t, GetUserHomeDir())
}
