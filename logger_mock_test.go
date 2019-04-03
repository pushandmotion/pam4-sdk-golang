package pam4sdk

import (
	"testing"

	"github.com/3dsinteractive/testify/assert"
)

func TestMockLogger(t *testing.T) {
	var interf ILogger
	interf = NewMockLogger()
	assert.NotNil(t, interf)
}
