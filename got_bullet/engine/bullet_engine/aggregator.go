package bullet_engine

import (
	"errors"
	"fmt"
	"time"

	"vixac.com/got/engine"
)

/*
*
The aggregator is going to contain the business logic that maps events to
changes in the aggstore.
*/

type Aggregator struct {
	summaryStore SummaryStoreInterface
}

type AggregatorInterface interface {
	ItemAdded(e AddItemEvent) error
	ItemStateChanged(e StateChangeEvent) error
	ItemDeleted(e ItemDeletedEvent) error
	ItemMoved(e ItemMovedEvent) error
	ItemEdited(e EditItemEvent) error
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
	upserts := make(map[engine.SummaryId]engine.Summary)
	upserts[e.Id] = engine.NewLeafSummary(e.State, e.Deadline, time.Now(), []engine.Tag{})

	//here we walk through the notion table: https://www.notion.so/Summary-2b69775b667e804886a8caafc3497136
	if enrichedEvent.ParentIsLeaf() {

		parentIndex := len(e.Ancestry) - 1
		parentSummaryId := e.Ancestry[parentIndex]
		parentSummary, ok := ancestorAggs[parentSummaryId]
		if !ok {
			return errors.New("missing summary for parent of state-changed items")
		}

		parentState := enrichedEvent.ParentState() //we need this to decrement its state from its ancestors.
		if parentState == nil {
			return errors.New("missing dev state. The parent should have had a state at the moment ")
		}
		parentSummary.State = nil //this is how we state that this is now a group.
		parentId := enrichedEvent.ParentId()
		upserts[*parentId] = parentSummary

		//decrement the parents state on all ancestors
		for _, a := range enrichedEvent.Ancestry {
			if a.Id != *parentId {
				change := engine.NewCountChange(*parentState, false)
				a.Summary.ApplyChange(change)

				//if we've added an active item then all its parents are deactivated
				if e.State == engine.Active && a.Summary.State != nil && *a.Summary.State == engine.Active {
					a.Summary.State = nil

				}
				upserts[a.Id] = a.Summary
			}
		}
	}
	//increment all parents with the new state
	change := engine.NewCountChange(e.State, true)
	for _, a := range enrichedEvent.Ancestry {
		existingUpsert, ok := upserts[a.Id]
		if ok {
			existingUpsert.ApplyChange(change)
			upserts[a.Id] = existingUpsert
		} else {
			a.Summary.ApplyChange(change)
			upserts[a.Id] = a.Summary
		}
	}

	//apply all the created changes.
	return a.summaryStore.UpsertManySummaries(upserts)
}

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
	changedItemSummary.State = e.NewState

	nowTime, err := engine.NewDateTime(time.Now())
	if err != nil {
		return err
	}
	changedItemSummary.UpdatedDate = &nowTime
	upserts := make(map[engine.SummaryId]engine.Summary)
	upserts[e.Id] = changedItemSummary

	hasAncestors := len(e.Ancestry) > 0
	if !hasAncestors { //just upsert this item and move on.
		return a.summaryStore.UpsertManySummaries(upserts)
	}

	parentIndex := len(e.Ancestry) - 1
	parentSummaryId := e.Ancestry[parentIndex]
	if parentSummaryId == engine.SummaryId(TheRootNoteInt32) {
		return a.summaryStore.UpsertManySummaries(upserts)
	}
	parentSummary, ok := ancestorAggs[parentSummaryId]
	if !ok {
		return errors.New("missing summary for parent of state-changed items")
	}

	//step 1  we decrement the old state and increment the new for all its ancestors
	//we check that the old state was active and the new state isnt active.
	isChangeFromActive := e.OldState == engine.Active && (e.NewState == nil || *e.NewState != engine.Active)

	//theres a chance we need to convert the parent to active state.
	isParentInNeedOfPromotingToAcive := isChangeFromActive && parentSummary.Counts.Active == 1

	parentHasNoOtherActiveChildren := parentSummary.Counts != nil && parentSummary.Counts.Active > 1
	isParentInNeedOfDemotingToGroup := isChangeFromActive && parentSummary.State != nil && parentHasNoOtherActiveChildren
	//do we bubble this? I think we let the user make these changes.

	var incChange engine.AggregateCountChange
	if e.NewState != nil {
		incChange = engine.NewCountChange(*e.NewState, true)
	}

	decChange := engine.NewCountChange(e.OldState, false)

	combined := incChange.Combine(decChange)
	ancestorInc := combined
	parentInc := combined

	var activeState engine.GotState = engine.Active

	if isParentInNeedOfDemotingToGroup {
		decrementParentExistingState := engine.NewCountChange(*parentSummary.State, false)
		ancestorInc = ancestorInc.Combine(decrementParentExistingState)
		//rid the parent of its state
		parentSummary.State = nil
		upserts[parentSummaryId] = parentSummary
	} else if isParentInNeedOfPromotingToAcive {

		incrementActiveDueToParentUpgrade := engine.NewCountChange(activeState, true)
		ancestorInc = ancestorInc.Combine(incrementActiveDueToParentUpgrade)
		parentSummary.State = &activeState
		upserts[parentSummaryId] = parentSummary
	}

	//so at this point we MIGHT have an upsert for the parent already,
	//and we have increments established for the ancestors and for the parent.
	for _, summaryId := range e.Ancestry {
		if summaryId == engine.SummaryId(TheRootNoteInt32) {
			continue
		}
		if summaryId == parentSummaryId {
			existingUpsert, ok := upserts[parentSummaryId]
			if ok {
				existingUpsert.ApplyChange(parentInc)
				upserts[summaryId] = existingUpsert
			} else {
				parentSummary.ApplyChange(parentInc)
				upserts[summaryId] = parentSummary
			}
		} else {
			ancestorSummary, ok := ancestorAggs[summaryId]
			if !ok {
				return errors.New("missing agg")
			}
			ancestorSummary.ApplyChange(ancestorInc)
			upserts[summaryId] = ancestorSummary
		}
	}

	return a.summaryStore.UpsertManySummaries(upserts)
}

