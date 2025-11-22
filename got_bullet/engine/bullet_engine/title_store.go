package bullet_engine

import (
	"errors"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
)

type TitleStoreInterface interface {
	AddItem(id int64, title string) error
	TitleFor(id int64) (*string, error)
	RemoveItem(id int64) error
}
type BulletTitleStore struct {
	Depot bullet_interface.DepotClientInterface
}

func NewBulletTitleStore(client bullet_interface.DepotClientInterface) (*BulletTitleStore, error) {
	return &BulletTitleStore{
		Depot: client,
	}, nil
}

func (s *BulletTitleStore) AddItem(id int64, title string) error {
	req := bullet_interface.DepotRequest{
		Key:   id,
		Value: title,
	}
	return s.Depot.DepotInsertOne(req)
}
func (s *BulletTitleStore) TitleFor(id int64) (*string, error) {
	keys := []int64{id}
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

	if title, ok := resp.Values[id]; ok {
		return &title, nil
	}
	return nil, nil
}

func (s *BulletTitleStore) RemoveItem(id int64) error {
	return errors.New("delete depot not working yet")

}
