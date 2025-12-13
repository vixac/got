package bullet_engine

import (
	"errors"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/ids"
)

type LongFormStoreInterface interface {
	UpsertItem(id int32, title string) error
	LongFormFor(id int32) (*string, error)
	LongFormForMany(ids []int32) (map[int32]string, error)
	RemoveItem(id int32) error
}
type BulletLongFormStore struct {
	Namespace int32
	Depot     bullet_interface.DepotClientInterface
}

func NewBulletLongFormStore(client bullet_interface.DepotClientInterface, namespaceId int32) (LongFormStoreInterface, error) {
	return &BulletLongFormStore{
		Namespace: namespaceId,
		Depot:     client,
	}, nil
}

// VX:TODO this implementation is identical to title store. Its just a string id pair.
// VX:TODO make this
func (s *BulletLongFormStore) UpsertItem(id int32, title string) error {
	namespacedId := bullet_stl.MakeNamespacedId(s.Namespace, id)
	req := bullet_interface.DepotRequest{
		Key:   namespacedId,
		Value: title,
	}
	return s.Depot.DepotInsertOne(req)
}

func (s *BulletLongFormStore) LongFormForMany(ids []int32) (map[int32]string, error) {

	var int64Ids []int64
	for _, v := range ids {
		namespacedId := bullet_stl.MakeNamespacedId(s.Namespace, v)
		int64Ids = append(int64Ids, namespacedId)
	}
	req := bullet_interface.DepotGetManyRequest{
		Keys: int64Ids,
	}
	resp, err := s.Depot.DepotGetMany(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	int32Map := make(map[int32]string)
	for k, v := range resp.Values {
		id := bullet_stl.ParseNamespacedId(k)
		int32Map[id.Id] = v
	}
	return int32Map, nil
}

func (s *BulletLongFormStore) LongFormFor(id int32) (*string, error) {
	namespacedId := bullet_stl.MakeNamespacedId(s.Namespace, id)
	keys := []int64{namespacedId}
	req := bullet_interface.DepotGetManyRequest{
		Keys: keys,
	}
	resp, err := s.Depot.DepotGetMany(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}

	if title, ok := resp.Values[namespacedId]; ok {
		return &title, nil
	}
	return nil, nil
}

func (s *BulletLongFormStore) RemoveItem(id int32) error {
	return errors.New("delete depot not working yet")
}
