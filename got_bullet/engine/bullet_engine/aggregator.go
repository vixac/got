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
		appliedChange := newParentSummary.ApplyChange(NewCountChange(e.State, true))
		//upsert the parent
		parentId := enrichedEvent.ParentId()
		upserts[*parentId] = appliedChange
	}

	for _, a := range enrichedEvent.Ancestry {
		change := NewCountChange(e.State, true)
		increments[a.Id] = change
	}

	//now we have all the increment maths, we just need to convert it to upserts.
	//apply increments to upserts
	for id, inc := range increments {
		summary, ok := upserts[id]
		if !ok {
			newSummaryCount := summary.ApplyChange(inc)
			upserts[id] = newSummaryCount
		} else {
			//create new upsert
			existingSummary, ok := ancestorAggs[id]
			if !ok {
				fmt.Printf("VX: Error finding id %d\n", id)
				return errors.New("dev error. Summary should exist in agg")
			}
			updatedSummary := existingSummary.ApplyChange(inc)
			upserts[id] = updatedSummary
		}
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
