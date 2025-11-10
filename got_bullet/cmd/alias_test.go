package cmd

import (
	"errors"
	"testing"

	"gotest.tools/assert"
)

func TestAliasCommand_MissingGID(t *testing.T) {

	var p = MockMessenger{}
	var e = MockEngine{}
	mockDeps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}
	cmd := buildAliasCommand(mockDeps)
	cmd.SetArgs([]string{}) // no --gid flag
	_ = cmd.Execute()
	assert.Equal(t, len(p.errors), 1)

	if p.errors[0].Message != "Missing gid" {
		t.Errorf("wrong message: %v", p.errors[0].Message)
	}
}
func TestAliasCommand_MissingAlias(t *testing.T) {

	var p = MockMessenger{}
	var e = MockEngine{}
	mockDeps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}
	cmd := buildAliasCommand(mockDeps)
	cmd.SetArgs([]string{"--gid", "abc"}) // no --alias flag
	_ = cmd.Execute()
	if len(p.errors) != 1 {
		t.Errorf("expected 1 error")
		return
	}
	if p.errors[0].Message != "Missing alias" {
		t.Errorf("wrong message: %v", p.errors[0].Message)
	}
}

func TestAliasCommandValidButEngineThrows(t *testing.T) {

	var p = MockMessenger{}
	var e = MockEngine{}
	e.errorToThrow = errors.New("test error")
	mockDeps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}
	cmd := buildAliasCommand(mockDeps)
	cmd.SetArgs([]string{"--gid", "abc", "--alias", "new_name"}) // no --alias flag
	_ = cmd.Execute()
	assert.Equal(t, len(p.errors), 1)
	assert.Equal(t, p.errors[0].Message, "test error")
	assert.Equal(t, e.aliasAlias, "new_name")
	assert.Equal(t, e.aliasGid, "abc")
}

func TestAliasCommand_Valid(t *testing.T) {

	var p = MockMessenger{}
	var e = MockEngine{}
	mockDeps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}
	cmd := buildAliasCommand(mockDeps)
	cmd.SetArgs([]string{"--gid", "abc", "--alias", "new_name"}) // no --alias flag
	_ = cmd.Execute()
	if len(p.errors) != 0 {
		t.Errorf("expected  no errors. 0 != %d", len(p.errors))
		return
	}
	assert.Equal(t, e.aliasAlias, "new_name")
	assert.Equal(t, e.aliasGid, "abc")
	assert.Equal(t, len(p.messages), 1)
	assert.Equal(t, p.messages[0].Message, "Success: abc is now aliased to new_name.")
}
