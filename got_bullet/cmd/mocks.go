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
func (m *MockMessenger) PrintInLine(line []console.Message) {
	//m.messages = append(m.messages, message)
}
func (m *MockMessenger) Error(message console.Message) {
	m.errors = append(m.errors, message)
}

type MockEngine struct {
	aliasGid   string
	aliasAlias string

	summaryLookup *engine.GidLookup
	doneLookup    []engine.GidLookup
	resolveLookup engine.GidLookup
	errorToThrow  error

	nodeIdToReturn *engine.NodeId
	gotIdToReturn  *engine.GotId
	unaliasAlias   string

	moveLookup    engine.GidLookup
	moveNewParent engine.GidLookup

	// New fields for till command
	createParent      *engine.GidLookup
	createDate        *engine.DateLookup
	createCompletable bool
	heading           string
}

func (m *MockEngine) OpenThenTimestamp(lookup engine.GidLookup) error {
	return errors.New("not impl")
}
func (m *MockEngine) TagItem(lookup engine.GidLookup, tag engine.TagLookup) error {
	return errors.New("not impl")
}

func (e MockEngine) ScheduleItem(lookup engine.GidLookup, dateLookup engine.DateLookup) error {
	return errors.New("not impl")
}
func (e *MockEngine) FetchItemsBelow(lookup *engine.GidLookup, sortByPath bool, states []engine.GotState) (*engine.GotFetchResult, error) {
	return nil, errors.New("not impl")
}

func (m *MockEngine) LookupAliasForMany(gid []string) (map[string]*string, error) {
	return nil, errors.New("not impl")
}

func (m *MockEngine) Unalias(alias string) (*engine.GotId, error) {
	m.unaliasAlias = alias
	return m.gotIdToReturn, m.errorToThrow
}
func (m *MockEngine) Summary(lookup *engine.GidLookup) (*engine.GotItemDisplay, error) {
	m.summaryLookup = lookup
	return nil, m.errorToThrow
}

func (m *MockEngine) MarkResolved(lookup []engine.GidLookup) error {
	m.doneLookup = lookup
	return m.errorToThrow
}
func (m *MockEngine) MarkActive(lookup engine.GidLookup) (*engine.NodeId, error) {
	m.resolveLookup = lookup
	return m.nodeIdToReturn, m.errorToThrow
}
func (m *MockEngine) MarkAsNote(lookup engine.GidLookup) (*engine.NodeId, error) {
	m.resolveLookup = lookup
	return m.nodeIdToReturn, m.errorToThrow
}

func (m *MockEngine) Delete(lookup engine.GidLookup) error {
	m.resolveLookup = lookup
	return m.errorToThrow
}
func (m *MockEngine) Alias(lookup engine.GidLookup, alias string) (bool, error) {
	m.resolveLookup = lookup
	m.aliasAlias = alias
	return false, m.errorToThrow
}

func (m *MockEngine) EditTitle(lookup engine.GidLookup, newHeading string) error {
	return m.errorToThrow
}

func (m *MockEngine) LookupAliasForGid(gid string) (*string, error) {
	return nil, errors.New("not impl")
}

func (m *MockEngine) Lookup(alias string) (*engine.GotId, error) {
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
