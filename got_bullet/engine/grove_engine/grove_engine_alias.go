package grove_engine

import (
	"vixac.com/got/engine"
)

func (g *GroveEngine) LookupAliasForMany(gid []string) (map[string]*string, error) {
	return g.AliasStore.LookupAliasForMany(gid)
}
func (g *GroveEngine) Unalias(alias string) (*engine.GotId, error) {
	return g.AliasStore.Unalias(alias)
}
func (g *GroveEngine) Alias(lookup engine.GidLookup, alias string) (bool, error) {
	gid, err := g.GidLookup.InputToGid(&lookup)
	if err != nil || gid == nil {
		return false, err
	}
	return g.AliasStore.Alias(*gid, alias)
}
