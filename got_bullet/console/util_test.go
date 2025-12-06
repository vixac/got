package console

import (
	"testing"

	"gotest.tools/assert"
)

func TestFitString(t *testing.T) {

	assert.Equal(t, FitString("abcdef", 8, ":"), "abcdef::")
	assert.Equal(t, FitString("a", 2, ":"), "a:")

}
