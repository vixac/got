package bullet_engine

import (
	"errors"

	"vixac.com/got/engine"
)

func (e *EngineBullet) Collpase(lookup engine.GidLookup) error {
	gid, err := e.GidLookup.InputToGid(&lookup)
	if err != nil || gid == nil {
		return err
	}
	//VX:TODO so this needs to add the item to a set of all collapsed items
	//which get chceked.
	return errors.New("Not impl.")

}
