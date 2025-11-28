package bullet_engine

import (
	"fmt"

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

// the idea here is that this listens to events and propagates updates to the aggregation
type BulletAggListener struct {
	store SummaryStoreInterface
}

func NewBulletAggListener(store SummaryStoreInterface) EventListenerInterface {
	return &BulletAggListener{
		store: store,
	}
}

// / LISTENER EVENTS
func (b *BulletAggListener) ItemAdded(e AddItemEvent) error {
	fmt.Printf("VX:TODO unhandled event ")
	return nil
}

func (b *BulletAggListener) ItemStateChanged(e StateChangeEvent) error {
	fmt.Printf("VX:TODO unhandled event ")
	return nil
}
func (b *BulletAggListener) ItemDeleted(e ItemDeletedEvent) error {
	fmt.Printf("VX:TODO unhandled event ")
	return nil
}
func (b *BulletAggListener) ItemMoved(e ItemMovedEvent) error {
	fmt.Printf("VX:TODO unhandled event ")
	return nil
}

func (b *BulletAggListener) onAdd(e AddItemEvent) error {

	allAncestors, err := b.store.Fetch(e.Ancestry)
	if err != nil {
		return err
	}
	for aggId, a := range allAncestors {
		newCount := allAncestors[aggId].Counts.ChangeState(e.State, 1)
		newAnc := a.UpdatedCount(newCount)
		b.store.UpsertAggregate(aggId, newAnc)

	}
	//each new item gets an empty aggregate.
	b.store.UpsertAggregate(e.Id, Aggregate{})
	return nil
}
