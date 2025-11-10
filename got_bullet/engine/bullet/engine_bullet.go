package bullet

import (
	"errors"

	wayfinder "github.com/vixac/firbolg_clients/bullet/wayfinder"
	"vixac.com/got/engine"
)

const (
	aliasBucket int32 = 1001
	nodeBucket  int32 = 1002
)

type EngineBullet struct {
	WayFinder wayfinder.WayFinderClientInterface
}

func (e *EngineBullet) Summary(lookup engine.GidLookup) (*engine.GotSummary, error) {

	query := wayfinder.WayFinderPrefixQueryRequest{
		BucketId: nodeBucket,
	}
	res, err := e.WayFinder.WayFinderQueryByPrefix(query)
	if err != nil {
		return nil, err
	}
	println("VX: res is", res)
	return nil, errors.New("not implemeneted")

}

func (e *EngineBullet) Alias(gid string, alias string) (bool, error) {
	return false, errors.New(("not impl"))
}

func (e *EngineBullet) Resolve(lookup engine.GidLookup) (*engine.NodeId, error) {
	//check if the gid is an exact match for an item id
	//check int32 parse, check its length is the right length

	//aliases can't start with a number.
	return nil, errors.New("not impl")
}

func (e *EngineBullet) Delete(lookup engine.GidLookup) (*engine.NodeId, error) {
	//check if the gid is an exact match for an item id
	//check int32 parse, check its length is the right length

	//aliases can't start with a number.
	return nil, errors.New("not impl")
}

func (e *EngineBullet) Unalias(alias string) (*engine.NodeId, error) {
	//check if the gid is an exact match for an item id
	//check int32 parse, check its length is the right length

	//aliases can't start with a number.
	return nil, errors.New("not impl")
}

func (e *EngineBullet) Move(lookup engine.GidLookup, newParent engine.GidLookup) (*engine.NodeId, error) {
	return nil, errors.New("not impl")
}

/**
The cache is going to work like this

asc(n)
and its only ascendants and descen

*/
/**
ok im goina do this the buck way..
it takes a wayfinder and uses it
and the business logic for using the wayfinder is reusable.
its a cool design.
*/
