package bullet_engine

import (
	"errors"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/ids"

	"vixac.com/got/engine"
)

const (
	aliasBucket int32 = 1001
	nodeBucket  int32 = 1002
)

type EngineBullet struct {
	Client bullet_interface.BulletClientInterface
}

func (e *EngineBullet) Summary(lookup *engine.GidLookup) (*engine.GotSummary, error) {

	query := bullet_interface.TrackGetItemsByPrefixRequest{
		BucketID: nodeBucket,
		Prefix:   ":",
	}
	res, err := e.Client.TrackGetManyByPrefix(query)
	if err != nil {
		return nil, err
	}

	var foundId *int64 = nil
	for bucket, values := range res.Values {
		if bucket != nodeBucket {
			continue
		}
		for _, v := range values {
			foundId = &v.Value
		}

	}
	if foundId != nil {
		var keys []int64
		keys = append(keys, *foundId)
		manyReq := bullet_interface.DepotGetManyRequest{
			Keys: keys,
		}
		e.Client.DepotGetMany(manyReq)
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

func (e *EngineBullet) CreateBuck(parent *engine.GidLookup, date *engine.DateLookup, completable bool, heading string) (*engine.NodeId, error) {
	//VX:TODO this should hit both the keys and also hit depot too for the heading.

	newId, err := e.NextId()
	if err != nil {
		return nil, err
	}
	stringId, err := bullet_stl.BulletIdIntToaasci(newId)
	if err != nil {
		return nil, err
	}
	//VX:TODO thats not quite right but lets go with it.
	err = e.Client.TrackInsertOne(nodeBucket, ":"+stringId, newId, nil, nil)

	if err != nil {
		return nil, err
	}
	depotReq := bullet_interface.DepotRequest{
		Key:   newId,
		Value: heading,
	}
	e.Client.DepotInsertOne(depotReq)

	//lets
	return nil, err
}

func (e *EngineBullet) Lookup(alias string) (*engine.NodeId, error) {
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
