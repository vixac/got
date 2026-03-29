package grove_engine

import "vixac.com/got/engine/engine_util"

func (e *GroveEngine) CreateStoreFile() (string, error) {
	return engine_util.CreateStoreFile(e, e.LongFormStore)
}

func (e *GroveEngine) RestoreFromFile(filename string) error {
	return engine_util.RestoreFromFile(filename, e)
}
