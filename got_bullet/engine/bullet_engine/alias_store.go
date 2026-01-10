package bullet_engine

import (
	"errors"
	"fmt"
	"log"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/containers"
	"vixac.com/got/engine"
)

type BulletAliasStore struct {
	TwoWay *bullet_stl.TwoWayListImpl
}

// The interface for all aliasing functionality. The alias takes a gid here not a lookup. Its an attempt to keep gidlookups outside this store.
type AliasStoreInterface interface {
	Lookup(alias string) (*engine.GotId, error)
	LookupAliasForGid(gid string) (*string, error)
	LookupAliasForMany(gid []string) (map[string]*string, error)
	Unalias(alias string) (*engine.GotId, error)
	Alias(gid string, alias string) (bool, error)
}

func NewBulletAliasStore(track bullet_interface.TrackClientInterface, bucketId int32) (AliasStoreInterface, error) {
	twoWay, err := bullet_stl.NewBulletTwoWayList(track, bucketId, "alias", ">", "<")
	if err != nil {
		return nil, err
	}
	return &BulletAliasStore{
		TwoWay: twoWay,
	}, nil

}

func (e *BulletAliasStore) LookupAliasForMany(gid []string) (map[string]*string, error) {

	var objects []bullet_stl.ListObject
	for _, g := range gid {
		if g == "0" { //this is the last id. This is lookup leaking into here.
			continue
		}
		if g[0] == '0' {
			fmt.Printf("VX: WTF WHY A ZERO '%s'\n", g)
		}
		objects = append(objects, bullet_stl.ListObject{Value: g})
	}

	res, err := e.TwoWay.GetSubjectsViaObjectForMany(objects)

	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	resMap := make(map[string]*string)
	for gid, alias := range res {
		if alias == nil {
			continue
		}
		resMap[gid.Value] = &alias.Value

	}
	return resMap, nil

}

func (e *BulletAliasStore) LookupAliasForGid(gid string) (*string, error) {
	if gid[0] == '0' {
		fmt.Printf("VX: WTF WHY A ZERO 2 %s\n", gid)
	}
	obj, err := e.TwoWay.GetSubjectViaObject(bullet_stl.ListObject{Value: gid})
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

	if gid[0] == '0' {
		fmt.Printf("VX: WTF WHY A ZERO 3 %s\n", gid)
		log.Fatal("NOPE")

	}
	existing, err := a.Lookup(alias)

	if err != nil {
		return false, err
	}
	if existing != nil {
		return false, errors.New("this alias is already being used. unalias it first")
	}

	return true, a.TwoWay.Upsert(bullet_stl.ListSubject{Value: alias}, bullet_stl.ListObject{Value: gid})
}
