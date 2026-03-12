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
	if request.GidLookupInput != nil { //last Id symbol
		parent := engine.GidLookup{Input: *request.GidLookupInput}
		fetchedParent, err := e.GidLookup.InputToGid(&parent)

		if err != nil {
			return nil, err
		}
		if fetchedParent == nil {
			return nil, errors.New("could not find parent")
		}
		parentGotId = fetchedParent
	}

	newId, err := e.IgGenerator.NextId()
	newId32 := int32(newId)
	if int64(newId32) != newId {
		return nil, errors.New("Error. We appear to have ran out of int32 id space.")
	}

	if err != nil {
		return nil, err
	}
	stringId, err := bullet_stl.BulletIdIntToaasci(int64(newId))
	if err != nil {
		return nil, err
	}
	gotId := engine.GotId{
		AasciValue: stringId,
		IntValue:   newId32,
	}

	var deadline *engine.DateTime = nil
	if request.OverrideSettings != nil && request.OverrideSettings.ScheduleDate != nil {
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

	// if the heading is a valid alias, we just create the alias
	// and dont add it as a heading.

	var headingToStore = request.Heading
	if engine.IsValidAlias(headingToStore) {
		_, err := e.AliasStore.Alias(gotId, headingToStore)
		if err != nil {
			return nil, err
		}
	}

	//add item heading to depot
	err = e.TitleStore.UpsertItem(newId32, headingToStore)
	if err != nil {
		return nil, err
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
