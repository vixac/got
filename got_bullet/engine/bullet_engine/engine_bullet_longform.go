package bullet_engine

import (
	"fmt"

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

func (e *EngineBullet) NotesFor(lookup engine.GidLookup) (*engine.LongFormBlockResult, error) {
	fmt.Printf("VX: Lookup is %s\n", lookup.Input)
	gid, err := e.GidLookup.InputToGid(&lookup)
	if err != nil || gid == nil {
		return nil, err
	}
	return e.LongFormStore.LongFormNotesFor(*gid)
}

/*


func (e *EngineBullet) LongFormForMany(ids []engine.GotId) (map[engine.GotId]engine.LongFormBlockResult, error) {
	return e.LongFormStore.LongFormForMany(ids)
}

func (e *EngineBullet) RemoveAllItemsFromLongStoreUnder(id engine.GotId) error {
	return e.LongFormStore.RemoveAllItemsFromLongStoreUnder(id)
}
*/
