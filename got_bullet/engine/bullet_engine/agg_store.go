package bullet_engine

import (
	"errors"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/ids"
	"vixac.com/got/engine"
)

// VX:TODO https://chatgpt.com/s/t_69236b33d8748191a6dd3c8e67a46a5e
// / This is to make clear that Agg doesn't know about GotIds specifically, although no doubt intValue GotIds will be used 1-1 for Agg.
// / Agg will just operate with int32 and namespaces which are also int32, and will use firbolg_clients MakeNamespacedId to create the int64 that depot needs.
type AggId int32

type AggStoreInterface interface {
	UpsertAggregate(id AggId, agg Aggregate) error
	Fetch(ids []AggId) (map[AggId]Aggregate, error)
	Delete(ids []AggId) error
}

// First pass of the kinds of things we'll count
type AggCount struct {
	Complete int `json:"c,omitempty"`
	Active   int `json:"a,omitempty"`
	Notes    int `json:"n,omitempty"`
}

type Deadline struct {
	Date string `json:"d,omitempty"`
}

type DatedTask struct {
	Deadline Deadline `json:"d"`
	Id       AggId    `json:"i,omitempty"`
}

// VX:TODO its either state OR its counts.
// deadline is separate. Maybe it doesn't belong here but we'll see.
type Aggregate struct {
	State    engine.GotState `json:"s,omitempty"`
	Counts   AggCount        `json:"c"`
	Deadline Deadline        `json:"d"`
}

func (c AggCount) ChangeState(state engine.GotState, inc int) AggCount {
	comp := c.Complete
	active := c.Active
	notes := c.Notes
	if state == engine.Active {
		active += inc
	} else if state == engine.Complete {
		comp += inc
	} else if state == engine.Note {
		notes += inc
	}
	return AggCount{
		Complete: comp,
		Active:   active,
		Notes:    notes,
	}
}
func (c AggCount) changeActive(inc int) AggCount {
	return AggCount{
		c.Complete,
		c.Active + inc,
		c.Notes,
	}
}
func (c AggCount) changeNotes(inc int) AggCount {
	return AggCount{
		c.Complete,
		c.Active,
		c.Notes + inc,
	}
}
func (c AggCount) changeComplete(inc int) AggCount {
	return AggCount{
		c.Complete + inc,
		c.Active,
		c.Notes,
	}
}
func (a *Aggregate) UpdatedCount(newCount AggCount) Aggregate {
	return Aggregate{
		State:    a.State,
		Counts:   newCount,
		Deadline: a.Deadline,
	}
}

// namespaces are like bucket Ids but they move separately so the fact that they're both int32 is coincidence. namespcae is a way to
// use the int64 space of depot ids to be <namespace><id>. This is fine because 2,147,483,647 is the positive total of int32. Thats plenty for got. If we need to host spaces higher than that, we probably want to
// break stuff up into spearated sections and use mirroring.
type BulletAggStore struct {
	codec     Codec[Aggregate]
	Client    bullet_interface.DepotClientInterface
	Namespace int32
}

func NewBulletAggStore(codec Codec[Aggregate], client bullet_interface.DepotClientInterface, namespace int32) (AggStoreInterface, error) {
	return &BulletAggStore{
		codec:     codec,
		Client:    client,
		Namespace: namespace,
	}, nil
}

func (a *BulletAggStore) aggIdToNamespacedId(id AggId) int64 {
	return bullet_stl.MakeNamespacedId(a.Namespace, int32(id))
}

func (a *BulletAggStore) namespacedIdToAgg(spaced int64) AggId {
	namespaced := bullet_stl.ParseNamespacedId(spaced)
	return AggId(namespaced.Id)
}

func (a *BulletAggStore) UpsertAggregate(id AggId, agg Aggregate) error {
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

func (a *BulletAggStore) Fetch(ids []AggId) (map[AggId]Aggregate, error) {
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

	result := make(map[AggId]Aggregate)
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

func (a *BulletAggStore) Delete(ids []AggId) error {
	return errors.New("not impl")
}
