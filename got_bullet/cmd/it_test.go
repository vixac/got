package cmd

import (
	"errors"
	"testing"

	"gotest.tools/assert"
)

func TestItCommand_MissingValid(t *testing.T) {
	var p = MockMessenger{}
	var e = MockEngine{}
	deps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}

	cmd := buildItCommand(deps)
	cmd.SetArgs([]string{"Finish", "report"})
	err := cmd.Execute()
	assert.NilError(t, err)
	assert.Equal(t, e.heading, "Finish report")
	assert.Equal(t, e.createCompletable, false)
}

func TestItCommand_MissingHeading(t *testing.T) {
	var p = MockMessenger{}
	var e = MockEngine{}
	deps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}

	cmd := buildItCommand(deps)
	cmd.SetArgs([]string{""})
	_ = cmd.Execute()
	assert.Equal(t, p.errors[0].Message, "missing heading")
}

func TestItommand_CreateBuck_ErrorFromEngine(t *testing.T) {
	var p = MockMessenger{}
	var e = MockEngine{}
	e.errorToThrow = errors.New("db error")

	deps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}

	cmd := buildItCommand(deps)
	cmd.SetArgs([]string{"Finish", "report"})
	_ = cmd.Execute()
	assert.Equal(t, len(p.errors), 1)
	assert.Equal(t, p.errors[0].Message, "db error")
	assert.Equal(t, e.createCompletable, false)
}
