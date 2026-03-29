package engine_util

import (
	"sort"

	"vixac.com/got/engine"
)

func NotesForMany(e engine.GotFetchInterface, lookup *engine.GidLookup, longform engine.LongFormStoreInterface) (*engine.LongFormBlockResult, error) {
	items, err := e.FetchItemsBelow(lookup, false, []engine.GotState{engine.Active, engine.Complete}, false)
	if err != nil {
		return nil, err
	}

	var gotIds []engine.GotId
	for _, v := range items.Result {
		gotIds = append(gotIds, v.GotId)
	}
	mapOfResults, err := longform.LongFormForMany(gotIds)
	if err != nil {
		return nil, err
	}
	var blocks []engine.LongFormBlock
	for _, v := range mapOfResults {
		blocks = append(blocks, v.Blocks...)
	}

	//sorted with latest first
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].Edited.Before(blocks[j].Edited)
	})
	res := engine.LongFormBlockResult{
		Blocks: blocks,
	}

	return &res, nil

}
