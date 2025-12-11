package bullet_engine

import (
	"errors"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
)

type TitleStoreInterface interface {
	UpsertItem(id int32, title string) error
	TitleFor(id int32) (*string, error)

	TitleForMany(ids []int32) (map[int32]string, error)
	RemoveItem(id int32) error
}
type BulletTitleStore struct {
	Depot bullet_interface.DepotClientInterface
}

func NewBulletTitleStore(client bullet_interface.DepotClientInterface) (TitleStoreInterface, error) {
	return &BulletTitleStore{
		Depot: client,
	}, nil
}

func (s *BulletTitleStore) UpsertItem(id int32, title string) error {
	req := bullet_interface.DepotRequest{
		Key:   int64(id),
		Value: title,
	}
	return s.Depot.DepotInsertOne(req)
}

func (s *BulletTitleStore) TitleForMany(ids []int32) (map[int32]string, error) {

	var int64Ids []int64
	for _, v := range ids {
		int64Ids = append(int64Ids, int64(v))
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
		int32Map[int32(k)] = v
	}
	return int32Map, nil
}

func (s *BulletTitleStore) TitleFor(id int32) (*string, error) {
	keys := []int64{int64(id)}
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

	if title, ok := resp.Values[int64(id)]; ok {
		return &title, nil
	}
	return nil, nil
}

func (s *BulletTitleStore) RemoveItem(id int32) error {
	return errors.New("delete depot not working yet")

}
