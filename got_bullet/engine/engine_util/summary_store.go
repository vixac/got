package engine_util

import (
	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/containers"
	"vixac.com/got/engine"
)

func NewSummaryId(gotId engine.GotId) engine.SummaryId {
	return engine.SummaryId(gotId.IntValue)
}

type BulletSummaryStore struct {
	codec      Codec[engine.Summary]
	Collection bullet_stl.Collection
}

func NewBulletSummaryStore(bucketId int32, track bullet_interface.TrackClientInterface, depot bullet_interface.DepotClientInterface, codec Codec[engine.Summary]) (engine.SummaryStoreInterface, error) {
	coll := bullet_stl.NewBulletCollection(bucketId, track, depot)
	return &BulletSummaryStore{
		codec:      codec,
		Collection: coll,
	}, nil
}

func (a *BulletSummaryStore) UpsertManySummaries(aggs map[engine.SummaryId]engine.Summary) error {
	var idStrings []string
	for id := range aggs {
		idStrings = append(idStrings, idToStr(int32(id)))
	}
	existing, err := a.Collection.ItemsForKeys(idStrings)
	if err != nil {
		return err
	}
	existingByKey := make(map[string]bullet_stl.CollectionId)
	for collId := range existing {
		existingByKey[collId.Key] = collId
	}
	for id, summary := range aggs {
		payload, err := a.codec.Encode(summary)
		if err != nil {
			return err
		}
		key := idToStr(int32(id))
		if collId, ok := existingByKey[key]; ok {
			err = a.Collection.EditPayload(collId, payload)
		} else {
			_, err = a.Collection.CreateItemUnder(key, payload)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *BulletSummaryStore) UpsertSummary(id engine.SummaryId, agg engine.Summary) error {
	return a.UpsertManySummaries(map[engine.SummaryId]engine.Summary{id: agg})
}

func (a *BulletSummaryStore) Fetch(ids []engine.SummaryId) (map[engine.SummaryId]engine.Summary, error) {
	var idStrings []string
	for _, id := range ids {
		idStrings = append(idStrings, idToStr(int32(id)))
	}
	resp, err := a.Collection.ItemsForKeys(idStrings)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	result := make(map[engine.SummaryId]engine.Summary)
	for collId, payload := range resp {
		id, err := strToId(collId.Key)
		if err != nil {
			return nil, err
		}
		var summary engine.Summary
		err = a.codec.Decode(payload, &summary)
		if err != nil {
			return nil, err
		}
		result[engine.SummaryId(id)] = summary
	}
	return result, nil
}

func (a *BulletSummaryStore) Delete(ids []engine.SummaryId) error {
	var idStrings []string
	for _, id := range ids {
		idStrings = append(idStrings, idToStr(int32(id)))
	}
	res, err := a.Collection.ItemsForKeys(idStrings)
	if err != nil || res == nil {
		return err
	}
	var collIds []bullet_stl.CollectionId
	for k := range res {
		collIds = append(collIds, k)
	}
	return a.Collection.DeleteItems(collIds)
}
