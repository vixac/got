package grove_engine

import (
	"errors"
	"time"

	"vixac.com/got/console"
	"vixac.com/got/engine"
	"vixac.com/got/engine/engine_util"
)

func (g *GroveEngine) FetchItemsBelow(lookup *engine.GidLookup, sortByPath bool, states []engine.GotState, hideUnderCollapsed bool) (*engine.GotFetchResult, error) {
	//now := time.Now()
	statesToFetch := make(map[engine.GotState]bool)
	for _, v := range states {
		statesToFetch[v] = true
	}
	//0->1 numberstore gid -> numberstore
	//0-> alias store gid -> alias store
	parentGid, err := g.GidLookup.InputToGid(lookup)
	if err != nil || parentGid == nil {
		return nil, err
	}

	parentIsRoot := parentGid.IntValue == engine_util.TheRootNoteInt32

	nodes, err := g.GroveStore.FetchBelow(parentGid)
	if err != nil {
		return nil, err
	}

	var descendantPlusParentIds []engine.GotId
	if !parentIsRoot { //unless its the root, we want to render the parent too.
		descendantPlusParentIds = append(descendantPlusParentIds, *parentGid)
	}
	for _, n := range nodes {
		descendantPlusParentIds = append(descendantPlusParentIds, n.Id)
	}
	paths, err := g.GroveStore.FetchAncestorsForMany(descendantPlusParentIds)
	if err != nil {
		return nil, err
	}
	infos, err := g.InfoStore.InfoForMany(descendantPlusParentIds)
	if err != nil {
		return nil, err
	}

	var idStrings []string
	for _, id := range descendantPlusParentIds {
		idStrings = append(idStrings, id.AasciValue)
	}
	aliases, err := g.LookupAliasForMany(idStrings)

	gotPaths := make(map[engine.GotId]engine.GotPath)
	for _, idWithPath := range paths {
		var path []engine.PathItem
		for _, pathItem := range idWithPath.Path {
			aliasOfPathItem := aliases[pathItem.AasciValue]
			path = append(path, engine.PathItem{
				Id:    pathItem.AasciValue,
				Alias: aliasOfPathItem,
			})
		}
		gotPaths[idWithPath.Id] = engine.GotPath{
			Ancestry: path,
		}
	}

	aggs, err := g.GroveStore.AggregatesOfDescendantsForMany(descendantPlusParentIds)
	if err != nil {
		return nil, err
	}

	individualStates, err := g.GroveStore.IndividualStateForMany(descendantPlusParentIds)
	if err != nil {
		return nil, err
	}
	//VX:TODO unfortunately fetches all notes here.
	longForms, err := g.LongFormStore.LongFormForMany(descendantPlusParentIds)
	if err != nil {
		return nil, err
	}

	var collapsedIds = make(map[engine.GotId]bool)
	for id, info := range infos.InfoMap {
		if info.Flags["collapsed"] == true {
			collapsedIds[id] = true
		}
	}

	var displays []engine.GotItemDisplay
	//add a display node for each id
	var parent *engine.GotItemDisplay
	for _, id := range descendantPlusParentIds {

		var thePath *engine.GotPath = nil
		path, ok := gotPaths[id]
		if ok {
			thePath = &path
		}
		info, ok := infos.InfoMap[id]
		if !ok {
			return nil, errors.New("missing info.. ")
		}
		var theAlias = ""
		alias, ok := aliases[id.AasciValue]
		if ok {
			theAlias = *alias
		}
		individual, ok := individualStates[id]
		if !ok {
			return nil, errors.New("missing agg")
		}
		state := individual

		var flagsArray []string
		for k, _ := range info.Flags {
			flagsArray = append(flagsArray, k)
		}
		summary := engine.NewSummary(state, info.Deadline, &info.CreatedDate, &info.UpdatedDate, info.Tags, flagsArray)

		agg, ok := aggs[id]
		if ok && !agg.IsLeaf {
			summary.Counts = &engine.AggCount{
				Complete: agg.Counts[engine.Complete],
				Active:   agg.Counts[engine.Active],
			}
		}
		_, hasTNote := longForms[id]

		shouldShow := true //we check state and collapsed flags on ancestors to decide if we're showing this item.

		_, stateIsBeingDisplayed := statesToFetch[*summary.State]
		if !stateIsBeingDisplayed {
			shouldShow = false
		}

		//check the ancestor paths
		if hideUnderCollapsed {
			for _, pathItem := range path.Ancestry {
				ancestorId, err := engine.NewGotId(pathItem.Id)
				if err != nil {
					return nil, err
				}
				_, isCollapsed := collapsedIds[*ancestorId]
				if isCollapsed {
					shouldShow = false
				}
			}
		}

		var displayDeadlineStr = ""
		var deadlineToken console.Token = console.TokenAlert{}

		if summary.State != nil {
			displayDeadline, t, err := engine_util.Deadline(summary.Deadline, *summary.State, time.Now())
			if err != nil {
				return nil, err
			}
			displayDeadlineStr = displayDeadline
			deadlineToken = t
		}

		if shouldShow {

			//VX:Note NumberGo is added add by EnrichWithNumberGos
			display := engine.GotItemDisplay{
				GotId:         id,
				DisplayGid:    "0" + id.AasciValue,
				Path:          thePath,
				Title:         info.Title,
				Alias:         theAlias,
				Deadline:      displayDeadlineStr,
				DeadlineToken: deadlineToken,
				SummaryObj:    &summary,
				HasTNote:      hasTNote,
			}
			if id == *parentGid {
				parent = &display
			} else {
				displays = append(displays, display)
			}
		}

	}

	var sorted []engine.GotItemDisplay
	if sortByPath {
		sorted = engine_util.SortTheseIntoDFS(displays)

	} else {
		sorted = engine_util.SortByUpdated(displays)
	}
	return engine_util.EnrichWithNumberGos(g.NumberGoStore, sorted, parent)

}
