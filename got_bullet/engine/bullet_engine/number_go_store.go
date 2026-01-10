package bullet_engine

import (
	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/ids"
	"vixac.com/got/engine"
)

// The store that holds on to the meanings of the number goes, so when user
// can use them async
type NumberGoStoreInterface interface {
	AssignNumberPairs(pairs []NumberGoPair) error
	GidFor(number int) (*engine.GotId, error)
}

type NumberGoPair struct {
	Number int
	Gid    engine.Gid
}

type BulletNumberGoStore struct {
	Namespace int32
	Depot     bullet_interface.DepotClientInterface
}

func NewBulletNumberGoStore(client bullet_interface.DepotClientInterface, namespaceId int32) (NumberGoStoreInterface, error) {
	return &BulletNumberGoStore{
		Namespace: namespaceId,
		Depot:     client,
	}, nil
}

func (n *BulletNumberGoStore) AssignNumberPairs(pairs []NumberGoPair) error {

	var reqs []bullet_interface.DepotRequest
	for _, p := range pairs {
		namespacedId := bullet_stl.MakeNamespacedId(n.Namespace, int32(p.Number))
		reqs = append(reqs, bullet_interface.DepotRequest{
			Key:   namespacedId,
			Value: p.Gid.Id,
		})
	}
	return n.Depot.DepotUpsertMany(reqs)
}

func (n *BulletNumberGoStore) GidFor(number int) (*engine.GotId, error) {
	namespacedId := bullet_stl.MakeNamespacedId(n.Namespace, int32(number))

	keys := []int64{namespacedId}
	manyReq := bullet_interface.DepotGetManyRequest{
		Keys: keys,
	}
	res, err := n.Depot.DepotGetMany(manyReq)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}

	value, ok := res.Values[namespacedId]
	if !ok {
		return nil, nil
	}

	return engine.NewGotId(value)
}
