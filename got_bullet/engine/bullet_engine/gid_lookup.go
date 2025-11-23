package bullet_engine

import (
	"strconv"
	"unicode"

	"vixac.com/got/engine"
)

type GidLookupInterface interface {
	InputToGid(lookup *engine.GidLookup) (*engine.GotId, error)
}

// VX:TODO wants number<GO> lookup
type BulletGidLookup struct {
	AliasStore    engine.GotAliasInterface
	NumberGoStore NumberGoStoreInterface
}

func NewBulletGidLookup(aliasStore engine.GotAliasInterface, numberGoStore NumberGoStoreInterface) (*BulletGidLookup, error) {
	return &BulletGidLookup{AliasStore: aliasStore, NumberGoStore: numberGoStore}, nil
}

func CheckNumber(p []byte) bool {
	r := string(p)
	sep := 0
	for _, b := range r {
		if unicode.IsNumber(b) {
			continue
		}
		if b == rune('.') {
			if sep > 0 {
				return false
			}
			sep++
			continue
		}
		return false
	}
	return true
}
func (b *BulletGidLookup) InputToGid(lookup *engine.GidLookup) (*engine.GotId, error) {
	if lookup == nil || len(lookup.Input) == 0 {
		return engine.NewGotId(TheRootNode.Value)
	}
	/**
	The lookup can be one of the following:
	- A gid, we know this because it's prefixed with 0 <-- this is a harmless prefix to The Aasci -> Int ids algorithm
	- A number<GO> lookup from the last list printed. We know this because its prefixed with 1->9 (its a number < 0)
	- An alias. if it starts with an alphanumeric, its an alias
	*/
	firstChar := lookup.Input[0]
	//this is a gid
	if firstChar == '0' {
		//we trim the first character and move on
		restOfString := lookup.Input[1:]
		return engine.NewGotId(restOfString)
	}

	//this is a number<GO> lookup
	if CheckNumber([]byte(lookup.Input)) {
		number, err := strconv.Atoi(lookup.Input)
		if err != nil {
			return nil, err
		}
		return b.NumberGoStore.GidFor(number)

	}
	return b.AliasStore.Lookup(lookup.Input)

}
