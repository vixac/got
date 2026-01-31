package cmd

import (
	"errors"
	"testing"

	"gotest.tools/assert"
)

func TestRemoveommand_MissingGID(t *testing.T) {

	var p = MockMessenger{}
	var e = MockEngine{}
	mockDeps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}
	cmd := buildRemoveCommand(mockDeps)
	cmd.SetArgs([]string{}) // no gid
	_ = cmd.Execute()
	assert.Equal(t, len(p.errors), 1)

	assert.Equal(t, p.errors[0].Message, "Expected at least one lookup as input")
}

func TestRemoveCommand_Valid(t *testing.T) {

	var p = MockMessenger{}
	var e = MockEngine{}
	mockDeps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}
	cmd := buildRemoveCommand(mockDeps)
	cmd.SetArgs([]string{"abc"}) // no gid
	_ = cmd.Execute()
	assert.Equal(t, len(p.errors), 0)
	assert.Equal(t, len(p.messages), 1)
	assert.Equal(t, p.messages[0].Message, "Success: abc is removed.")
}

func TestRemoveCommand_ValidButThrows(t *testing.T) {

	var p = MockMessenger{}
	var e = MockEngine{}
	e.errorToThrow = errors.New("test error")
	mockDeps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}
	cmd := buildRemoveCommand(mockDeps)
	cmd.SetArgs([]string{"abc"}) // no gid
	_ = cmd.Execute()
	assert.Equal(t, len(p.errors), 1)
	assert.Equal(t, p.errors[0].Message, "test error")
}
