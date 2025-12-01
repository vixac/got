package bullet_engine

import (
	"errors"
	"fmt"
	"strconv"

	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/containers"
	"vixac.com/got/engine"
)

const (
	firstId       int32 = 360 //maps to a1 with bullet_stl.AasciBulletIdToInt
	bucketId            = 100
	listName            = "next-id-list"
	separator           = ""
	latestSubject       = "latest"
)

// VX:TODO test, maybe put somewhere else too.
func (e *EngineBullet) NextId() (int32, error) {

	list, err := bullet_stl.NewBulletOneWayList(e.Client, bucketId, listName, separator)
	if err != nil {
		return 0, err
	}
	latest := bullet_stl.ListSubject{Value: latestSubject}
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
