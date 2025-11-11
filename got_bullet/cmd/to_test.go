package cmd

import (
	"errors"
	"testing"

	"gotest.tools/assert"
)

func TestToCommand_MissingValid(t *testing.T) {
	var p = MockMessenger{}
	var e = MockEngine{}
	deps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}

	cmd := buildToCommand(deps)
	cmd.SetArgs([]string{"Finish", "report"})
	err := cmd.Execute()
	assert.NilError(t, err)
	assert.Equal(t, e.heading, "Finish report")
	assert.Equal(t, e.createCompletable, true)
}

func TestToCommand_MissingHeading(t *testing.T) {
	var p = MockMessenger{}
	var e = MockEngine{}
	deps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}

	cmd := buildToCommand(deps)
	cmd.SetArgs([]string{""})
	_ = cmd.Execute()
	assert.Equal(t, p.errors[0].Message, "missing heading")
}

func TestToCommand_CreateBuck_ErrorFromEngine(t *testing.T) {
	var p = MockMessenger{}
	var e = MockEngine{}
	e.errorToThrow = errors.New("db error")

	deps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}

	cmd := buildToCommand(deps)
	cmd.SetArgs([]string{"Finish", "report"})
	_ = cmd.Execute()
	assert.Equal(t, len(p.errors), 1)
	assert.Equal(t, p.errors[0].Message, "db error")
	assert.Equal(t, e.createCompletable, true)
}
