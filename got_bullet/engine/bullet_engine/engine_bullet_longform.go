package bullet_engine

// This is the better pattern and we should stick to these.
var _ LongFormStoreInterface = (*EngineBullet)(nil)

func (e *EngineBullet) UpsertItem(id int32, title string) error {
	return e.LongFormStore.UpsertItem(id, title)
}
func (e *EngineBullet) LongFormFor(id int32) (*string, error) {
	return e.LongFormStore.LongFormFor(id)
}

func (e *EngineBullet) LongFormForMany(ids []int32) (map[int32]string, error) {
	return e.LongFormStore.LongFormForMany(ids)
}

func (e *EngineBullet) RemoveItem(id int32) error {
	return e.LongFormStore.RemoveItem(id)
}
