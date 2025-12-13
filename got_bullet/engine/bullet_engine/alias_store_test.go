package bullet_engine

import (
	"testing"

	"gotest.tools/assert"
	"vixac.com/got/engine"
)

func TestAliasValidity(t *testing.T) {
	assert.Equal(t, engine.IsValidAlias("bob"), true)
	assert.Equal(t, engine.IsValidAlias("1bob"), false)
	assert.Equal(t, engine.IsValidAlias("alice bob"), false)
}
