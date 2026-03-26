package engine_util

import (
	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/containers"
	"vixac.com/got/engine"
)

const scheduleKey = "sk"

type BuckStoreInterface interface {
	//insert or replace the info for this got id
	UpsertInfo(id engine.GotId, info BuckInfo) error
	//delete every info that exists for this list of ids
	DeleteInfoMany(ids []engine.GotId) error

	//fetch all info for these ids, and list any missing info in the Missing field
	InfoForMany(ids []engine.GotId) (*InfoManyResponse, error)
}

type InfoManyResponse struct {
	InfoMap map[engine.GotId]BuckInfo
	Missing []engine.GotId
}

/*
All the larger content that one buck might have.
*/
type BuckInfo struct {
	Title    string           `json:"t,omitempty"`
	Deadline *engine.DateTime `json:"d,omitempty"` //VX:TODO consider a separate store of deadlines to ids
	//Alias  we have aliasstore for now, but we could move the gid -> alias mapping to here.
	CreatedDate engine.DateTime `json:"c,omitempty"`
	UpdatedDate engine.DateTime `json:"u,omitempty"`
	Tags        []engine.Tag    `json:"tags,omitempty"`
	Flags       map[string]bool `json:"f,omitempty"`
}

type BuckStore struct {
	Codec      Codec[BuckInfo]
	Collection bullet_stl.Collection
}

func NewBuckStore(bucketId int32, track bullet_interface.TrackClientInterface, depot bullet_interface.DepotClientInterface, codec Codec[BuckInfo]) (BuckStoreInterface, error) {
	coll := bullet_stl.NewBulletCollection(bucketId, track, depot)
	return &BuckStore{
		Codec:      codec,
		Collection: coll,
	}, nil
}

func (b *BuckStore) UpsertInfo(id engine.GotId, info BuckInfo) error {
	return nil

}
func (b *BuckStore) DeleteInfoMany(ids []engine.GotId) error {
	return nil

}
func (b *BuckStore) InfoForMany(ids []engine.GotId) (*InfoManyResponse, error) {
	return nil, nil
}
