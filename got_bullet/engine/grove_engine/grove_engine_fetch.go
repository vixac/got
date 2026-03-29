package grove_engine

import (
	"errors"
	"fmt"

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
	fmt.Printf("VX: paths %+v \n", paths)
	infos, err := g.InfoStore.InfoForMany(descendantPlusParentIds)
	if err != nil {
		return nil, err
	}
	fmt.Printf("VX: infos %+v \n", infos)

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

	aggregates, err := g.GroveStore.AggregatesForMany(descendantPlusParentIds)
	if err != nil {
		return nil, err
	}
	fmt.Printf("VX: theres the aggs %+v\n", aggregates)

	displays := make(map[engine.GotId]engine.GotItemDisplay)
	//add a display node for each id
	for _, id := range descendantPlusParentIds {
		if id == *parentGid {
			//VX:TODO this is the parent node. We want to render it but not WITH the others.
		}
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

		displays[id] = engine.GotItemDisplay{
			GotId:      id,
			DisplayGid: "0" + id.AasciValue,
			Path:       thePath,
			Title:      info.Title,
			Alias:      theAlias,
			Deadline:   info.Deadline.Date,
		}
	}

	return nil, nil

}
