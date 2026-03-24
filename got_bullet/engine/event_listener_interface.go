package engine

// VX:TODO can we rename Agg entirely to state. and consider the aggregation to be an extension of state.
type ItemEventType int

// VX:TODO these are the event types
const (
	EventTypeAdd         = 100
	EventTypeChangeState = 101
)

type AddItemEvent struct {
	Id               SummaryId
	State            GotState
	Ancestry         []SummaryId
	Deadline         *DateTime
	OverrideSettings *CreateOverrideSettings
}
type StateChangeEvent struct {
	Id       SummaryId
	OldState GotState
	NewState *GotState //here we pass nil if the item was removed.
	Ancestry []SummaryId
}

type ItemDeletedEvent struct {
	Id       SummaryId
	State    GotState
	Ancestry []SummaryId
}
type ItemMovedEvent struct {
	Id          SummaryId
	OldAncestry []SummaryId
	NewAncestry []SummaryId
}

type EditItemEvent struct {
	Id SummaryId
}

// VX:TODO flesh this out with all the events that might be interesting
type EventListenerInterface interface {
	ItemAdded(e AddItemEvent) error
	ItemStateChanged(e StateChangeEvent) error
	ItemDeleted(e ItemDeletedEvent) error
	ItemMoved(e ItemMovedEvent) error
	ItemEdited(e EditItemEvent) error
}
