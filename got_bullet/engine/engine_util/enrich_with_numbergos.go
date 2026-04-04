package engine_util

import "vixac.com/got/engine"

func EnrichWithNumberGos(store engine.NumberGoStoreInterface, summaries []engine.GotItemDisplay, parent *engine.GotItemDisplay) (*engine.GotFetchResult, error) {

	var expandedSummaries []engine.GotItemDisplay
	var pairs []engine.NumberGoPair
	//here we enrich the itemdisplays by adding the number go, now that we know the sort order.
	for i, s := range summaries {

		num := i + 1
		pairs = append(pairs, engine.NumberGoPair{
			Number: num,
			Gid:    s.GotId.AasciValue,
		})

		copy := s
		copy.NumberGo = num
		expandedSummaries = append(expandedSummaries, copy)
	}

	//the number go order is saved so it can be used in subsequent calls
	err := store.AssignNumberPairs(pairs)
	if err != nil {
		return nil, err
	}

	//the summaries injected dont have number go assigned.
	res := engine.GotFetchResult{Result: expandedSummaries, Parent: parent}
	return &res, nil

}
