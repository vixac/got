package bullet_engine

import (
	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/containers"
	"vixac.com/got/engine"
)

type BulletAliasStore struct {
	TwoWay *bullet_stl.TwoWayListImpl
}

func NewBulletAliasStore(track bullet_interface.TrackClientInterface, bucketId int32) (engine.GotAliasInterface, error) {
	twoWay, err := bullet_stl.NewBulletTwoWayList(track, bucketId, "alias", ">", "<")
	if err != nil {
		return nil, err
	}
	return &BulletAliasStore{
		TwoWay: twoWay,
	}, nil

}

func (e *BulletAliasStore) LookupAliasForGid(gid string) (*string, error) {
	obj, err := e.TwoWay.GetOSubjectViaObject(bullet_stl.ListObject{Value: gid})
	if err != nil {
		return nil, err
	}
	if obj == nil {
		return nil, nil
	}
	return &obj.Value, nil

}
func (a *BulletAliasStore) Lookup(alias string) (*engine.GotId, error) {
	object, err := a.TwoWay.GetObjectViaSubject(bullet_stl.ListSubject{Value: alias})
	if err != nil {
		return nil, err
	}
	if object == nil {
		return nil, nil
	}
	return engine.NewGotId(object.Value)
}

// VX:TODO no need for gotid here
func (a *BulletAliasStore) Unalias(alias string) (*engine.GotId, error) {
	return nil, a.TwoWay.DeleteViaSub(bullet_stl.ListSubject{Value: alias})
}

// VX:TODO no need for bool here.
func (a *BulletAliasStore) Alias(gid string, alias string) (bool, error) {
	return true, a.TwoWay.Upsert(bullet_stl.ListSubject{Value: alias}, bullet_stl.ListObject{Value: gid})
}
