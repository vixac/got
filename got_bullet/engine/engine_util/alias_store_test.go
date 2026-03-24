package engine_util

import (
	"testing"

	"gotest.tools/assert"
)

func TestAliasValidity(t *testing.T) {
	assert.Equal(t, IsValidAlias("bob"), true)
	assert.Equal(t, IsValidAlias("1bob"), false)
	assert.Equal(t, IsValidAlias("alice bob"), false)
}
