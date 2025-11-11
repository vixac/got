package cmd

import (
	"errors"
	"testing"

	"gotest.tools/assert"
)

func TestTillCommand_MissingAlias(t *testing.T) {
	var p = MockMessenger{}
	var e = MockEngine{}
	deps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}

	cmd := buildTillCommand(deps)
	cmd.SetArgs([]string{"2025-10-22", "Finish", "report"}) // no --for flag

	_ = cmd.Execute()
	assert.Equal(t, len(p.errors), 1)
	assert.Equal(t, p.errors[0].Message, "missing alias")
}

func TestTillCommand_MissingDate(t *testing.T) {
	var p = MockMessenger{}
	var e = MockEngine{}
	deps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}

	cmd := buildTillCommand(deps)
	cmd.SetArgs([]string{"--for", "parentAlias"}) // missing date + heading
	_ = cmd.Execute()
	assert.Equal(t, p.errors[0].Message, "missing args")
}

func TestTillCommand_MissingHeading(t *testing.T) {
	var p = MockMessenger{}
	var e = MockEngine{}
	deps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}

	cmd := buildTillCommand(deps)
	cmd.SetArgs([]string{"--for", "parentAlias", "2025-10-22"}) // no heading
	_ = cmd.Execute()
	assert.Equal(t, p.errors[0].Message, "missing args")
}

func TestTillCommand_CreateBuck_ErrorFromEngine(t *testing.T) {
	var p = MockMessenger{}
	var e = MockEngine{}
	e.errorToThrow = errors.New("db error")

	deps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}

	cmd := buildTillCommand(deps)
	cmd.SetArgs([]string{"--for", "parentAlias", "2025-10-22", "Finish", "report"})
	_ = cmd.Execute()
	assert.Equal(t, len(p.errors), 1)
	assert.Equal(t, p.errors[0].Message, "db error")

	// Confirm the engine was called correctly
	assert.Equal(t, e.createParent.Input, "parentAlias")
	assert.Equal(t, e.createDate.UserInput, "2025-10-22")
	assert.Equal(t, e.createCompletable, true)
}

func TestTillCommand_Valid(t *testing.T) {
	var p = MockMessenger{}
	var e = MockEngine{}
	deps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}

	cmd := buildTillCommand(deps)
	cmd.SetArgs([]string{"--for", "parentAlias", "2025-10-22", "Finish", "report"})
	_ = cmd.Execute()

	// Ensure engine was called with the expected parameters
	assert.Equal(t, e.createParent.Input, "parentAlias")
	assert.Equal(t, e.createDate.UserInput, "2025-10-22")
	assert.Equal(t, e.createCompletable, true)

	// Should not print any errors
	assert.Equal(t, len(p.errors), 0)
}
