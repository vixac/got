package engine_util

import (
	"errors"
	"strings"
	"unicode"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/containers"
	"vixac.com/got/engine"
)

type BulletAliasStore struct {
	TwoWay *bullet_stl.TwoWayListImpl
}

func NewBulletAliasStore(track bullet_interface.TrackClientInterface, bucketId int32) (engine.AliasStoreInterface, error) {
	twoWay, err := bullet_stl.NewBulletTwoWayList(track, bucketId, "alias", ">", "<")
	if err != nil {
		return nil, err
	}
	return &BulletAliasStore{
		TwoWay: twoWay,
	}, nil

}

func CheckNumber(p []byte) bool {
	r := string(p)

	sep := 0
	for i, b := range r {

		if unicode.IsNumber(b) {
			continue
		}

		if b == '-' {
			if i != 0 {
				return false
			}
			continue
		}

		if b == '.' {
			if sep > 0 {
				return false
			}
			sep++
			continue
		}

		return false
	}

	return len(r) > 0
}
func IsValidAlias(input string) bool {
	if len(input) == 0 {
		return false
	}
	spaces := strings.Contains(input, " ")
	if spaces {
		return false
	}
	bytes := []byte(input)
	firstCharIsNumber := CheckNumber([]byte{bytes[0]})
	if firstCharIsNumber {
		return false
	}
	return true

}

func (e *BulletAliasStore) LookupAliasForMany(gid []string) (map[string]*string, error) {
	var objects []bullet_stl.ListObject
	for _, g := range gid {
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
func (a *BulletAliasStore) Alias(id engine.GotId, alias string) (bool, error) {

	existing, err := a.Lookup(alias)

	if err != nil {
		return false, err
	}
	if existing != nil {
		return false, errors.New("this alias is already being used. unalias it first")
	}

	return true, a.TwoWay.Upsert(bullet_stl.ListSubject{Value: alias}, bullet_stl.ListObject{Value: id.AasciValue})
}
