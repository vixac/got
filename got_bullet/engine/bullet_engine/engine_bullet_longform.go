package bullet_engine

import "vixac.com/got/engine"

/*
type LongFormStoreInterface interface {
	AppendNote(id engine.GotId, block LongFormBlock) error
	LongFormNotesFor(id engine.Gotid) (*LongFormBlockResult, error)
	LongFormForMany(ids []engine.GotId) (map[int32]LongFormBlockResult, error)
	RemoveAllItemsFromLongStoreUnder(id engine.GotId) error
}

*/
// This is the better pattern and we should stick to these.
var _ engine.LongFormStoreInterface = (*EngineBullet)(nil)

func (e *EngineBullet) AppendNote(id engine.GotId, block engine.LongFormBlock) error {
	return e.LongFormStore.AppendNote(id, block)
}
func (e *EngineBullet) LongFormNotesFor(id engine.GotId) (*engine.LongFormBlockResult, error) {
	return e.LongFormStore.LongFormNotesFor(id)
}

func (e *EngineBullet) LongFormForMany(ids []engine.GotId) (map[engine.GotId]engine.LongFormBlockResult, error) {
	return e.LongFormStore.LongFormForMany(ids)
}

func (e *EngineBullet) RemoveAllItemsFromLongStoreUnder(id engine.GotId) error {
	return e.LongFormStore.RemoveAllItemsFromLongStoreUnder(id)
}
