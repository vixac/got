package cmd

import (
	"errors"
	"testing"

	"gotest.tools/assert"
)

func TestMoveCommand_MissingGID(t *testing.T) {

	var p = MockMessenger{}
	var e = MockEngine{}
	mockDeps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}
	cmd := buildMvCommand(mockDeps)
	cmd.SetArgs([]string{}) // no --gid flag
	_ = cmd.Execute()
	assert.Equal(t, len(p.errors), 1)
	assert.Equal(t, p.errors[0].Message, "Invalid args. just 2 please.")
}
func TestMoveCommand_MissingParent(t *testing.T) {

	var p = MockMessenger{}
	var e = MockEngine{}
	mockDeps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}
	cmd := buildMvCommand(mockDeps)
	cmd.SetArgs([]string{"item"})
	_ = cmd.Execute()
	assert.Equal(t, len(p.errors), 1)
	assert.Equal(t, p.errors[0].Message, "Invalid args. just 2 please.")
}

func TestMoveCommandValidButEngineThrows(t *testing.T) {

	var p = MockMessenger{}
	var e = MockEngine{}
	e.errorToThrow = errors.New("test error")
	mockDeps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}
	cmd := buildMvCommand(mockDeps)
	cmd.SetArgs([]string{"abc", "new_parent"})
	_ = cmd.Execute()
	assert.Equal(t, len(p.errors), 1)
	assert.Equal(t, p.errors[0].Message, "test error")
	assert.Equal(t, e.moveLookup.Input, "abc")
	assert.Equal(t, e.moveNewParent.Input, "new_parent")
}

func TestMoveCommand_Valid(t *testing.T) {

	var p = MockMessenger{}
	var e = MockEngine{}
	mockDeps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}
	cmd := buildMvCommand(mockDeps)
	cmd.SetArgs([]string{"abc", "new_parent"})
	_ = cmd.Execute()
	assert.Equal(t, len(p.errors), 0)
	assert.Equal(t, e.moveLookup.Input, "abc")
	assert.Equal(t, e.moveNewParent.Input, "new_parent")
	assert.Equal(t, len(p.messages), 1)
	assert.Equal(t, p.messages[0].Message, "Success: abc moved to new parent new_parent")
}

func TestMoveCommand_ValidButWithExistingParent(t *testing.T) {

	var p = MockMessenger{}
	var e = MockEngine{}
	mockDeps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}
	cmd := buildMvCommand(mockDeps)
	cmd.SetArgs([]string{"abc", "new_parent"})
	_ = cmd.Execute()
	if len(p.errors) != 0 {
		t.Errorf("expected  no errors. 0 != %d", len(p.errors))
		return
	}
	assert.Equal(t, e.moveLookup.Input, "abc")
	assert.Equal(t, e.moveNewParent.Input, "new_parent")
	assert.Equal(t, len(p.messages), 1)
	assert.Equal(t, p.messages[0].Message, "Success: abc moved to new parent new_parent")
}
