package bullet_engine

import (
	"errors"
	"fmt"
	"time"

	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/ids"
	"vixac.com/got/engine"
)

// lets rewrite this maybe.
func (e *EngineBullet) FetchItemsBelow(lookup *engine.GidLookup, sortByPath bool, states []engine.GotState) (*engine.GotFetchResult, error) {

	now := time.Now()
	statesToFetch := make(map[engine.GotState]bool)
	for _, v := range states {
		statesToFetch[v] = true
	}
	//0->1 numberstore gid -> numberstore
	//0-> alias store gid -> alias store
	parentGid, err := e.GidLookup.InputToGid(lookup)
	if err != nil || parentGid == nil {
		return nil, err
	}

	parentIsRoot := parentGid.IntValue == TheRootNoteInt32

	//1.gid->ancestor (object -> subject)
	//2.all descendants: allpairs for full key

	//VX:TODO we should be able to apply the state filtering here so that complete items aren't fetched unless necessary. Because many complete , few active.
	all, err := e.AncestorList.FetchImmediatelyUnder(*parentGid)
	if err != nil {
		return nil, err
	}

	// we tolerate having no children here because we want to render the parent no matter what.
	var allIds map[string]AncestorLookupResult
	if all != nil {
		allIds = all.Ids
	} else {
		allIds = make(map[string]AncestorLookupResult)
	}

	var plusParent = 0
	if !parentIsRoot {
		plusParent = 1
	}
	//get string ids of all items to do the alias lookup
	stringIds := make([]string, len(allIds)+plusParent)

	i := 0
	for k := range allIds {
		stringIds[i] = k
		i++
	}
	//fetch theparent too if its not the root node.
	if !parentIsRoot {
		stringIds[len(allIds)] = parentGid.AasciValue
	}

	aliasMap, err := e.AliasStore.LookupAliasForMany(stringIds)
	if err != nil {
		return nil, err
	}

	//VX:TODO so there's no ancestor path available for the parent. its a bug basically. Because aren't fetching the PARENT when we call FetchImmediatelyUnder. In theory we could return the parent perhaps? Not sure.
	var intIds []int32
	ancestorPaths := make(map[int32]engine.GotPath)
	if all != nil {
		for id, ancestorLookup := range allIds {

			intId, err := bullet_stl.AasciBulletIdToInt(id)
			if err != nil {
				return nil, err
			}
			intIds = append(intIds, int32(intId))

			path := ancestorPathFor(&ancestorLookup, aliasMap)

			if path != nil {
				ancestorPaths[int32(intId)] = *path
			}
		}
	}

	//this is a bit of a workaround beacuse the intIds come from the for loop above, which is on the result of FetchimmediatelyUnder, which has no parent.
	if !parentIsRoot {
		intIds = append(intIds, parentGid.IntValue)
	}
	var summaryIds []engine.SummaryId
	for _, v := range intIds {
		summaryIds = append(summaryIds, engine.SummaryId(v))
	}
	summaries, err := e.SummaryStore.Fetch(summaryIds)
	if err != nil {
		return nil, err
	}

	var collapsedIds = make(map[engine.SummaryId]bool)
	for _, id := range summaryIds {
		summary, ok := summaries[id]
		if !ok {
			continue //VX:TODO this is an error
		}
		if summary.Flags["collapsed"] {
			fmt.Printf("VX: THIS IS COLLAPSED: %d\n", id)
			collapsedIds[id] = true
		}
	}

	//titleStore: allIds -> title
	titles, err := e.TitleStore.TitleForMany(intIds)
	if err != nil {
		return nil, err
	}

	//just needed to see if we present the note emoji. Unfortunately we're loading
	//the actual notes on here.
	//VX:TODO we just need to know if theres a note, not load the content.
	longForms, err := e.LongFormForMany(intIds)
	if err != nil {
		return nil, err
	}

	var parentItemDisplay *engine.GotItemDisplay = nil

	//build ancestors
	var itemDisplays []engine.GotItemDisplay
	for k, v := range titles {

		stringId, err := bullet_stl.BulletIdIntToaasci(int64(k)) //VX:TODO can we just look this up from above?
		if err != nil {
			return nil, err
		}

		var alias string = ""
		found, ok := aliasMap[stringId]
		if ok {
			alias = *found
		}
		var path *engine.GotPath = nil
		if foundPath, ok := ancestorPaths[k]; ok {
			path = &foundPath
		}

		gotId, err := engine.NewGotId(stringId)
		if err != nil {
			return nil, err
		}
		summaryId := NewSummaryId(*gotId)
		summary, ok := summaries[summaryId]
		if !ok {
			return nil, errors.New("missing summary in fetchItems Below")
		}

		//here we filter complete leafs from the jobs list, and their notes.
		//VX:Note we want to have completes
		//not even appear in the search, because thats more scalable.

		_, hasLongForm := longForms[k]

		//this is the parent, so we populate parentItemDisplay and then continnue.
		if k == parentGid.IntValue { //we will render parents separtely
			displayItem, err := itemDisplay(summary, now, *gotId, v, alias, path, hasLongForm)
			if err != nil {
				return nil, err
			}
			parentItemDisplay = displayItem
			continue
		}

		//now we decide we're showing the descendant.

		//this is shouldShow logic.
		pathLen := len(path.Ancestry)
		var isParentComplete = false
		if pathLen > 0 {
			parentId := path.Ancestry[pathLen-1].Id
			backToInt, _ := bullet_stl.AasciBulletIdToInt(parentId) //so many conversions. VX:TODO just create a 2 way map or whatever. Maybe that map is its own type.
			parentSummary, ok := summaries[engine.SummaryId(backToInt)]
			if ok {
				if parentSummary.State != nil && *parentSummary.State == engine.Complete {
					isParentComplete = true
				}
			}
		}

		var shouldShow = false
		shouldFetchComplete := statesToFetch[engine.GotState(engine.Complete)]
		if summary.State == nil {
			shouldShow = true
		} else {
			if statesToFetch[*summary.State] {
				shouldShow = true
			}
			//this is an edge case where if you're not rendering complete nodes, you also don't want to render notes under complete nodes.
			if *summary.State == engine.Note && isParentComplete && !shouldFetchComplete {
				shouldShow = false
			}
		}
		//of those can would be shown based on their state, we hide the ones that are under a collapsed parent
		if shouldShow {
			for _, pathItem := range path.Ancestry {
				ancestorId := pathItem.Id
				backToInt, _ := bullet_stl.AasciBulletIdToInt(ancestorId) //so many conversions. VX:TODO just create a 2 way map or whatever. Maybe that map is its own type.
				if parentGid != nil && backToInt == int64(parentGid.IntValue) {
					continue //we don't care if the parent of the request is collapsed, because they've called for it
				}
				ancestorSummary, ok := summaries[engine.SummaryId(backToInt)]
				if !ok {
					return nil, errors.New("Missing summary for ancestor")
				}
				if ancestorSummary.Flags != nil && ancestorSummary.Flags["collapsed"] {
					fmt.Printf("VX: we are hiding this node beacuse its parent is collapsed %d\n", backToInt)
					shouldShow = false
				}

			}
		}

		//finally, if this is a descendant that we should show, we add it.
		if shouldShow {
			displayItem, err := itemDisplay(summary, now, *gotId, v, alias, path, hasLongForm)
			if err != nil {
				return nil, err
			}
			itemDisplays = append(itemDisplays, *displayItem)
		}
	}
	var sorted []engine.GotItemDisplay
	if sortByPath {
		sorted = SortTheseIntoDFS(itemDisplays)

	} else {
		sorted = SortByUpdated(itemDisplays)
	}
	return e.renderSummaries(sorted, parentItemDisplay)
}
