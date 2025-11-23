package bullet_engine

import (
	"strconv"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/containers"
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
	OneWay bullet_stl.OneWayList
}

func NewBulletNumberGoStore(client bullet_interface.TrackClientInterface, bucketId int32) (NumberGoStoreInterface, error) {
	oneWay, err := bullet_stl.NewBulletOneWayList(client, bucketId, "numbergoes", ">")
	if err != nil {
		return nil, err
	}
	return &BulletNumberGoStore{
		OneWay: oneWay,
	}, nil
}

func (n *BulletNumberGoStore) AssignNumberPairs(pairs []NumberGoPair) error {
	//VX:TODO this should be a bulk insert
	for _, pair := range pairs {
		numberStr := strconv.Itoa(pair.Number)
		err := n.OneWay.Upsert(bullet_stl.ListSubject{numberStr}, bullet_stl.ListObject{pair.Gid.Id})
		if err != nil {
			return err
		}
	}
	return nil
}
func (n *BulletNumberGoStore) GidFor(number int) (*engine.GotId, error) {
	numberStr := strconv.Itoa(number)
	object, err := n.OneWay.GetObject(bullet_stl.ListSubject{Value: numberStr})
	if err != nil {
		return nil, err
	}
	if object == nil {
		return nil, nil
	}
	gid, err := engine.NewGotId(object.Value)
	if err != nil {
		return nil, err
	}
	return gid, nil
}
