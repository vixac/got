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
	//aliases too.

	var idsWeWantAliasesFor []string
	for _, id := range descendantPlusParentIds {
		idsWeWantAliasesFor = append(idsWeWantAliasesFor, id.AasciValue)
	}

	//if a parent was passed in, we want the alias of all its ancestors so we can
	//render the parents full path
	if parentGid != nil {
		for _, idWithPath := range paths {
			if idWithPath.Id == *parentGid {
				for _, parentAncestor := range idWithPath.Path {
					idsWeWantAliasesFor = append(idsWeWantAliasesFor, parentAncestor.AasciValue)
				}

			}
		}
	}

	aliases, err := g.LookupAliasForMany(idsWeWantAliasesFor)

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
		if id.AasciValue != parentGid.AasciValue && info.Flags["collapsed"] == true {
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

		//VX:Note days ago str
		//createdStr, err := humanizeDateTime(summary.CreatedDate, now)
		createdStr, err := summary.CreatedDate.JsonDateToReadable()
		if err != nil {
			return nil, err
		}
		if shouldShow {
			isParent := id == *parentGid
			//VX:Note NumberGo is added add by EnrichWithNumberGos
			display := engine.GotItemDisplay{
				GotId:         id,
				Created:       createdStr,
				DisplayGid:    id.DisplayAasci(),
				Path:          thePath,
				Title:         info.Title,
				Alias:         theAlias,
				Deadline:      displayDeadlineStr,
				DeadlineToken: deadlineToken,
				SummaryObj:    &summary,
				HasTNote:      hasTNote,
				IsParent:      isParent,
			}
			if isParent {
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

func humanizeDateTime(date *engine.DateTime, now time.Time) (string, error) {
	dateUnix, err := date.ToDate()
	if err != nil {
		return "", err
	}
	var dateStr = ""
	if dateUnix != nil {
		dateStr, _ = console.HumanizeDate(time.Time(*dateUnix), now)
	}
	return dateStr, nil
}
