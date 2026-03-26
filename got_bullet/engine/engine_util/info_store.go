package engine_util

import (
	"errors"
	"time"

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

func NewBuckInfo(title string, deadline *engine.DateTime, created engine.DateTime, updated engine.DateTime, tags []engine.Tag, flags map[string]bool) BuckInfo {
	return BuckInfo{
		Title:       title,
		Deadline:    deadline,
		CreatedDate: created,
		UpdatedDate: updated,
		Tags:        tags,
		Flags:       flags,
	}
}

type BuckStore struct {
	Codec      Codec[BuckInfo]
	Collection bullet_stl.Collection
}

func NewBuckStore(bucketId int32, track bullet_interface.TrackClientInterface, depot bullet_interface.DepotClientInterface, codec Codec[BuckInfo]) BuckStoreInterface {
	coll := bullet_stl.NewBulletCollection(bucketId, track, depot)
	return &BuckStore{
		Codec:      codec,
		Collection: coll,
	}
}

func (b *BuckStore) UpsertInfo(id engine.GotId, info BuckInfo) error {
	key := idToStr(id)
	existing, err := b.Collection.AllItemsUnderPrefix(key)
	if err != nil {
		return err
	}

	encoded, err := b.Codec.Encode(info)
	if err != nil {
		return err
	}

	now := time.Now()
	if len(existing) == 0 {
		_, err := b.Collection.CreateItemUnder(key, encoded, &now)
		return err
	}
	if len(existing) != 1 {
		return errors.New("upserting to a key that is not unique")
	}

	var theCollId bullet_stl.CollectionId
	for k := range existing {
		theCollId = k
	}
	return b.Collection.EditPayload(theCollId, encoded, &now)
}

func (b *BuckStore) DeleteInfoMany(ids []engine.GotId) error {
	var keys []string
	for _, id := range ids {
		keys = append(keys, idToStr(id))
	}
	res, err := b.Collection.ItemsForKeys(keys)
	if err != nil || res == nil {
		return err
	}

	collids := make([]bullet_stl.CollectionId, 0, len(res))
	for k := range res {
		collids = append(collids, k)
	}
	return b.Collection.DeleteItems(collids)
}

func (b *BuckStore) InfoForMany(ids []engine.GotId) (*InfoManyResponse, error) {
	keyToId := make(map[string]engine.GotId)
	var keys []string
	for _, id := range ids {
		key := idToStr(id)
		keys = append(keys, key)
		keyToId[key] = id
	}

	resp, err := b.Collection.ItemsForKeys(keys)
	if err != nil {
		return nil, err
	}

	infoMap := make(map[engine.GotId]BuckInfo)
	foundKeys := make(map[string]bool)

	for collId, item := range resp {
		var info BuckInfo
		if err := b.Codec.Decode(item.Payload, &info); err != nil {
			return nil, err
		}
		gotId := keyToId[collId.Key]
		infoMap[gotId] = info
		foundKeys[collId.Key] = true
	}

	var missing []engine.GotId
	for key, id := range keyToId {
		if !foundKeys[key] {
			missing = append(missing, id)
		}
	}

	return &InfoManyResponse{
		InfoMap: infoMap,
		Missing: missing,
	}, nil
}
