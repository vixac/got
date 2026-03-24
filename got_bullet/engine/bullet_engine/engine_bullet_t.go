package bullet_engine

import (
	"errors"
	"fmt"

	"vixac.com/got/engine"
	"vixac.com/got/engine/engine_util"
)

func (e *EngineBullet) OpenThenTimestamp(lookup engine.GidLookup) error {

	gid, err := e.GidLookup.InputToGid(&lookup)
	if err != nil || gid == nil {
		return err
	}
	summaryId := engine.SummaryId(gid.IntValue)
	exists, err := e.SummaryStore.Fetch([]engine.SummaryId{summaryId})
	if err != nil {
		return err
	}
	if exists != nil {
		_, ok := exists[summaryId]
		if !ok {
			return errors.New("This gid does not exist.")
		}
	}

	existing, err := e.LongFormStore.LongFormNotesFor(*gid)
	if err != nil {
		return err
	}

	commentedOutNotes := engine_util.ConsolidateBlocksIntoCommentedString(existing)

	withRemovedComments, err := engine_util.OpenTextEditorWithCommentedOutString(commentedOutNotes)
	if err != nil {
		return err
	}
	if withRemovedComments == nil {
		fmt.Printf("VX: No comments made.")
		return nil
	}

	newId, err := e.LongFormStore.AppendNote(*gid, *withRemovedComments)
	if err != nil {
		return err
	}
	fmt.Printf("VX: New note: %s\n", newId.ToString())

	//we send the edit event so the update time gets changed
	e.publishEditEvent(engine.EditItemEvent{Id: engine.SummaryId(gid.IntValue)})
	return nil
}
