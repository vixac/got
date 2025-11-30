package bullet_engine

import (
	"errors"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/ids"
	"vixac.com/got/engine"
)

// VX:TODO https://chatgpt.com/s/t_69236b33d8748191a6dd3c8e67a46a5e

// VX:TODO this should probbaly get gids or gotids or whatever.
// / This is to make clear that Agg doesn't know about GotIds specifically, although no doubt intValue GotIds will be used 1-1 for Agg.
// / Agg will just operate with int32 and namespaces which are also int32, and will use firbolg_clients MakeNamespacedId to create the int64 that depot needs.
type SummaryId int32

func NewSummaryId(gotId engine.GotId) SummaryId {
	return SummaryId(gotId.IntValue)
}

type SummaryStoreInterface interface {
	UpsertAggregate(id SummaryId, agg Aggregate) error
	UpsertManyAggregates(aggs map[SummaryId]Aggregate) error
	Fetch(ids []SummaryId) (map[SummaryId]Aggregate, error)
	Delete(ids []SummaryId) error
}

// namespaces are like bucket Ids but they move separately so the fact that they're both int32 is coincidence. namespcae is a way to
// use the int64 space of depot ids to be <namespace><id>. This is fine because 2,147,483,647 is the positive total of int32. Thats plenty for got. If we need to host spaces higher than that, we probably want to
// break stuff up into spearated sections and use mirroring.
type BulletSummaryStore struct {
	codec     Codec[Aggregate]
	Client    bullet_interface.DepotClientInterface
	Namespace int32
}

func NewBulletSummaryStore(codec Codec[Aggregate], client bullet_interface.DepotClientInterface, namespace int32) (SummaryStoreInterface, error) {
	return &BulletSummaryStore{
		codec:     codec,
		Client:    client,
		Namespace: namespace,
	}, nil
}

func (a *BulletSummaryStore) aggIdToNamespacedId(id SummaryId) int64 {
	return bullet_stl.MakeNamespacedId(a.Namespace, int32(id))
}

func (a *BulletSummaryStore) namespacedIdToAgg(spaced int64) SummaryId {
	namespaced := bullet_stl.ParseNamespacedId(spaced)
	return SummaryId(namespaced.Id)
}

func (a *BulletSummaryStore) UpsertManyAggregates(aggs map[SummaryId]Aggregate) error {

	var reqs []bullet_interface.DepotRequest
	for id, agg := range aggs {
		json, err := a.codec.Encode(agg)
		if err != nil {
			return err
		}
		spaced := a.aggIdToNamespacedId(id)
		reqs = append(reqs, bullet_interface.DepotRequest{
			Key:   spaced,
			Value: json,
		})
	}
	return a.Client.DepotUpsertMany(reqs)
}

// VX:TODO RM or call many if we want to keep it
func (a *BulletSummaryStore) UpsertAggregate(id SummaryId, agg Aggregate) error {

	json, err := a.codec.Encode(agg)
	if err != nil {
		return err
	}
	spaced := a.aggIdToNamespacedId(id)
	req := bullet_interface.DepotRequest{
		Key:   spaced,
		Value: json,
	}
	return a.Client.DepotInsertOne(req)
}

func (a *BulletSummaryStore) Fetch(ids []SummaryId) (map[SummaryId]Aggregate, error) {
	var keys []int64
	for _, id := range ids {
		spaced := a.aggIdToNamespacedId(id)
		keys = append(keys, spaced)
	}
	manyReq := bullet_interface.DepotGetManyRequest{
		Keys: keys,
	}
	resp, err := a.Client.DepotGetMany(manyReq)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}

	result := make(map[SummaryId]Aggregate)
	for k, v := range resp.Values {
		aggObj := &Aggregate{}
		err := a.codec.Decode(v, aggObj)
		if err != nil {
			return nil, err
		}
		aggId := a.namespacedIdToAgg(k)
		result[aggId] = *aggObj
	}
	return result, nil

}

func (a *BulletSummaryStore) Delete(ids []SummaryId) error {
	return errors.New("not impl")
}
