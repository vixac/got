package engine_util

import (
	"errors"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	"vixac.com/got/engine"
)

// everything in one json body
type NumberGoBlock struct {
	Pairs map[int]string `json:"p"` //numberGo -> gid
}

type BulletNumberGoStore struct {
	Codec   Codec[NumberGoBlock]
	DepotId int64
	Depot   bullet_interface.DepotClientInterface
}

func NewBulletNumberGoStore(client bullet_interface.DepotClientInterface, codec Codec[NumberGoBlock], depotId int64) (engine.NumberGoStoreInterface, error) {
	return &BulletNumberGoStore{
		DepotId: depotId,
		Codec:   codec,
		Depot:   client,
	}, nil
}

func (n *BulletNumberGoStore) AssignNumberPairs(pairs []engine.NumberGoPair) error {
	pairMap := make(map[int]string)

	for _, p := range pairs {

		pairMap[p.Number] = p.Gid
	}
	block := NumberGoBlock{
		Pairs: pairMap,
	}

	json, err := n.Codec.Encode(block)
	if err != nil {
		return err
	}
	req := bullet_interface.DepotRequest{
		Key:   n.DepotId,
		Value: json,
	}
	return n.Depot.DepotUpsertMany([]bullet_interface.DepotRequest{req})
}

func (n *BulletNumberGoStore) GidFor(number int) (*engine.GotId, error) {

	keys := []int64{n.DepotId}
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

	json, ok := res.Values[n.DepotId]
	if !ok {
		return nil, nil
	}

	var block NumberGoBlock
	err = n.Codec.Decode(json, &block)
	if err != nil {
		return nil, err
	}
	value, ok := block.Pairs[number]
	if !ok {
		return nil, errors.New("missing number go id")
	}
	return engine.NewGotId(value)
}
