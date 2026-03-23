package bullet_engine

import (
	"sort"

	"vixac.com/got/engine"
)

func (e *EngineBullet) JotNote(lookup engine.GidLookup, note string) (engine.LongFormKey, error) {

	gid, err := e.GidLookup.InputToGid(&lookup)
	if err != nil || gid == nil {
		return engine.LongFormKey{}, err
	}

	id, err := e.LongFormStore.AppendNote(*gid, note)
	return *id, err
}

func (e *EngineBullet) NotesForMany(lookup *engine.GidLookup) (*engine.LongFormBlockResult, error) {
	items, err := e.FetchItemsBelow(lookup, false, []engine.GotState{engine.Active, engine.Complete}, false)
	if err != nil {
		return nil, err
	}

	var gotIds []engine.GotId
	for _, v := range items.Result {
		gotIds = append(gotIds, v.GotId)
	}
	mapOfResults, err := e.LongFormStore.LongFormForMany(gotIds)
	if err != nil {
		return nil, err
	}
	var blocks []engine.LongFormBlock
	for _, v := range mapOfResults {
		blocks = append(blocks, v.Blocks...)
	}

	//sorted with latest first
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].Edited.Before(blocks[j].Edited)
	})
	res := engine.LongFormBlockResult{
		Blocks: blocks,
	}

	return &res, nil
}

func (e *EngineBullet) NotesFor(lookup *engine.GidLookup, recurse bool) (*engine.LongFormBlockResult, error) {
	if lookup == nil || recurse {
		return e.NotesForMany(lookup)
	}
	gid, err := e.GidLookup.InputToGid(lookup)
	if err != nil {
		return nil, err
	}

	return e.LongFormStore.LongFormNotesFor(*gid)
}
