package bullet_engine

import (
	"vixac.com/got/engine"
)

// VX:TODO can we rename Agg entirely to state. and consider the aggregation to be an extension of state.
type ItemEventType int

// VX:TODO these are the event types
const (
	EventTypeAdd         = 100
	EventTypeChangeState = 101
)

// VX:TODO ADD EVENT
type AddItemEvent struct {
	Id       engine.SummaryId
	State    engine.GotState
	Ancestry []engine.SummaryId
	Deadline *engine.DateTime
}
type StateChangeEvent struct {
	Id       engine.SummaryId
	OldState engine.GotState
	NewState *engine.GotState //here we pass nil if the item was removed.
	Ancestry []engine.SummaryId
}

type ItemDeletedEvent struct {
	Id       engine.SummaryId
	State    engine.GotState
	Ancestry []engine.SummaryId
}
type ItemMovedEvent struct {
	Id        engine.SummaryId
	OldParent *engine.SummaryId
	NewParent *engine.SummaryId
}

type EditItemEvent struct {
	Id engine.SummaryId
}

// VX:TODO flesh this out with all the events that might be interesting
type EventListenerInterface interface {
	ItemAdded(e AddItemEvent) error
	ItemStateChanged(e StateChangeEvent) error
	ItemDeleted(e ItemDeletedEvent) error
	ItemMoved(e ItemMovedEvent) error
	ItemEdited(e EditItemEvent) error
}
