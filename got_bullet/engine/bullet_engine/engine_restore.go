package bullet_engine

import (
	"fmt"

	"vixac.com/got/engine"
	"vixac.com/got/engine/engine_util"
)

func (e *EngineBullet) CreateStoreFile() error {

	allItems, err := e.FetchItemsBelow(nil, true, []engine.GotState{engine.Active, engine.Complete}, false)
	if err != nil {
		return err
	}
	//fmt.Printf("There are %d items \n", len(allItems.Result))
	var everyId []int32
	var everyAasciId []string
	for _, display := range allItems.Result {
		everyId = append(everyId, display.GotId.IntValue)
		everyAasciId = append(everyAasciId, display.GotId.AasciValue)
	}
	everyLongForm, err := e.LongFormStore.LongFormForMany(everyId)
	if err != nil {
		return err
	}
	//fmt.Printf("VX: there are %d longforms \n", len(everyLongForm))

	//VX:TODO not needed.
	/*
		allAliases, err := e.LookupAliasForMany(everyAasciId)
		if err != nil {
			return err
		}
	*/
	//fmt.Printf("VX: there are %d aliases \n", len(allAliases))

	//VX:TODO finish. I need to do alias first.

	var createItemRequests []engine.CreateBuckRequest
	for _, item := range allItems.Result {
		//item.SummaryObj.Deadline
		//lookup := engine.GidLookup{Input: item.DisplayGid}
		var alias *string = nil
		var noAlias = true
		if item.Alias != "" {
			alias = &item.Alias
			noAlias = false
		}

		var flags []string
		for f, _ := range item.SummaryObj.Flags {
			flags = append(flags, f)
		}

		var longFormPtr *engine.LongFormBlockResult = nil
		longForm, ok := everyLongForm[item.GotId.IntValue]
		if ok {
			longFormPtr = &longForm
		}

		var updated string = ""
		var createdDate string = ""
		if item.SummaryObj != nil && item.SummaryObj.UpdatedDate != nil {
			updated = item.SummaryObj.UpdatedDate.Date
		}
		if item.SummaryObj != nil && item.SummaryObj.CreatedDate != nil {
			createdDate = item.SummaryObj.CreatedDate.Date
		}
		overrides := engine.CreateOverrideSettings{
			OverrideId:  &item.GotId.IntValue,
			UpdatedDate: updated,
			CreatedDate: createdDate, //item.SummaryObj.CreatedDate.Date,
			Alias:       alias,
			NoAlias:     noAlias,
			Tags:        item.SummaryObj.Tags,
			Flags:       flags,
			LongForm:    longFormPtr,
		}
		var state engine.GotState = engine.Active
		if item.SummaryObj != nil && item.SummaryObj.State != nil {
			state = *item.SummaryObj.State
		}
		req := engine.NewCreateBuckRequest(nil, nil, item.Title, state, &overrides)
		createItemRequests = append(createItemRequests, req)
	}
	codec := &engine_util.JSONCodec[[]engine.CreateBuckRequest]{}
	json, err := codec.Encode(createItemRequests)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", json)

	return nil

}

func (e *EngineBullet) RestoreFromFile(filename string) error {
	return nil
}
