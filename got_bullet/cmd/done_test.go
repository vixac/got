package cmd

import (
	"errors"
	"testing"

	"gotest.tools/assert"
)

func TestDoneCommand_MissingGID(t *testing.T) {

	var p = MockMessenger{}
	var e = MockEngine{}
	mockDeps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}
	cmd := buildDoneCommand(mockDeps)
	cmd.SetArgs([]string{}) // no gid
	_ = cmd.Execute()
	assert.Equal(t, len(p.errors), 1)

	assert.Equal(t, p.errors[0].Message, "Expected at least one lookup as input")
}

func TestDoneCommand_Valid(t *testing.T) {

	var p = MockMessenger{}
	var e = MockEngine{}
	mockDeps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}
	cmd := buildDoneCommand(mockDeps)
	cmd.SetArgs([]string{"abc"}) // no gid
	_ = cmd.Execute()
	assert.Equal(t, len(p.errors), 0)
	assert.Equal(t, len(p.messages), 1)
	assert.Equal(t, p.messages[0].Message, "Success: 1 items is marked complete.")
}

func TestDoneCommand_ValidButThrows(t *testing.T) {

	var p = MockMessenger{}
	var e = MockEngine{}
	e.errorToThrow = errors.New("test error")
	mockDeps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}
	cmd := buildDoneCommand(mockDeps)
	cmd.SetArgs([]string{"abc"}) // no gid
	_ = cmd.Execute()
	assert.Equal(t, len(p.errors), 1)
	assert.Equal(t, p.errors[0].Message, "test error")
}
