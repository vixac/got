package engine_util

import (
	"errors"
	"strconv"
	"time"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/containers"
	"vixac.com/got/engine"
)

type TitleStoreInterface interface {
	UpsertItem(id int32, title string) error
	TitleFor(id int32) (*string, error)

	TitleForMany(ids []int32) (map[int32]string, error)
	RemoveItem(id int32) error
}

// VX:TODO move to bullet_engine?
type BulletTitleStore struct {
	Collection bullet_stl.Collection
}

func NewBulletTitleStore(bucketId int32, track bullet_interface.TrackClientInterface, depot bullet_interface.DepotClientInterface) TitleStoreInterface {
	coll := bullet_stl.NewBulletCollection(bucketId, track, depot)
	return &BulletTitleStore{
		Collection: coll,
	}
}

func idToStr(id engine.GotId) string { //I need to
	return id.AasciValue //strconv.Itoa(int(id.IntValue))
}

func strToId(key string) (int32, error) {
	id, err := strconv.Atoi(key)
	return int32(id), err
}
func (s *BulletTitleStore) UpsertItem(id int32, title string) error {

	//first we check if it exists:
	idStr := strconv.Itoa(int(id))
	existing, err := s.Collection.AllItemsUnderPrefix(idStr)
	if err != nil {
		return err
	}

	//create the item
	now := time.Now()
	if existing == nil || len(existing) == 0 {
		_, err := s.Collection.CreateItemUnder(idStr, title, &now)
		return err
	}
	if len(existing) != 1 {
		return errors.New("Upserting to a key that is not unique.")

	}

	var theCollId bullet_stl.CollectionId
	for k, _ := range existing {
		theCollId = k
	}
	//edit the item.
	return s.Collection.EditPayload(theCollId, title, &now)
}

func (s *BulletTitleStore) TitleForMany(ids []int32) (map[int32]string, error) {
	var idStrings []string
	for _, id := range ids {
		idStr := strconv.Itoa(int(id))
		idStrings = append(idStrings, idStr)
	}
	resp, err := s.Collection.ItemsForKeys(idStrings)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}

	int32Map := make(map[int32]string)
	for k, v := range resp {
		intId, err := strToId(k.Key)
		if err != nil {
			return nil, err
		}
		int32Map[intId] = v.Payload
	}
	return int32Map, nil
}

func (s *BulletTitleStore) TitleFor(id int32) (*string, error) {
	var many []int32
	many = append(many, id)
	res, err := s.TitleForMany(many)
	if err != nil || len(res) == 0 {
		return nil, err
	}
	if len(res) != 1 {
		return nil, errors.New("too many results")
	}
	keys := make([]int32, 0, len(res))
	for k := range res {
		keys = append(keys, k)
	}
	result := res[id]
	return &result, nil
}

func (s *BulletTitleStore) RemoveItem(id int32) error {
	var keys []string
	idStr := strconv.Itoa(int(id))
	keys = append(keys, idStr)
	res, err := s.Collection.ItemsForKeys(keys)
	if err != nil || res == nil {
		return err
	}
	if len(res) != 1 {
		return errors.New("too many results")
	}
	collids := make([]bullet_stl.CollectionId, 0, len(res)) //there should be 1
	for k := range res {
		collids = append(collids, k)
	}

	return s.Collection.DeleteItems(collids)

}
