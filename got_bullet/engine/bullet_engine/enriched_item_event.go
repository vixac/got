package bullet_engine

import (
	"errors"
	"fmt"

	"vixac.com/got/engine"
)

/*
*
This is an item Event that has its ancestry fetched, and it means it knows what it is, and what it's parent is.
*/

// VX:TODO reuse for all events?
type EnrichedAddItemEvent struct {
	Event    AddItemEvent
	Ancestry []EnrichedSummary
}
type EnrichedSummary struct {
	Id      SummaryId
	Summary Summary
}

func NewEnrichedAddItemEvent(event AddItemEvent, summaries map[SummaryId]Summary) (EnrichedAddItemEvent, error) {

	var ancestry []EnrichedSummary

	for _, summaryId := range event.Ancestry {
		if summaryId == SummaryId(TheRootNoteInt32) {
			continue
		}
		summary, ok := summaries[summaryId]
		if !ok {
			fmt.Printf("VX: Error fetching for id %d\n", summaryId)
			return EnrichedAddItemEvent{}, errors.New("missing summary when creating enriched item")
		}
		ancestry = append(ancestry, EnrichedSummary{
			Id:      summaryId,
			Summary: summary,
		})
	}
	return EnrichedAddItemEvent{
		Event:    event,
		Ancestry: ancestry,
	}, nil
}

func (e EnrichedAddItemEvent) Parent() *EnrichedSummary {
	if len(e.Ancestry) == 0 {
		return nil //parent is root
	}
	return &e.Ancestry[len(e.Ancestry)-1]
}

func (e EnrichedAddItemEvent) ParentIsLeaf() bool {
	parent := e.Parent()
	if parent == nil {
		return false
	}
	return parent.Summary.IsLeaf()

}
func (e EnrichedAddItemEvent) ParentId() *SummaryId {
	parent := e.Parent()
	if parent == nil {
		return nil
	}
	return &parent.Id
}
func (e EnrichedAddItemEvent) ParentIsRoot() bool {
	return e.Parent() == nil
}

func (e EnrichedAddItemEvent) ParentIsGroup() bool {
	parent := e.Parent()
	if parent == nil {
		return false //root is a group? but not really
	}
	return parent.Summary.Counts != nil
}

func (e EnrichedAddItemEvent) ParentState() *engine.GotState {
	parent := e.Parent()
	if parent == nil {
		return nil
	}
	return parent.Summary.State
}
