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
	Id       SummaryId
	State    engine.GotState
	Ancestry []SummaryId
	Deadline *Deadline
}
type StateChangeEvent struct {
	Id       SummaryId
	OldState engine.GotState
	NewState engine.GotState
}

type ItemDeletedEvent struct {
	Id SummaryId
}
type ItemMovedEvent struct {
	Id        SummaryId
	OldParent *SummaryId
	NewParent *SummaryId
}

// VX:TODO flesh this out with all the events that might be interesting
type EventListenerInterface interface {
	ItemAdded(e AddItemEvent) error
	ItemStateChanged(e StateChangeEvent) error
	ItemDeleted(e ItemDeletedEvent) error
	ItemMoved(e ItemMovedEvent) error
}
