package bullet_engine

import (
	"errors"

	"vixac.com/got/engine"
)

func (e *EngineBullet) Delete(lookup engine.GidLookup) error {
	// Convert lookup to GID
	gid, err := e.GidLookup.InputToGid(&lookup)
	if err != nil {
		return err
	}
	if gid == nil {
		return errors.New("could not resolve gid from lookup")
	}

	// Check if this item is a parent (has children)
	children, err := e.AncestorList.FetchImmediatelyUnder(*gid)
	if err != nil {
		return err
	}
	if children != nil && len(children.Ids) > 0 {
		return errors.New("cannot delete item: it has children")
	}

	summaryId := engine.SummaryId(gid.IntValue)

	// Fetch state and ancestry BEFORE deletion for the delete event
	summary, err := e.SummaryStore.Fetch([]engine.SummaryId{summaryId})
	if err != nil {
		return err
	}
	itemSummary, ok := summary[summaryId]
	if !ok {
		return errors.New("item not found in summary store")
	}

	// Get the state (default to Note if nil)
	var itemState engine.GotState = engine.Note
	if itemSummary.State != nil {
		itemState = *itemSummary.State
	}

	// Fetch ancestry
	ancestorResult, err := e.AncestorList.FetchAncestorsOf(*gid)
	if err != nil {
		return err
	}
	var ancestryIds []engine.SummaryId
	if ancestorResult != nil {
		for _, ancestor := range ancestorResult.Ids {
			ancestryIds = append(ancestryIds, engine.SummaryId(ancestor.IntValue))
		}
	}

	// Delete alias if it exists
	alias, err := e.LookupAliasForGid(gid.AasciValue)
	if err != nil {
		return err
	}
	if alias != nil {
		_, err = e.AliasStore.Unalias(*alias)
		if err != nil {
			return err
		}
	}

	// Delete longForm entry if it exists
	err = e.LongFormStore.RemoveItemFromLongStore(gid.IntValue)
	if err != nil {
		return err
	}

	// Delete summary
	err = e.SummaryStore.Delete([]engine.SummaryId{summaryId})
	if err != nil {
		return err
	}

	// Delete title
	err = e.TitleStore.RemoveItem(gid.IntValue)
	if err != nil {
		return err
	}

	// Delete from ancestor list
	err = e.AncestorList.RemoveItem(*gid)
	if err != nil {
		return err
	}

	// Publish delete event with state and ancestry
	return e.publishItemDeletedEvent(ItemDeletedEvent{
		Id:       summaryId,
		State:    itemState,
		Ancestry: ancestryIds,
	})
}
