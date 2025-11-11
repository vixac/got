package cmd

import (
	"errors"
	"testing"

	"gotest.tools/assert"
)

func TestNoteCommand_MissingArgs(t *testing.T) {
	var p = MockMessenger{}
	var e = MockEngine{}
	deps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}

	cmd := buildNoteCommand(deps)
	cmd.SetArgs([]string{}) // no args at all

	err := cmd.Execute()
	assert.ErrorContains(t, err, "missing args")
	assert.Equal(t, len(p.errors), 1)
	assert.Equal(t, p.errors[0].Message, "missing args")
}

func TestNoteCommand_MissingAlias(t *testing.T) {
	var p = MockMessenger{}
	var e = MockEngine{}
	deps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}

	cmd := buildNoteCommand(deps)
	cmd.SetArgs([]string{"", "Write", "notes"}) // empty alias

	err := cmd.Execute()
	assert.ErrorContains(t, err, "missing alias")
	assert.Equal(t, len(p.errors), 1)
	assert.Equal(t, p.errors[0].Message, "missing alias")
}

func TestNoteCommand_MissingHeading(t *testing.T) {
	var p = MockMessenger{}
	var e = MockEngine{}
	deps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}

	cmd := buildNoteCommand(deps)
	cmd.SetArgs([]string{"parentAlias"}) // no heading provided

	err := cmd.Execute()
	assert.ErrorContains(t, err, "missing args")
	assert.Equal(t, len(p.errors), 1)
	assert.Equal(t, p.errors[0].Message, "missing args")
}

func TestNoteCommand_CreateBuck_ErrorFromEngine(t *testing.T) {
	var p = MockMessenger{}
	var e = MockEngine{}
	e.errorToThrow = errors.New("db error")

	deps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}

	cmd := buildNoteCommand(deps)
	cmd.SetArgs([]string{"parentAlias", "Write", "notes"})
	err := cmd.Execute()
	assert.ErrorContains(t, err, "db error")
	assert.Equal(t, len(p.errors), 1)
	assert.Equal(t, p.errors[0].Message, "db error")

	// Confirm the engine call
	assert.Equal(t, e.createParent.Input, "parentAlias")
	assert.Assert(t, e.createDate == nil)
	assert.Equal(t, e.createCompletable, false)
}

func TestNoteCommand_Valid(t *testing.T) {
	var p = MockMessenger{}
	var e = MockEngine{}
	deps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}

	cmd := buildNoteCommand(deps)
	cmd.SetArgs([]string{"parentAlias", "Write", "notes"})
	err := cmd.Execute()
	assert.NilError(t, err)

	// Engine should be called correctly
	assert.Equal(t, e.createParent.Input, "parentAlias")
	assert.Assert(t, e.createDate == nil)
	assert.Equal(t, e.createCompletable, false)

	// Should not print errors
	assert.Equal(t, len(p.errors), 0)
}
