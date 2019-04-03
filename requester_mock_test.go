package pam4sdk

import (
	"testing"

	"github.com/3dsinteractive/testify/assert"
)

func TestMockRequester(t *testing.T) {
	var interf IRequester
	interf = NewMockRequester()
	assert.NotNil(t, interf)

	var cfg IRequesterConfig
	cfg = NewMockRequesterConfig()
	assert.NotNil(t, cfg)
}
