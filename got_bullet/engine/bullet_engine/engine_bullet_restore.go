package bullet_engine

import (
	"vixac.com/got/engine/engine_util"
)

func (e *EngineBullet) CreateStoreFile() (string, error) {
	return engine_util.CreateStoreFile(e, e.LongFormStore)
}

func (e *EngineBullet) RestoreFromFile(filename string) error {
	return engine_util.RestoreFromFile(filename, e)
}
