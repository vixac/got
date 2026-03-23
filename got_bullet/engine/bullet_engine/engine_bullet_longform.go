package bullet_engine

import (
	"errors"
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

func (e *EngineBullet) NotesFor(lookup *engine.GidLookup, recurse bool) (*engine.LongFormBlockResult, error) {
	if lookup == nil {
		fmt.Printf("VX: TODO HANDLE NIL LOOKUP")
		return nil, errors.New("VX:TODO NIL LOOKUP")
	}
	fmt.Printf("VX: Lookup is %s\n", lookup.Input)
	//VX:TODO handle nil lookup
	gid, err := e.GidLookup.InputToGid(lookup)
	if err != nil {
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
