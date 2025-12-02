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
}

func NewAggregator(summaryStore SummaryStoreInterface) (*Aggregator, error) {
	return &Aggregator{
		summaryStore: summaryStore,
	}, nil
}

func (a *Aggregator) ItemAdded(e AddItemEvent) error {

	ancestorAggs, err := a.summaryStore.Fetch(e.Ancestry)
	if err != nil {
		return err
	}

	enrichedEvent, err := EnrichSummaries(e.Ancestry, ancestorAggs)
	if err != nil {
		return err
	}

	//step 1. We create the new summary object for the new item
	upserts := make(map[SummaryId]Summary)
	upserts[e.Id] = NewLeafSummary(e.State, e.Deadline)
	//here we walk through the notion table: https://www.notion.so/Summary-2b69775b667e804886a8caafc3497136
	if enrichedEvent.ParentIsLeaf() {
		//convert parent to group with a count 1 for e.state
		parentState := enrichedEvent.ParentState()
		if parentState == nil {
			return errors.New("missing dev state")
		}
		newParentSummary := Summary{
			State:    nil,
			Deadline: enrichedEvent.Parent().Summary.Deadline,
		}
		parentId := enrichedEvent.ParentId()
		fmt.Printf("VX: Leaf parent is changed to %+v from original %+v \n", newParentSummary, *enrichedEvent.Parent())
		upserts[*parentId] = newParentSummary

		//decrement the parents state on all ancestors
		for _, a := range enrichedEvent.Ancestry {
			if a.Id != *parentId {
				change := NewCountChange(*parentState, false)
				fmt.Printf("VX: because a leaf changed to group, we are decrementing")
				a.Summary.ApplyChange(change)
				upserts[a.Id] = a.Summary
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

	//increment all parents with the new state
	change := NewCountChange(e.State, true)
	for _, a := range enrichedEvent.Ancestry {
		existingUpsert, ok := upserts[a.Id]
		if ok {
			fmt.Printf("VX: we have an upsert for this one already: %d, %+v", a.Id, existingUpsert)
			existingUpsert.ApplyChange(change)
			upserts[a.Id] = existingUpsert
			fmt.Printf("VX: summary is now: %d, %+v", a.Id, existingUpsert)
		} else {
			fmt.Printf("VX: icnremting for the first time: %+v\n", a.Summary.Counts)
			a.Summary.ApplyChange(change)
			fmt.Printf("VX: incremened for the first time is now: %+v\n", a.Summary.Counts)
			upserts[a.Id] = a.Summary
		}
	}

	//apply all the created changes.
	return a.summaryStore.UpsertManyAggregates(upserts)
}

// VX:TODO this is ready to try out.
func (a *Aggregator) ItemStateChanged(e StateChangeEvent) error {

	idsIncludingThis := e.Ancestry
	idsIncludingThis = append(idsIncludingThis, e.Id) //the last item is *THIS*, it's on the end which is wierd.
	ancestorAggs, err := a.summaryStore.Fetch(idsIncludingThis)
	if err != nil {
		return err
	}
	//step 1. change the state of this leaf.
	changedItemSummary, ok := ancestorAggs[e.Id]
	if !ok {
		return errors.New("missing summary for state-changed item.s")
	}
	changedItemSummary.State = &e.NewState
	upserts := make(map[SummaryId]Summary)
	upserts[e.Id] = changedItemSummary

	//step 1  we decrement the old state and increment the new for all its ancestors
	incChange := NewCountChange(e.NewState, true)
	decChange := NewCountChange(e.OldState, false)
	combined := incChange.Combine(decChange)
	for _, summaryId := range e.Ancestry {
		summary, ok := ancestorAggs[summaryId]
		if !ok {
			return errors.New("missing summary in state-change for ancestor")
		}
		summary.ApplyChange(combined)
		upserts[summaryId] = summary
	}
	return a.summaryStore.UpsertManyAggregates(upserts)
}

func (a *Aggregator) ItemDeleted(e ItemDeletedEvent) error {
	fmt.Printf("VX:TODO unhandled event ")
	return nil
}

func (a *Aggregator) ItemMoved(e ItemMovedEvent) error {
	fmt.Printf("VX:TODO unhandled event ")
	return nil
}
