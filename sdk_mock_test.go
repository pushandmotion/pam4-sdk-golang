package pam4sdk

import (
	"testing"

	"github.com/3dsinteractive/testify/assert"
)

func TestMockSdk(t *testing.T) {
	var interf ISdk
	interf = NewMockSdk()
	assert.NotNil(t, interf)
}
