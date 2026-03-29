package grove_engine

import (
	"errors"

	"vixac.com/got/engine"
)

func (g *GroveEngine) DeleteMany(lookups []engine.GidLookup) error {
	return errors.New("Delete many not impl")
}
func (g *GroveEngine) Move(lookup engine.GidLookup, newParent engine.GidLookup) error {
	return errors.New("Move  many not impl")
}
