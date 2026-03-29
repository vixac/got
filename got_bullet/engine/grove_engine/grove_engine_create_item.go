package grove_engine

import (
	"errors"
	"time"

	"vixac.com/got/engine"
	"vixac.com/got/engine/engine_util"
)

// This is not idempodent, and if it errors out part of the way, then you're in crazy town.
func (g *GroveEngine) CreateBuck(request engine.CreateBuckRequest) (*engine.GotId, error) {
	//grab the parent, which might be nil
	var lookup *engine.GidLookup = nil
	if request.GidLookupInput != nil {
		lookup = &engine.GidLookup{Input: *request.GidLookupInput}
	}
	parentGotId, err := g.deriveLookup(lookup)
	if err != nil {
		return nil, err
	}
	newId, err := g.claimNewId(request.OverrideSettings)
	if err != nil || newId == nil {
		return nil, err
	}

	err = g.aliasBuckIfNeeded(request, *newId)
	if err != nil {
		return nil, err
	}

	err = g.maybeWriteLongforms(request.OverrideSettings)
	if err != nil {
		return nil, err
	}

	err = g.writeBuckInfo(request, *newId)
	if err != nil {
		return nil, err
	}

	groveReq := GotStoreCreateRequest{
		Id:     *newId,
		State:  request.InitialState,
		Parent: parentGotId,
	}
	err = g.GroveStore.CreateBuck(groveReq)
	if err != nil {
		return nil, err
	}
	return newId, nil
}

func (g *GroveEngine) writeBuckInfo(request engine.CreateBuckRequest, id engine.GotId) error {
	now, _ := engine.NewDateTime(time.Now())
	created, err := createdTimeOrNil(request.OverrideSettings)
	if err != nil {
		return err
	}
	updated, err := updatedTimeOrNil(request.OverrideSettings)
	if err != nil {
		return err
	}
	deadline, err := deriveDeadline(request)
	if err != nil {
		return err
	}
	if created == nil {
		created = &now
	}
	if updated == nil {
		updated = &now
	}

	var tags []engine.Tag
	if request.OverrideSettings != nil {
		tags = request.OverrideSettings.Tags
	}
	flags := make(map[string]bool)
	if request.OverrideSettings != nil {
		for _, f := range request.OverrideSettings.Flags {
			flags[f] = true
		}
	}

	buckInfo := engine_util.NewBuckInfo(request.Heading, deadline, *created, *updated, tags, flags)
	err = g.InfoStore.UpsertInfo(id, buckInfo)
	return err
}

func (g *GroveEngine) maybeWriteLongforms(override *engine.CreateOverrideSettings) error {
	if override != nil && override.LongForm != nil && len(override.LongForm) != 0 {

		for _, restoreBlock := range override.LongForm {
			blockId, err := engine.NewLongFormKeyFromString(restoreBlock.KeyString)
			if err != nil {
				return err
			}
			editTime, err := engine.EpochMillisStringToDate(restoreBlock.EditMillis)
			if err != nil {
				return err
			}
			block := engine.LongFormBlock{
				Id:      *blockId,
				Content: restoreBlock.Content,
				Edited:  *editTime,
			}
			err = g.LongFormStore.InsertBlock(block)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *GroveEngine) aliasBuckIfNeeded(request engine.CreateBuckRequest, id engine.GotId) error {
	if request.HasOverride() && request.OverrideSettings.Alias != nil {
		_, err := g.AliasStore.Alias(id, *request.OverrideSettings.Alias)
		if err != nil {
			return err
		}
	} else if request.HasOverride() && request.OverrideSettings.NoAlias == true {
		//do nothing. This buck explicitlyu has no radiu
	} else if engine_util.IsValidAlias(request.Heading) {
		_, err := g.AliasStore.Alias(id, request.Heading)
		if err != nil {
			return err
		}
	}
	return nil
}
func deriveDeadline(request engine.CreateBuckRequest) (*engine.DateTime, error) {

	if request.HasOverride() && request.OverrideSettings.ScheduleDate != nil {
		return request.OverrideSettings.ScheduleDate, nil
	} else if request.ScheduleLookupInput != nil {
		dateTime, err := engine.NewDeadlineFromDateLookup(*request.ScheduleLookupInput, time.Now())
		if err != nil {
			return nil, err
		}
		return &dateTime, nil
	}
	return nil, nil
}

// creates a new id and sets it as the last (if its the last)
func (g *GroveEngine) claimNewId(override *engine.CreateOverrideSettings) (*engine.GotId, error) {

	var newId int32
	if override != nil && override.OverrideId != nil {
		//we need the lastId to be kept up to date with the ids being thrown intot he system. we want the highest id to be set as the last Id.
		newId = int32(*override.OverrideId)
		err := g.IdGenerator.SetLastIdIfLower(int64(newId))
		if err != nil {
			return nil, err
		}

	} else {
		newIdFromNext, err := g.IdGenerator.NextId()
		if int64(int32(newIdFromNext)) != newIdFromNext {
			return nil, errors.New("Error. We appear to have ran out of int32 id space.")
		}

		if err != nil {
			return nil, err
		}
		newId = int32(newIdFromNext)
	}
	return engine.NewGotIdFromInt(newId)
}

// if the lookup is nil, return nil.
// if the lookup was entered, return the gid or error
func (g *GroveEngine) deriveLookup(lookup *engine.GidLookup) (*engine.GotId, error) {
	if lookup == nil {
		return nil, nil
	}
	fetched, err := g.GidLookup.InputToGid(lookup)

	if err != nil {
		return nil, err
	}
	if fetched == nil {
		return nil, errors.New("could not find parent")
	}
	return fetched, nil

}

func stringToTime(dateString string) (engine.DateTime, error) {
	createdDate, err := engine.NewTimeFromString(dateString)
	if err != nil {
		return engine.DateTime{}, err
	}
	return engine.NewDateTime(time.Time(*createdDate))

}
func createdTimeOrNil(override *engine.CreateOverrideSettings) (*engine.DateTime, error) {

	if override != nil {
		date, err := stringToTime(override.CreatedDate)
		return &date, err
	}
	return nil, nil
}
func updatedTimeOrNil(override *engine.CreateOverrideSettings) (*engine.DateTime, error) {

	if override != nil {
		date, err := stringToTime(override.UpdatedDate)
		return &date, err
	}
	return nil, nil
}
