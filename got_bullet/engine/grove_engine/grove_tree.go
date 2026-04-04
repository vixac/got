package grove_engine

import (
	"errors"

	"vixac.com/got/engine"
)

func (g *GroveEngine) DeleteMany(lookups []engine.GidLookup) error {
	return errors.New("Delete many not impl")
}

func (g *GroveEngine) Move(lookup engine.GidLookup, newParent engine.GidLookup) error {
	gid, err := g.GidLookup.InputToGid(&lookup)
	if err != nil || gid == nil {
		return err
	}
	parent, err := g.GidLookup.InputToGid(&newParent)
	if err != nil {
		return err
	}
	return g.GroveStore.Move(*gid, *parent)
}
