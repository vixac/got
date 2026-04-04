package grove_engine

import (
	"errors"
	"sort"
	"strconv"
	"time"

	"vixac.com/got/engine"
	"vixac.com/got/engine/engine_util"
)

func (g *GroveEngine) fetchAndDepthSortAncestry(gids []engine.GotId) ([]GotIdWithPath, error) {
	pairs, err := g.GroveStore.FetchAncestorsForMany(gids)

	if err != nil {
		return nil, err
	}
	//sorted for leaf nodes first.
	sort.Slice(pairs, func(i, j int) bool {
		return len(pairs[i].Path) > len(pairs[j].Path)
	})
	return pairs, nil
}

// resolves all gidlookups into gotids and then sorts them to deepest first.
func (g *GroveEngine) ResolveBulkLookupsReverseDepthSorted(lookups []engine.GidLookup) ([]GotIdWithPath, error) {
	var gids []engine.GotId
	for _, lookup := range lookups {
		gid, err := g.GidLookup.InputToGid(&lookup)
		if err != nil || gid == nil {
			return nil, err
		}
		gids = append(gids, *gid)
	}
	sortedPairs, err := g.fetchAndDepthSortAncestry(gids)
	if err != nil {
		return nil, err
	}
	return sortedPairs, nil
}

func (g *GroveEngine) MarkResolved(lookups []engine.GidLookup) error {
	sortedPairs, err := g.ResolveBulkLookupsReverseDepthSorted(lookups)
	if err != nil {
		return err
	}
	complete := engine.GotState(engine.Complete)
	now := time.Now()
	millis := now.UnixNano()
	millisStr := strconv.FormatInt(millis, 10)
	var ids []engine.GotId
	for _, pair := range sortedPairs {
		ids = append(ids, pair.Id)
	}

	stateMap, err := g.GroveStore.IndividualStateForMany(ids)
	if err != nil {
		return err
	}
	var statePairs []GotIdWithState
	for _, id := range ids {
		state, ok := stateMap[id]
		if !ok {
			return errors.New("Missing state.")
		}
		statePairs = append(statePairs, GotIdWithState{
			Id:    id,
			State: state,
		})
	}
	return g.GroveStore.BulkChangeState(statePairs, complete, millisStr)
}

// grab the info for this lookup, assuming it maps to 1 info item, and update the timeedited on the info object before returning it (this doesnt save anything)
func (g *GroveEngine) updatedEditTimeInfoForLookup(lookup engine.GidLookup) (*engine_util.BuckInfo, *engine.GotId, error) {
	gid, err := g.GidLookup.InputToGid(&lookup)
	if err != nil || gid == nil {
		return nil, nil, err
	}

	infoMap, err := g.InfoStore.InfoForMany([]engine.GotId{*gid})
	if err != nil {
		return nil, nil, err
	}
	info, ok := infoMap.InfoMap[*gid]
	if !ok {
		return nil, nil, errors.New("missing info for gid")
	}
	//apply the date to the edited field of this
	now, err := engine.NewDateTime(time.Now())
	if err != nil {
		return nil, nil, err
	}
	info.UpdatedDate = now

	return &info, gid, nil

}

func (g *GroveEngine) EditTitle(lookup engine.GidLookup, newHeading string) error {

	info, gid, err := g.updatedEditTimeInfoForLookup(lookup)
	if err != nil {
		return err
	}
	info.Title = newHeading
	return g.InfoStore.UpsertInfo(*gid, *info)
}

func (g *GroveEngine) ScheduleItem(lookup engine.GidLookup, dateLookup engine.DateLookup) error {
	info, gid, err := g.updatedEditTimeInfoForLookup(lookup)
	if err != nil {
		return err
	}

	deadline, err := engine.NewDeadlineFromDateLookup(dateLookup.UserInput, time.Now())
	if err != nil {
		return err
	}
	info.Deadline = &deadline
	return g.InfoStore.UpsertInfo(*gid, *info)
}
func (g *GroveEngine) TagItem(lookup engine.GidLookup, tagLookup engine.TagLookup) error {

	info, gid, err := g.updatedEditTimeInfoForLookup(lookup)
	if err != nil {
		return err
	}

	tagLiteralString := tagLookup.Input
	tagLiteral := engine.TagLiteral{
		Display: tagLiteralString,
		Token:   "", //not fucking with token yet.
	}

	tag := engine.Tag{
		Literal:    &tagLiteral,
		Identifier: nil,
	}

	//ignore if duplicate
	for _, t := range info.Tags {
		if t.EqualTo(tag) {
			return nil
		}
	}
	info.Tags = append(info.Tags, tag)
	return g.InfoStore.UpsertInfo(*gid, *info)
}

func (g *GroveEngine) ToggleCollapse(lookup engine.GidLookup, collapsed bool) error {
	info, gid, err := g.updatedEditTimeInfoForLookup(lookup)
	if err != nil {
		return err
	}

	if info.Flags == nil {
		info.Flags = make(map[string]bool)
	}
	if collapsed {
		info.Flags["collapsed"] = true
	} else {
		delete(info.Flags, "collapsed")
	}
	return g.InfoStore.UpsertInfo(*gid, *info)
}
