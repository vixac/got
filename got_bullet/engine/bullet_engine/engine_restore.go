package bullet_engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"

	"vixac.com/got/engine"
	"vixac.com/got/engine/engine_util"
)

func (e *EngineBullet) CreateStoreFile() error {

	allItems, err := e.FetchItemsBelow(nil, true, []engine.GotState{engine.Active, engine.Complete}, false)
	if err != nil {
		return err
	}
	var everyId []engine.GotId
	var everyAasciId []string
	for _, display := range allItems.Result {
		everyId = append(everyId, display.GotId)
		everyAasciId = append(everyAasciId, display.GotId.AasciValue)
	}
	everyLongForm, err := e.LongFormStore.LongFormForMany(everyId)
	if err != nil {
		return err
	}

	//sorted for top level first.
	sortedItemLList := allItems.Result
	sort.Slice(sortedItemLList, func(i, j int) bool {
		return sortedItemLList[i].Path.Depth() < sortedItemLList[j].Path.Depth()
	})

	var createItemRequests []engine.CreateBuckRequest
	for _, item := range sortedItemLList {
		var lookup *engine.GidLookup = nil
		if item.Path != nil && item.Path.Depth() > 0 {
			lastItemInPath := item.Path.Ancestry[item.Path.Depth()-1]
			if lastItemInPath.Id != "0" { //this id actually gets looked up so we don't want that.
				lookup = &engine.GidLookup{Input: "0" + lastItemInPath.Id}
			}
		}
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
		longForm, ok := everyLongForm[item.GotId]
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
		req := engine.NewCreateBuckRequest(lookup, nil, item.Title, state, &overrides)
		createItemRequests = append(createItemRequests, req)
	}
	codec := &engine_util.JSONCodec[[]engine.CreateBuckRequest]{}
	json, err := codec.Encode(createItemRequests)
	if err != nil {
		return err
	}
	//here we print the restore to std out (VX:TODO return to caller so it can use the deps.Printer)
	fmt.Printf("%s\n", json)
	return nil

}

func (e *EngineBullet) RestoreFromFile(filename string) error {

	var res []engine.CreateBuckRequest
	var data []byte
	data, _ = ioutil.ReadFile(filename)

	err := json.Unmarshal(data, &res)
	if err != nil {
		return err
	}
	for _, createReq := range res {
		gotId, err := e.CreateBuck(createReq)
		if err != nil {
			return err
		}
		fmt.Printf("VX: restored %s\n", gotId.AasciValue)
	}
	return nil
}
