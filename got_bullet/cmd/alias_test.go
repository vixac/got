package cmd

import (
	"testing"

	"vixac.com/got/console"
	"vixac.com/got/engine"
)

// --- Mock Dependencies ---

type mockMessenger struct {
	messages []console.Message
	errors   []console.Message
}

func (m *mockMessenger) Print(message console.Message) {
	m.messages = append(m.messages, message)
}
func (m *mockMessenger) Error(message console.Message) {
	m.errors = append(m.errors, message)
}

type mockEngine struct {
}

func (m *mockEngine) Summary(lookup engine.GidLookup) (*engine.GotSummary, error) {
	println("VX: Not mocked Summary")
	return nil, nil
}
func (m *mockEngine) Resolve(lookup engine.GidLookup) (*engine.NodeId, error) {
	println("VX: Not mocked Resolve")
	return nil, nil
}
func (m *mockEngine) Alias(gid string, alias string) (bool, error) {
	println("VX: Not mocked Alias")
	return false, nil
}

// --- Tests ---

func TestDoneCommand_MissingGID(t *testing.T) {

	var p = mockMessenger{}
	var e = mockEngine{}
	mockDeps := RootDependencies{
		Printer: &p,
		Engine:  &e,
	}
	cmd := buildAliasCommand(mockDeps)

	cmd.SetArgs([]string{}) // no --gid flag

	_ = cmd.Execute()

	if len(p.errors) == 0 {
		t.Errorf("expected missing gid message, got no error")
	}
	if len(p.errors) != 1 {
		t.Errorf("expected 1 error")
		return
	}
	if p.errors[0].Message != "Missing gid" {

		t.Errorf("wrong message: %v", p.errors[0].Message)
	}
	//VX:TODO read the error and confirm its correct
}

//unconfirmed.
/*
func TestDoneCommand_Success(t *testing.T) {
	m := &mockMessenger{}
	cmd := cmd.NewDoneCommand(m, tasks)

	cmd.SetArgs([]string{"--gid", "123"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if tasks.completedGID != "123" {
		t.Errorf("expected task to be completed with gid 123, got %v", tasks.completedGID)
	}

	found := false
	for _, msg := range m.output {
		if strings.Contains(msg, "Task completed") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected completion message, got %v", m.output)
	}
}

func TestDoneCommand_FailureFromService(t *testing.T) {
	m := &mockMessenger{}
	tasks := &mockTaskService{err: errors.New("db unavailable")}
	cmd := cmd.NewDoneCommand(m, tasks)

	cmd.SetArgs([]string{"--gid", "123"})

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "db unavailable") {
		t.Errorf("expected service error, got %v", err)
	}

	found := false
	for _, msg := range m.output {
		if strings.Contains(msg, "Failed to complete task") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected error message, got %v", m.output)
	}
}
*/
