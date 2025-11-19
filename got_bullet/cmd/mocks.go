package cmd

import (
	"errors"

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

	nodeIdToReturn *engine.NodeId
	unaliasAlias   string

	moveLookup    engine.GidLookup
	moveNewParent engine.GidLookup

	// New fields for till command
	createParent      *engine.GidLookup
	createDate        *engine.DateLookup
	createCompletable bool
	heading           string
}

func (m *MockEngine) Unalias(alias string) (*engine.NodeId, error) {
	m.unaliasAlias = alias
	return m.nodeIdToReturn, m.errorToThrow
}
func (m *MockEngine) Summary(lookup engine.GidLookup) (*engine.GotSummary, error) {
	m.summaryLookup = lookup
	return nil, m.errorToThrow
}
func (m *MockEngine) Resolve(lookup engine.GidLookup) (*engine.NodeId, error) {
	m.resolveLookup = lookup
	return m.nodeIdToReturn, m.errorToThrow
}
func (m *MockEngine) Delete(lookup engine.GidLookup) (*engine.NodeId, error) {
	m.resolveLookup = lookup
	return m.nodeIdToReturn, m.errorToThrow
}
func (m *MockEngine) Alias(gid string, alias string) (bool, error) {
	m.aliasGid = gid
	m.aliasAlias = alias
	return false, m.errorToThrow
}

func (m *MockEngine) Lookup(alias string) (*engine.NodeId, error) {
	return nil, errors.New("not impl")

}
func (m *MockEngine) Move(lookup engine.GidLookup, newParent engine.GidLookup) (*engine.NodeId, error) {
	m.moveLookup = lookup
	m.moveNewParent = newParent
	return m.nodeIdToReturn, m.errorToThrow
}

func (m *MockEngine) CreateBuck(parent *engine.GidLookup, date *engine.DateLookup, completable bool, heading string) (*engine.NodeId, error) {
	m.createParent = parent
	m.createDate = date
	m.createCompletable = completable
	m.heading = heading
	return m.nodeIdToReturn, m.errorToThrow
}
