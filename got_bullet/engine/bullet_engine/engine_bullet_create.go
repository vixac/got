package bullet_engine

import (
	"errors"
	"time"

	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/ids"
	"vixac.com/got/engine"
)

func (e *EngineBullet) CreateBuck(request engine.CreateBuckRequest) (*engine.GotId, error) {

	//lookup parent first because if you're looking up lastId, the lastId will change half way through this func
	var parentGotId *engine.GotId = nil
	override := request.OverrideSettings != nil
	if request.GidLookupInput != nil { //last Id symbol
		var parentLookup *engine.GidLookup = &engine.GidLookup{Input: *request.GidLookupInput}
		fetchedParent, err := e.GidLookup.InputToGid(parentLookup)

		if err != nil {
			return nil, err
		}
		if fetchedParent == nil {
			return nil, errors.New("could not find parent")
		}
		parentGotId = fetchedParent

	}

	var newId int32
	if override && request.OverrideSettings.OverrideId != nil {
		//we need the lastId to be kept up to date with the ids being thrown intot he system. we want the highest id to be set as the last Id.
		newId = int32(*request.OverrideSettings.OverrideId)
		err := e.IgGenerator.SetLastIdIfLower(int64(newId))
		if err != nil {
			return nil, err
		}
	} else {

		newIdFromNext, err := e.IgGenerator.NextId()

		if int64(int32(newIdFromNext)) != newIdFromNext {
			return nil, errors.New("Error. We appear to have ran out of int32 id space.")
		}

		if err != nil {
			return nil, err
		}
		newId = int32(newIdFromNext)
	}

	stringId, err := bullet_stl.BulletIdIntToAasci(int64(newId))
	if err != nil {
		return nil, err
	}
	gotId := engine.GotId{
		AasciValue: stringId,
		IntValue:   newId,
	}

	var deadline *engine.DateTime = nil
	if override && request.OverrideSettings.ScheduleDate != nil {
		deadline = request.OverrideSettings.ScheduleDate
	} else if request.ScheduleLookupInput != nil {

		dateTime, err := engine.NewDeadlineFromDateLookup(*request.ScheduleLookupInput, time.Now())
		if err != nil {
			return nil, err
		}
		deadline = &dateTime
	}

	//add item to ancestry
	ancestry, err := e.AncestorList.AddItem(gotId, parentGotId)
	if err != nil {
		return nil, err
	}

	var headingToStore = request.Heading
	// if the heading is a valid alias, we just create the alias
	// and dont add it as a heading.
	if override && request.OverrideSettings.Alias != nil {
		_, err := e.AliasStore.Alias(gotId, *request.OverrideSettings.Alias)
		if err != nil {
			return nil, err
		}
	} else if override && request.OverrideSettings.NoAlias == true {
		//do nothing. This buck explicitlyu has no radiu
	} else if engine.IsValidAlias(headingToStore) {
		_, err := e.AliasStore.Alias(gotId, headingToStore)
		if err != nil {
			return nil, err
		}
	}

	//add item heading to depot
	err = e.TitleStore.UpsertItem(newId, headingToStore)
	if err != nil {
		return nil, err
	}

	//VX:TODO the shape of blocks has changed. the longform database is wrong.
	//if longform is present in the override, add that too.
	if request.OverrideSettings != nil && request.OverrideSettings.LongForm != nil {
		for _, b := range request.OverrideSettings.LongForm.Blocks {

			err = e.LongFormStore.InsertBlock(b)
			if err != nil {
				return nil, err
			}
		}
	}

	var summaryIds []engine.SummaryId
	if ancestry != nil {
		for _, a := range ancestry.Ids {
			summaryIds = append(summaryIds, engine.SummaryId(a.IntValue))
		}
	}

	var newState engine.GotState = request.InitialState

	e.publishAddEvent(AddItemEvent{
		Id:               engine.SummaryId(newId),
		State:            newState,
		Ancestry:         summaryIds,
		Deadline:         deadline,
		OverrideSettings: request.OverrideSettings,
	})
	return &engine.GotId{
		AasciValue: "0" + stringId,
		IntValue:   int32(newId),
	}, nil
}
