package grove_engine

import (
	"errors"
	"time"

	"vixac.com/got/engine"
	"vixac.com/got/engine/engine_util"
)

func (g *GroveEngine) MarkResolved(lookup []engine.GidLookup) error {
	/**
	the previous one works by calling  fetchAndDepthSortAncestry and then perform update
	1 at a time. not great.

	*/
	return errors.New(" MarkResolved Not impl")
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
