package bullet_engine

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	"vixac.com/got/engine"
)

// The store that holds on to the meanings of the number goes, so when user
// can use them async
type NumberGoStoreInterface interface {
	AssignNumberPairs(pairs []NumberGoPair) error
	GidFor(number int) (*engine.GotId, error)
}

type NumberGoPair struct {
	Number string `json:"n"`
	Gid    string `json:"g"`
}

// everything in one json body
type NumberGoBlock struct {
	Pairs map[string]string `json:"p"` //numberGo -> gid
}

type BulletNumberGoStore struct {
	Codec   Codec[NumberGoBlock]
	DepotId int64
	Depot   bullet_interface.DepotClientInterface
}

func NewBulletNumberGoStore(client bullet_interface.DepotClientInterface, codec Codec[NumberGoBlock], depotId int64) (NumberGoStoreInterface, error) {
	return &BulletNumberGoStore{
		DepotId: depotId,
		Codec:   codec,
		Depot:   client,
	}, nil
}

func (n *BulletNumberGoStore) AssignNumberPairs(pairs []NumberGoPair) error {
	pairMap := make(map[string]string)

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
	fmt.Printf("VX: JSON IS '%s'\n", json)
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

	fmt.Printf("VX: JSON fetched is  IS '%s'\n", json)
	var block NumberGoBlock
	err = n.Codec.Decode(json, &block)
	if err != nil {
		fmt.Printf("VXL decode error %s", err.Error())
		return nil, err
	}
	numberToStr := strconv.Itoa(number)
	value, ok := block.Pairs[numberToStr]
	if !ok {
		return nil, errors.New("missing number go id")
	}
	fmt.Printf("VX: val is %s\n", value)
	return engine.NewGotId(value)
}
