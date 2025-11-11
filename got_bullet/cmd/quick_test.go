package cmd

import (
	"errors"
	"testing"

	"gotest.tools/assert"
)

func TestQuickCommand_MissingValid(t *testing.T) {
	var p = MockMessenger{}
	var e = MockEngine{}
	deps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}

	cmd := buildQuickCommand(deps)
	cmd.SetArgs([]string{"2025-10-22", "Finish", "report"}) // no --for flag
	err := cmd.Execute()
	assert.NilError(t, err)
	assert.Equal(t, e.createDate.UserInput, "2025-10-22")
	assert.Equal(t, e.heading, "Finish report")
	assert.Equal(t, e.createCompletable, true)
}

func TestQuickCommand_MissingDate(t *testing.T) {
	var p = MockMessenger{}
	var e = MockEngine{}
	deps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}

	cmd := buildQuickCommand(deps)
	cmd.SetArgs([]string{"", "Finish", "report"}) // missing date + heading
	_ = cmd.Execute()
	assert.Equal(t, p.errors[0].Message, "missing date")
}

func TestQuickCommand_MissingHeading(t *testing.T) {
	var p = MockMessenger{}
	var e = MockEngine{}
	deps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}

	cmd := buildQuickCommand(deps)
	cmd.SetArgs([]string{"2025-10-22"}) // no heading
	_ = cmd.Execute()
	assert.Equal(t, p.errors[0].Message, "missing args")
}

func TestQuickCommand_CreateBuck_ErrorFromEngine(t *testing.T) {
	var p = MockMessenger{}
	var e = MockEngine{}
	e.errorToThrow = errors.New("db error")

	deps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}

	cmd := buildQuickCommand(deps)
	cmd.SetArgs([]string{"2025-10-22", "Finish", "report"})
	_ = cmd.Execute()
	assert.Equal(t, len(p.errors), 1)
	assert.Equal(t, p.errors[0].Message, "db error")

	// Confirm the engine was called correctly
	assert.Equal(t, e.createDate.UserInput, "2025-10-22")
	assert.Equal(t, e.createCompletable, true)
}
