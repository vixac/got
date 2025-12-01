package bullet_engine

import (
	"errors"
	"fmt"
)

/*
*
The aggregator is going to contain the business logic that maps events to
changes in the aggstore.
*/

type Aggregator struct {
	summaryStore SummaryStoreInterface
	// ancestorStore AncestorListInterface
}

func NewAggregator(summaryStore SummaryStoreInterface) (*Aggregator, error) {
	return &Aggregator{
		summaryStore: summaryStore,
	}, nil
}

/**
 */

func (a *Aggregator) ItemAdded(e AddItemEvent) error {

	ancestorAggs, err := a.summaryStore.Fetch(e.Ancestry)
	if err != nil {
		return err
	}

	enrichedEvent, err := NewEnrichedAddItemEvent(e, ancestorAggs)
	if err != nil {
		return err
	}

	//step 1. We create the new summary object for the new item
	upserts := make(map[SummaryId]Summary)
	upserts[e.Id] = NewLeafSummary(e.State, e.Deadline)
	//here we walk through the notion table: https://www.notion.so/Summary-2b69775b667e804886a8caafc3497136

	//some counters or somethign?
	//VX:TODO increment stateCount on all ancestors

	increments := make(map[SummaryId]AggregateCountChange)

	if enrichedEvent.ParentIsLeaf() {
		//convert parent to group with a count 1 for e.state
		//decrement all aggs with the parent state
		//time to convert it.
		parentState := enrichedEvent.ParentState()
		if parentState == nil {
			return errors.New("missing dev state")
		}
		newParentSummary := Summary{
			State:    nil,
			Deadline: enrichedEvent.Parent().Summary.Deadline,
		}
		//newParentSummary.ApplyChange(NewCountChange(e.State, true))
		//upsert the parent
		parentId := enrichedEvent.ParentId()
		fmt.Printf("VX: Leaf parent is changed to %+v -> %+v \n", newParentSummary, *enrichedEvent.Parent())
		upserts[*parentId] = newParentSummary

		//decrement the parents state on all ancestors
		for _, a := range enrichedEvent.Ancestry {
			if a.Id != *parentId {
				change := NewCountChange(*parentState, false)
				fmt.Printf("VX: because a leaf changed to group, we are decrementing")
				increments[a.Id] = change
			}
		}
	}
	for _, u := range upserts {
		if u.Counts != nil {
			fmt.Printf("VX: here is an upsert we need to insert before we do the addition: %+v\n", u.Counts)
		} else {
			fmt.Printf("VX: this upsert had no count: %+v\n", u)
		}

	}

	//VX:TODOT HIS BITIS WRONG

	//increment all parents with the new state
	for _, a := range enrichedEvent.Ancestry {
		change := NewCountChange(e.State, true)
		//existingUpsert, ok := upserts[a.Id]
		existingIncrement, ok := increments[a.Id]
		if ok { //update the existing upsert
			existingUpsert.ApplyChange(change)
			upserts[a.Id] = existingUpsert //put the change straight in
		} else {
			//store the change as an upsert for later? Why thogh
			increments[a.Id] = change
		}

	}

	//now we have all the increment maths, we just need to convert it to upserts.
	//apply increments to upserts
	for id, inc := range increments {
		existingSummary, ok := ancestorAggs[id]
		if !ok {
			fmt.Printf("VX: Error finding id %d\n", id)
			return errors.New("dev error. Summary should exist in agg")
		}

		existingUpsert, ok := upserts[id]
		if !ok { //no upsert to edit. So we create a version of the existing with the increment
			fmt.Printf("VX: no existingincrement for .. %d\n", id)
			existingSummary.ApplyChange(inc)
			if existingUpsert.State == nil {
				fmt.Printf("VX: this should still be aleaf: %+v \n", existingSummary)
			}
			upserts[id] = existingSummary
		} else { //we just apply the change to the upsert
			fmt.Printf("VX: existing upsert.. %+v\n", existingUpsert)
			//create new upsert

			existingUpsert.ApplyChange(inc)
			upserts[id] = existingUpsert
		}
	}
	for _, a := range upserts {
		fmt.Printf("VX: upserting this %+v\n", a)
	}

	//apply all the created changes.
	return a.summaryStore.UpsertManyAggregates(upserts)
}

func (a *Aggregator) ItemStateChanged(e StateChangeEvent) error {
	fmt.Printf("VX:TODO unhandled event ")
	return nil
}
func (a *Aggregator) ItemDeleted(e ItemDeletedEvent) error {
	fmt.Printf("VX:TODO unhandled event ")
	return nil
}

func (a *Aggregator) ItemMoved(e ItemMovedEvent) error {
	fmt.Printf("VX:TODO unhandled event ")
	return nil
}
