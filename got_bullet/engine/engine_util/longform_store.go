package engine_util

import (
	"errors"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/containers"
	"vixac.com/got/engine"
)

type BulletLongFormStore struct {
	Collection bullet_stl.Collection
}

func NewBulletLongFormStore(bucketId int32, track bullet_interface.TrackClientInterface, depot bullet_interface.DepotClientInterface) (engine.LongFormStoreInterface, error) {
	coll := bullet_stl.NewBulletCollection(bucketId, track, depot)
	return &BulletLongFormStore{Collection: coll}, nil
}

func (s *BulletLongFormStore) UpsertItem(id int32, block engine.LongFormBlock) error {
	idStr := idToStr(id)
	existing, err := s.Collection.AllItemsUnderPrefix(idStr)
	if err != nil {
		return err
	}
	if len(existing) == 0 {
		_, err := s.Collection.CreateItemUnder(idStr, block.Content)
		return err
	}
	if len(existing) != 1 {
		return errors.New("upserting to a key that is not unique")
	}
	var theCollId bullet_stl.CollectionId
	for k := range existing {
		theCollId = k
	}
	return s.Collection.EditPayload(theCollId, block.Content)
}

func (s *BulletLongFormStore) LongFormForMany(ids []int32) (map[int32]engine.LongFormBlockResult, error) {
	var idStrings []string
	for _, id := range ids {
		idStrings = append(idStrings, idToStr(id))
	}
	resp, err := s.Collection.ItemsForKeys(idStrings)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	result := make(map[int32]engine.LongFormBlockResult)
	for k, v := range resp {
		id, err := strToId(k.Key)
		if err != nil {
			return nil, err
		}
		block := engine.LongFormBlock{Content: v}
		result[id] = engine.LongFormBlockResult{Blocks: []engine.LongFormBlock{block}}
	}
	return result, nil
}

func (s *BulletLongFormStore) LongFormFor(id int32) (*engine.LongFormBlockResult, error) {
	res, err := s.LongFormForMany([]int32{id})
	if err != nil || len(res) == 0 {
		return nil, err
	}
	r := res[id]
	return &r, nil
}

func (s *BulletLongFormStore) RemoveAllItemsFromLongStore(id int32) error {
	res, err := s.Collection.ItemsForKeys([]string{idToStr(id)})
	if err != nil || res == nil {
		return err
	}
	var collIds []bullet_stl.CollectionId
	for k := range res {
		collIds = append(collIds, k)
	}
	return s.Collection.DeleteItems(collIds)
}
