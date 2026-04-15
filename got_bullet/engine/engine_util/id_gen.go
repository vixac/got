package engine_util

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/containers"
	"vixac.com/got/engine"
)

// 0 is a bit like a null terminator character. Beacuse ancestor list is forward only
// duplicate objects aren't a problem. We'll give this value to all nodes that are infact a leaf node.
var (
	TheRootNoteInt32 int32 = 0
	//if no parent is provided, then the root ancestor is provided.
	TheRootNode = bullet_stl.ListSubject{
		Value: "0",
	}
)

const (
	firstId int64 = 360 //maps to a1 with bullet_stl.AasciBulletIdToInt
)

type IdGenerator struct {
	Client        bullet_interface.BulletClientInterface
	BucketId      int32
	ListName      string
	Separator     string
	LatestSubject string
}

// Change of setup, here we're allowing the class to define its own bucketId etc. Perhaps not ideal.
func NewIdBulletGenerator(client bullet_interface.BulletClientInterface, bucketId int32, listName string, separator string, latestSubject string) engine.IdGeneratorInterface {
	return &IdGenerator{
		Client:        client,
		BucketId:      bucketId,
		ListName:      listName,
		Separator:     separator,
		LatestSubject: latestSubject,
	}
}

func (i *IdGenerator) SetLastIdIfLower(newId int64) error {
	list, err := bullet_stl.NewBulletOneWayList(i.Client, i.BucketId, i.ListName, i.Separator)
	if err != nil {
		return err
	}
	latest := bullet_stl.ListSubject{Value: i.LatestSubject}
	currentHighest, err := list.GetObject(latest)
	if err != nil {
		fmt.Printf("VX next Id failed at get object. %s\n", err.Error())
		return err
	}
	//basecase, we just use this id.
	if currentHighest == nil {
		return i.setNextId(newId)
	}
	valueInt, err := strconv.ParseInt(currentHighest.Value, 10, 32)
	if newId < valueInt {
		return nil

	}
	return i.setNextId(newId)
}

func (i *IdGenerator) LastId() (int64, error) {
	list, err := bullet_stl.NewBulletOneWayList(i.Client, i.BucketId, i.ListName, i.Separator)
	if err != nil {
		return 0, err
	}
	latest := bullet_stl.ListSubject{Value: i.LatestSubject}
	currentHighest, err := list.GetObject(latest)
	if err != nil {
		fmt.Printf("VX NEXT ID failed at get object. %s\n", err.Error())
		return 0, err
	}
	if currentHighest == nil {
		return 0, errors.New("no lastId found.")
	}
	valueInt, err := strconv.ParseInt(currentHighest.Value, 10, 32)
	return valueInt, err

}

// saves the next id
func (i IdGenerator) setNextId(next int64) error {
	//VX:TODO save that list etc.
	list, err := bullet_stl.NewBulletOneWayList(i.Client, i.BucketId, i.ListName, i.Separator)
	if err != nil {
		return err
	}
	latest := bullet_stl.ListSubject{Value: i.LatestSubject}
	str := fmt.Sprint(next)
	err = list.Upsert(latest, bullet_stl.ListObject{Value: str})
	if err != nil {
		return err
	}
	return nil
}

// VX:TODO test
func (i *IdGenerator) NextId() (int64, error) {

	list, err := bullet_stl.NewBulletOneWayList(i.Client, i.BucketId, i.ListName, i.Separator)
	if err != nil {
		return 0, err
	}
	latest := bullet_stl.ListSubject{Value: i.LatestSubject}
	currentHighest, err := list.GetObject(latest)
	if err != nil {
		fmt.Printf("VX next Id failed at get object.. Next it will return  %s\n", err.Error())
		return 0, err
	}

	//base case, start at the beginning.
	if currentHighest == nil {
		err = i.setNextId(firstId)
		if err != nil {
			return 0, err
		}
		return firstId, nil
	} else {
		//now increment
		value := currentHighest.Value
		valueInt, err := strconv.ParseInt(value, 10, 32)

		if err != nil {
			return 0, err
		}
		incrementedValue := valueInt + 1
		if !engine.FitsInInt32(incrementedValue) {
			return 0, errors.New("this id space has exhausted int32")
		}
		err = i.setNextId(incrementedValue)
		return incrementedValue, nil
	}
}
