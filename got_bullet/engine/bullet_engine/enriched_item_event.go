package bullet_engine

import (
	"errors"

	"vixac.com/got/engine"
)

/*
*
This is an item Event that has its ancestry fetched, and it means it knows what it is, and what it's parent is.
*/

type EnrichedAncestry struct {
	Ancestry []EnrichedSummary
}

type EnrichedSummary struct {
	Id      engine.SummaryId
	Summary engine.Summary
}

func EnrichSummaries(ancestry []engine.SummaryId, summaries map[engine.SummaryId]engine.Summary) (EnrichedAncestry, error) {
	var enriched []EnrichedSummary

	for _, summaryId := range ancestry {
		if summaryId == engine.SummaryId(TheRootNoteInt32) {
			continue
		}
		summary, ok := summaries[summaryId]
		if !ok {
			return EnrichedAncestry{}, errors.New("missing summary when creating enriched item")
		}
		enriched = append(enriched, EnrichedSummary{
			Id:      summaryId,
			Summary: summary,
		})
	}

	return EnrichedAncestry{
		Ancestry: enriched,
	}, nil
}

func (e EnrichedAncestry) Parent() *EnrichedSummary {
	if len(e.Ancestry) == 0 {
		return nil //parent is root
	}
	return &e.Ancestry[len(e.Ancestry)-1]
}

func (e EnrichedAncestry) ParentIsLeaf() bool {
	parent := e.Parent()
	if parent == nil {
		return false
	}
	return parent.Summary.IsLeaf()

}
func (e EnrichedAncestry) ParentId() *engine.SummaryId {
	parent := e.Parent()
	if parent == nil {
		return nil
	}
	return &parent.Id
}
func (e EnrichedAncestry) ParentIsRoot() bool {
	return e.Parent() == nil
}

func (e EnrichedAncestry) ParentIsGroup() bool {
	parent := e.Parent()
	if parent == nil {
		return false //root is a group? but not really
	}
	return parent.Summary.Counts != nil
}

func (e EnrichedAncestry) ParentState() *engine.GotState {
	parent := e.Parent()
	if parent == nil {
		return nil
	}
	return parent.Summary.State
}
