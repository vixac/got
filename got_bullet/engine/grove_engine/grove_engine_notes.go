package grove_engine

import (
	"errors"
	"fmt"

	"vixac.com/got/engine"
	"vixac.com/got/engine/engine_util"
)

func (e *GroveEngine) JotNote(lookup engine.GidLookup, note string) (engine.LongFormKey, error) {

	gid, err := e.GidLookup.InputToGid(&lookup)
	if err != nil || gid == nil {
		return engine.LongFormKey{}, err
	}

	id, err := e.LongFormStore.AppendNote(*gid, note)
	return *id, err
}

func (e *GroveEngine) NotesForMany(lookup *engine.GidLookup) (*engine.LongFormBlockResult, error) {
	return engine_util.NotesForMany(e, lookup, e.LongFormStore)
}

func (e *GroveEngine) NotesFor(lookup *engine.GidLookup, recurse bool) (*engine.LongFormBlockResult, error) {
	if lookup == nil || recurse {
		return e.NotesForMany(lookup)
	}
	gid, err := e.GidLookup.InputToGid(lookup)
	if err != nil {
		return nil, err
	}

	return e.LongFormStore.LongFormNotesFor(*gid)
}

func (e *GroveEngine) OpenThenTimestamp(lookup engine.GidLookup) error {

	gid, err := e.GidLookup.InputToGid(&lookup)
	if err != nil || gid == nil {
		return err
	}
	info, err := e.InfoStore.InfoForMany([]engine.GotId{*gid})
	if err != nil {
		return err
	}
	if len(info.InfoMap) == 0 {
		return errors.New("This gid does not exist.")
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
	//e.publishEditEvent(engine.EditItemEvent{Id: engine.SummaryId(gid.IntValue)})
	return nil
}
