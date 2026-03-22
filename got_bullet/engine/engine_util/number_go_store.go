package engine_util

import (
	"errors"
	"time"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/containers"
	"vixac.com/got/engine"
)

const numberGoKey = "ng"

// everything in one json body
type NumberGoBlock struct {
	Pairs map[int]string `json:"p"` //numberGo -> gid
}

type BulletNumberGoStore struct {
	Codec      Codec[NumberGoBlock]
	Collection bullet_stl.Collection
}

func NewBulletNumberGoStore(bucketId int32, track bullet_interface.TrackClientInterface, depot bullet_interface.DepotClientInterface, codec Codec[NumberGoBlock]) (engine.NumberGoStoreInterface, error) {
	coll := bullet_stl.NewBulletCollection(bucketId, track, depot)
	return &BulletNumberGoStore{
		Codec:      codec,
		Collection: coll,
	}, nil
}

func (n *BulletNumberGoStore) AssignNumberPairs(pairs []engine.NumberGoPair) error {
	pairMap := make(map[int]string)
	for _, p := range pairs {
		pairMap[p.Number] = p.Gid
	}
	payload, err := n.Codec.Encode(NumberGoBlock{Pairs: pairMap})
	if err != nil {
		return err
	}
	existing, err := n.Collection.AllItemsUnderPrefix(numberGoKey)
	if err != nil {
		return err
	}
	now := time.Now()
	if len(existing) == 0 {
		_, err = n.Collection.CreateItemUnder(numberGoKey, payload, &now)
		return err
	}
	var collId bullet_stl.CollectionId
	for k := range existing {
		collId = k
	}
	return n.Collection.EditPayload(collId, payload, &now)
}

func (n *BulletNumberGoStore) GidFor(number int) (*engine.GotId, error) {
	res, err := n.Collection.ItemsForKeys([]string{numberGoKey})
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	var payload string
	for _, v := range res {
		payload = v.Payload
	}
	var block NumberGoBlock
	err = n.Codec.Decode(payload, &block)
	if err != nil {
		return nil, err
	}
	value, ok := block.Pairs[number]
	if !ok {
		return nil, errors.New("missing number go id")
	}
	return engine.NewGotId(value)
}
