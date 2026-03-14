package bullet_engine

import "vixac.com/got/engine"

// This is the better pattern and we should stick to these.
var _ engine.LongFormStoreInterface = (*EngineBullet)(nil)

func (e *EngineBullet) UpsertItem(id int32, block engine.LongFormBlock) error {
	return e.LongFormStore.UpsertItem(id, block)
}
func (e *EngineBullet) LongFormFor(id int32) (*engine.LongFormBlockResult, error) {
	return e.LongFormStore.LongFormFor(id)
}

func (e *EngineBullet) LongFormForMany(ids []int32) (map[int32]engine.LongFormBlockResult, error) {
	return e.LongFormStore.LongFormForMany(ids)
}

func (e *EngineBullet) RemoveAllItemsFromLongStore(id int32) error {
	return e.LongFormStore.RemoveAllItemsFromLongStore(id)
}
