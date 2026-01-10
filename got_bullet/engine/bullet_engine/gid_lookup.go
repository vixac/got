package bullet_engine

import (
	"errors"
	"strconv"

	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/ids"
	"vixac.com/got/engine"
)

type GidLookupInterface interface {
	InputToGid(lookup *engine.GidLookup) (*engine.GotId, error)
}

// VX:TODO wants number<GO> lookup
type BulletGidLookup struct {
	AliasStore    engine.GotAliasInterface
	NumberGoStore NumberGoStoreInterface
	IdGenerator   IdGeneratorInterface
}

func NewBulletGidLookup(aliasStore engine.GotAliasInterface, numberGoStore NumberGoStoreInterface, idGen IdGeneratorInterface) (*BulletGidLookup, error) {
	return &BulletGidLookup{AliasStore: aliasStore, NumberGoStore: numberGoStore, IdGenerator: idGen}, nil
}

func (b *BulletGidLookup) InputToGid(lookup *engine.GidLookup) (*engine.GotId, error) {
	if lookup == nil || len(lookup.Input) == 0 {
		return engine.NewGotId(TheRootNode.Value)
	}
	/**
	The lookup can be one of the following:
	- exactly "0", this means its short hand for last. So we fetch the last created id.
	- A gid, we know this because it's prefixed with 0 <-- this is a harmless prefix to The Aasci -> Int ids algorithm
	- A number<GO> lookup from the last list printed. We know this because its prefixed with 1->9 (its a number < 0)
	- An alias. if it starts with an alphanumeric, its an alias
	*/

	//this is short hand for the last Id created
	if lookup.Input == "0" {
		lastId, err := b.IdGenerator.LastId()
		if err != nil {
			return nil, err
		}
		if lastId == 0 {
			return nil, errors.New("Invalid last id.")
		}
		// vx convert to string
		str, err := bullet_stl.BulletIdIntToaasci(lastId)
		if err != nil {
			return nil, err
		}
		return engine.NewGotId(str)
	}
	firstChar := lookup.Input[0]
	//this is a gid
	if firstChar == '0' {
		//we trim the first character and move on
		restOfString := lookup.Input[1:]
		return engine.NewGotId(restOfString)
	}

	//this is a number<GO> lookup
	if engine.CheckNumber([]byte(lookup.Input)) {
		number, err := strconv.Atoi(lookup.Input)
		if err != nil {
			return nil, err
		}
		return b.NumberGoStore.GidFor(number)

	}
	return b.AliasStore.Lookup(lookup.Input)

}
