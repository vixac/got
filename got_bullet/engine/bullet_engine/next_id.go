package bullet_engine

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/containers"
	"vixac.com/got/engine"
)

const (
	firstId int32 = 360 //maps to a1 with bullet_stl.AasciBulletIdToInt
	//bucketId            = 100
//	listName            = "next-id-list"
//	separator           = "" //VX:TODO this should be something else. I should migrate.
//	latestSubject       = "latest"
)

type IdGeneratorInterface interface {
	LastId() (int32, error) //fetches the last createdId
	NextId() (int32, error) //creates a new id, stores it as the lastId, and returns it
}
type IdGenerator struct {
	Client        bullet_interface.BulletClientInterface
	BucketId      int32
	ListName      string
	Separator     string
	LatestSubject string
}

// Change of setup, here we're allowing the class to define its own bucketId etc. Perhaps not ideal.
func NewIdBulletGenerator(client bullet_interface.BulletClientInterface, bucketId int32, listName string, separator string, latestSubject string) IdGeneratorInterface {
	return &IdGenerator{
		Client:        client,
		BucketId:      bucketId,
		ListName:      listName,
		Separator:     separator,
		LatestSubject: latestSubject,
	}
}

// VX:TODO
// VX:TODO this is all bad things.
func (i *IdGenerator) LastId() (int32, error) {
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
	valueInt, err := strconv.ParseInt(currentHighest.Value, 10, 32)
	return int32(valueInt), err

}

// VX:TODO test, maybe put somewhere else too.
func (i *IdGenerator) NextId() (int32, error) {

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

	//base case, start at the beginning.
	if currentHighest == nil {
		str := fmt.Sprint(firstId)
		err := list.Upsert(latest, bullet_stl.ListObject{Value: str})
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
		str := strconv.FormatInt(incrementedValue, 10)
		err = list.Upsert(latest, bullet_stl.ListObject{Value: str})
		if err != nil {
			return 0, err
		}
		return int32(incrementedValue), nil
	}
}