func (a *Aggregator) ItemDeleted(e ItemDeletedEvent) error {
	//we convert it to a statechanged event with no new state
	return a.ItemStateChanged(StateChangeEvent{
		Id:       e.Id,
		Ancestry: e.Ancestry,
		OldState: e.State,
		NewState: nil,
	})

}

func (a *Aggregator) ItemMoved(e ItemMovedEvent) error {
	// Fetch the moved item's summary to get its state
	movedItemSummary, err := a.summaryStore.Fetch([]engine.SummaryId{e.Id})
	if err != nil {
		return err
	}
	summary, ok := movedItemSummary[e.Id]
	if !ok {
		return errors.New("missing summary for moved item")
	}

	// If the item has no state (it's a group), we can't update counts based on it
	// Groups don't contribute to ancestor counts directly
	if summary.State == nil {
		return errors.New("cannot move a group node - only leaf nodes with state can be moved")
	}

	// Collect all ancestor IDs we need to fetch (excluding root)
	var allAncestorIds []engine.SummaryId
	for _, id := range e.OldAncestry {
		if id != engine.SummaryId(TheRootNoteInt32) {
			allAncestorIds = append(allAncestorIds, id)
		}
	}
	for _, id := range e.NewAncestry {
		if id != engine.SummaryId(TheRootNoteInt32) {
			allAncestorIds = append(allAncestorIds, id)
		}
	}

	// Fetch all ancestor summaries
	ancestorSummaries, err := a.summaryStore.Fetch(allAncestorIds)
	if err != nil {
		return err
	}

	upserts := make(map[engine.SummaryId]engine.Summary)

	// Decrement the moved item's state from all old ancestors
	decChange := engine.NewCountChange(*summary.State, false)
	for _, id := range e.OldAncestry {
		if id == engine.SummaryId(TheRootNoteInt32) {
			continue
		}
		ancestorSummary, ok := ancestorSummaries[id]
		if !ok {
			return fmt.Errorf("missing summary for old ancestor %d", id)
		}
		ancestorSummary.ApplyChange(decChange)
		upserts[id] = ancestorSummary
	}

	// Increment the moved item's state on all new ancestors
	incChange := engine.NewCountChange(*summary.State, true)
	for _, id := range e.NewAncestry {
		if id == engine.SummaryId(TheRootNoteInt32) {
			continue
		}
		// Check if we already have this ancestor in upserts (could be shared between old and new)
		if existingUpsert, ok := upserts[id]; ok {
			existingUpsert.ApplyChange(incChange)
			upserts[id] = existingUpsert
		} else {
			ancestorSummary, ok := ancestorSummaries[id]
			if !ok {
				return fmt.Errorf("missing summary for new ancestor %d", id)
			}
			ancestorSummary.ApplyChange(incChange)
			upserts[id] = ancestorSummary
		}
	}

	if len(upserts) == 0 {
		return nil
	}

	return a.summaryStore.UpsertManySummaries(upserts)
}

func (a *Aggregator) ItemEdited(e EditItemEvent) error {
	nowTime, err := engine.NewDateTime(time.Now())
	if err != nil {
		return err
	}

	ids := []engine.SummaryId{e.Id}
	list, err := a.summaryStore.Fetch(ids)
	if err != nil {
		return err
	}
	if len(list) != 1 {
		return errors.New("no item to update")
	}
	item := list[e.Id]
	//
	item.UpdatedDate = &nowTime

	upserts := make(map[engine.SummaryId]engine.Summary)
	upserts[e.Id] = item
	return a.summaryStore.UpsertManySummaries(upserts)
}
