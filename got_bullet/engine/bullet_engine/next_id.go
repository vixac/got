package bullet_engine

import (
	"fmt"
	"strconv"

	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/containers"
)

const (
	firstId       int64 = 37
	bucketId            = 100
	listName            = "next-id-list"
	separator           = ""
	latestSubject       = "latest"
)

// VX:TODO test, maybe put somewhere else too.
func (e *EngineBullet) NextId() (int64, error) {

	fmt.Printf("VX NEXT ID \n")
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

		fmt.Printf("VX NEXT ID is upserting?\n")
		str := strconv.FormatInt(firstId, 10)
		err := list.Upsert(latest, bullet_stl.ListObject{Value: str})
		if err != nil {
			return 0, err
		}
		return firstId, nil
	} else {
		fmt.Printf("VX NEXT ID got an id %s\n", currentHighest.Value)
		//now increment
		value := currentHighest.Value
		valueInt, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return 0, err
		}
		incrementedValue := valueInt + 1
		str := strconv.FormatInt(incrementedValue, 10)
		err = list.Upsert(latest, bullet_stl.ListObject{Value: str})
		if err != nil {
			return 0, err
		}
		return incrementedValue, nil
	}
}
