package bullet_engine

import "vixac.com/got/engine"

func (e *EngineBullet) JotNote(lookup engine.GidLookup, note string) (engine.LongFormKey, error) {

	gid, err := e.GidLookup.InputToGid(&lookup)
	if err != nil || gid == nil {
		return engine.LongFormKey{}, err
	}

	id, err := e.LongFormStore.AppendNote(*gid, note)
	return *id, err
}

/*
func (e *EngineBullet) LongFormNotesFor(id engine.GotId) (*engine.LongFormBlockResult, error) {
	return e.LongFormStore.LongFormNotesFor(id)
}

func (e *EngineBullet) LongFormForMany(ids []engine.GotId) (map[engine.GotId]engine.LongFormBlockResult, error) {
	return e.LongFormStore.LongFormForMany(ids)
}

func (e *EngineBullet) RemoveAllItemsFromLongStoreUnder(id engine.GotId) error {
	return e.LongFormStore.RemoveAllItemsFromLongStoreUnder(id)
}
*/
