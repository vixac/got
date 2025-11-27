package bullet_engine

import (
	"fmt"

	"vixac.com/got/engine"
)

type ItemEventType int

// VX:TODO these are the event types
const (
	EventTypeAdd         = 100
	EventTypeChangeState = 101
)

// VX:TODO ADD EVENT
type ItemEvent struct {
	Type     ItemEventType
	Id       AggId
	State    engine.GotState
	Ancestry []AggId
	Deadline *Deadline
}

// VX:TODO flesh this out with all the events that might be interesting
type AggListenerInterface interface {
	ItemEvent(e ItemEvent) error
}

// the idea here is that this listens to events and propagates updates to the aggregation
type BulletAggListener struct {
	store AggStoreInterface
}

func NewBulletAggListener(store AggStoreInterface) AggListenerInterface {
	return &BulletAggListener{
		store: store,
	}
}

func (b *BulletAggListener) ItemEvent(e ItemEvent) error {
	if e.Type == EventTypeAdd {
		return b.onAdd(e)
	}
	fmt.Printf("VX:TODO unhandled event ")
	return nil
}

func (b *BulletAggListener) onAdd(e ItemEvent) error {

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
