package cmd

import (
	"vixac.com/got/console"
	"vixac.com/got/engine"
)

// --- Mock Dependencies ---

type MockMessenger struct {
	messages []console.Message
	errors   []console.Message
}

func (m *MockMessenger) Print(message console.Message) {
	m.messages = append(m.messages, message)
}
func (m *MockMessenger) Error(message console.Message) {
	m.errors = append(m.errors, message)
}

type MockEngine struct {
	aliasGid   string
	aliasAlias string

	summaryLookup engine.GidLookup
	resolveLookup engine.GidLookup
	errorToThrow  error

	unaliasAlias string
}

func (m *MockEngine) Unalias(alias string) (*engine.NodeId, error) {
	m.unaliasAlias = alias
	return nil, m.errorToThrow
}
func (m *MockEngine) Summary(lookup engine.GidLookup) (*engine.GotSummary, error) {
	m.summaryLookup = lookup
	return nil, m.errorToThrow
}
func (m *MockEngine) Resolve(lookup engine.GidLookup) (*engine.NodeId, error) {
	m.resolveLookup = lookup
	return nil, m.errorToThrow
}
func (m *MockEngine) Delete(lookup engine.GidLookup) (*engine.NodeId, error) {
	m.resolveLookup = lookup
	return nil, m.errorToThrow
}
func (m *MockEngine) Alias(gid string, alias string) (bool, error) {
	m.aliasGid = gid
	m.aliasAlias = alias
	return false, m.errorToThrow
}
