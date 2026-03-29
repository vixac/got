package bullet_engine

import (
	"vixac.com/got/engine"
	"vixac.com/got/engine/engine_util"
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
	return engine_util.NotesForMany(e, lookup, e.LongFormStore)
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
