package engine_util

import (
	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	"vixac.com/got/engine"
)

type BulletLongFormStore struct {
	Namespace int32
	Depot     bullet_interface.DepotClientInterface
}

func NewBulletLongFormStore(client bullet_interface.DepotClientInterface, namespaceId int32) (engine.LongFormStoreInterface, error) {
	return &BulletLongFormStore{
		Namespace: namespaceId,
		Depot:     client,
	}, nil
}

func (s *BulletLongFormStore) UpsertItem(id int32, block engine.LongFormBlock) error {
	//VX:TODO convert to using Collection
	/*
		namespacedId := bullet_stl.MakeNamespacedId(s.Namespace, id)
		req := bullet_interface.DepotRequest{
			Key:   namespacedId,
			Value: block.Content,
		}
		return s.Depot.DepotInsertOne(req)
	*/
	return nil
}

func (s *BulletLongFormStore) LongFormForMany(ids []int32) (map[int32]engine.LongFormBlockResult, error) {

	//VX:TODO convert to using Collection
	/*
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
		int32Map := make(map[int32]engine.LongFormBlockResult)
		for k, v := range resp.Values {
			id := bullet_stl.ParseNamespacedId(k)
			//this implemention of longform just uses a single block.
			block := engine.LongFormBlock{
				Content: v,
			}
			int32Map[id.Id] = engine.LongFormBlockResult{
				Blocks: []engine.LongFormBlock{block},
			}
		}
		return int32Map, nil
	*/
	return nil, nil
}

func (s *BulletLongFormStore) LongFormFor(id int32) (*engine.LongFormBlockResult, error) {
	/*
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
			block := engine.LongFormBlock{
				Content: title,
			}
			res := engine.LongFormBlockResult{
				Blocks: []engine.LongFormBlock{block},
			}
			return &res, nil
		}
		return nil, nil
	*/
	return nil, nil
}

func (s *BulletLongFormStore) RemoveAllItemsFromLongStore(id int32) error {
	return nil
	//VX:TODO convert to using Collection
	/*
		namespacedId := bullet_stl.MakeNamespacedId(s.Namespace, id)
		req := bullet_interface.DepotDeleteRequest{
			Key: namespacedId,
		}
		return s.Depot.DepotDeleteOne(req)
	*/
}
