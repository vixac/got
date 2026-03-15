package bullet_engine

import (
	"fmt"

	"vixac.com/got/engine"
)

func (e *EngineBullet) CreateStoreFile(filename string) error {

	allItems, err := e.FetchItemsBelow(nil, true, []engine.GotState{engine.Active, engine.Complete}, false)
	if err != nil {
		return err
	}
	fmt.Printf("There are %d items \n", len(allItems.Result))
	var everyId []int32
	for _, display := range allItems.Result {
		everyId = append(everyId, display.GotId.IntValue)
	}
	everyLongForm, err := e.LongFormStore.LongFormForMany(everyId)
	if err != nil {
		return err
	}
	fmt.Printf("VX: there are %d longforms \n", len(everyLongForm))

	//VX:TODO finish. I need to do alias first.

	/*
		var createItemRequests []engine.CreateBuckRequest
		for _, item := range allItems.Result {
			//req := engine.NewCreateBuckRequest(item.)
		}
	*/
	return nil

}

func (e *EngineBullet) RestoreFromFile(filename string) error {
	return nil
}
