package bullet_engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/ids"

	"vixac.com/got/console"
	"vixac.com/got/engine"
)

const (
	aliasBucket    int32 = 1001
	nodeBucket     int32 = 1002
	ancestorBucket int32 = 1003
	numberGoBucket int32 = 1004
	titleBucket    int32 = 0 //backwards compatability
	longFormBucket int32 = 1005
)

const (
	aggregateNamespace int32 = 2000
)

type EngineBullet struct {
	Client        bullet_interface.BulletClientInterface
	AncestorList  AncestorListInterface
	TitleStore    TitleStoreInterface
	GidLookup     GidLookupInterface
	AliasStore    engine.GotAliasInterface
	NumberGoStore NumberGoStoreInterface
	SummaryStore  SummaryStoreInterface
	LongFormStore LongFormStoreInterface

	EventListeners []EventListenerInterface //these will listen to events broadcasted by engineBullet

	//interface conformance
	//	LongFormStoreInterface
}

func (e *EngineBullet) ScheduleItem(lookup engine.GidLookup, dateLookup engine.DateLookup) error {
	gid, err := e.GidLookup.InputToGid(&lookup)
	if err != nil || gid == nil {
		return err
	}

	summaryId := engine.SummaryId(gid.IntValue)

	ids := []engine.SummaryId{summaryId}
	items, err := e.SummaryStore.Fetch(ids)
	if err != nil {
		return err
	}
	summary, ok := items[summaryId]
	if !ok {
		return errors.New("missing summary")
	}

	deadline, err := engine.NewDeadlineFromDateLookup(dateLookup.UserInput, time.Now())
	if err != nil {
		return err
	}
	summary.Deadline = &deadline
	return e.SummaryStore.UpsertSummary(summaryId, summary)
}

// lets rewrite this maybe.
func (e *EngineBullet) FetchItemsBelow(lookup *engine.GidLookup, descendantType int, states []int) (*engine.GotFetchResult, error) {

	//0->1 numberstore gid -> numberstore
	//0-> alias store gid -> alias store
	gid, err := e.GidLookup.InputToGid(lookup)
	if err != nil || gid == nil {
		return nil, err
	}

	//1.gid->ancestor (object -> subject)
	//2.all descendants: allpairs for full key
	all, err := e.AncestorList.FetchImmediatelyUnder(*gid)
	if err != nil || all == nil {
		return nil, err
	}

	//get string ids of all items to do the alias lookup
	stringIds := make([]string, len(all.Ids))

	i := 0
	for k := range all.Ids {
		stringIds[i] = k
		i++
	}

	aliasMap, err := e.AliasStore.LookupAliasForMany(stringIds)
	if err != nil {
		return nil, err
	}

	var intIds []int32
	ancestorPaths := make(map[int32]engine.GotPath)
	for id, ancestorLookup := range all.Ids {

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
	//titleStore: allIds -> title
	titles, err := e.TitleStore.TitleForMany(intIds)
	if err != nil {
		return nil, err
	}

	var summaryIds []engine.SummaryId
	for _, v := range intIds {
		summaryIds = append(summaryIds, engine.SummaryId(v))
	}
	summaries, err := e.SummaryStore.Fetch(summaryIds)
	if err != nil {
		return nil, err
	}

	//just needed to see if we present the note emoji. Unfortunately we're loading
	//the actual notes on here.
	longForms, err := e.LongFormForMany(intIds)
	if err != nil {
		return nil, err
	}

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

		gotId, err := engine.NewGotId(stringId)
		if err != nil {
			return nil, err
		}
		summaryId := NewSummaryId(*gotId)
		summary, ok := summaries[summaryId]
		var summaryPointer *engine.Summary = nil
		if ok {
			summaryPointer = &summary
		} else {
			return nil, errors.New("missing summary in fetchItems Below")
		}

		//here we filter complete leafs from the jobs list, and their notes.
		//VX:Note we want to have completes
		//not even appear in the search, because thats more scalable.
		isComplete := summary.State != nil && *summary.State == engine.Complete
		isNote := summary.State != nil && *summary.State == engine.Note

		isHiddenNote := isNote && isParentComplete
		_, hasLongForm := longForms[k]

		var displayDeadline = ""
		var deadlineToken console.Token = console.TokenSecondary{}

		//if theres a deadline and either its a group or its an active job
		if summary.Deadline != nil && (summary.State == nil || (summary.State != nil && *summary.State == engine.Active)) {

			var date console.RFC3339Time
			dateBytes := []byte(summary.Deadline.Date)
			err := json.Unmarshal(dateBytes, &date)
			if err != nil {
				return nil, err
			}

			deadStr, spaceTime := console.HumanizeDate(time.Time(date), time.Now())
			displayDeadline = deadStr
			deadlineToken = ToToken(spaceTime)
		}
		if !isComplete && !isHiddenNote {
			itemDisplays = append(itemDisplays, engine.GotItemDisplay{
				Gid:           stringId,
				Title:         v,
				Path:          path,
				Alias:         alias,
				SummaryObj:    summaryPointer,
				HasTNote:      hasLongForm,
				Deadline:      displayDeadline,
				DeadlineToken: deadlineToken,
			})
		}

	}
	sorted := SortTheseIntoDFS(itemDisplays)
	return e.renderSummaries(sorted)
}

func ToToken(s console.SpaceTime) console.Token {
	switch s.TimeType {
	case console.PastMany:
		return console.TokenAlert{}
	case console.FutureMany:
		return console.TokenBrand{}
	default:
		return console.TokenWarning{}

	}
}

func (e *EngineBullet) Summary(lookup *engine.GidLookup) (*engine.GotItemDisplay, error) {

	gid, err := e.GidLookup.InputToGid(lookup)
	if err != nil {
		return nil, err
	}
	if gid == nil {
		return nil, errors.New("no gid")
	}

	title, err := e.TitleStore.TitleFor(gid.IntValue)
	if err != nil {
		return nil, errors.New("no title")
	}
	if title == nil {
		return nil, nil
	}

	ancestorResult, err := e.AncestorList.FetchAncestorsOf(*gid)
	if err != nil {
		return nil, errors.New("error fetching ancestors")
	}
	path, err := e.ancestorPathFrom(ancestorResult)
	if err != nil {
		return nil, err
	}

	return &engine.GotItemDisplay{
		Gid:   gid.AasciValue,
		Title: *title,
		Path:  path,
	}, nil

}

// adds the items to the number go store as well as
func (e *EngineBullet) renderSummaries(summaries []engine.GotItemDisplay) (*engine.GotFetchResult, error) {

	var expandedSummaries []engine.GotItemDisplay
	var pairs []NumberGoPair
	for i, s := range summaries {

		num := i + 1
		pairs = append(pairs, NumberGoPair{
			Number: num,
			Gid:    engine.Gid{Id: s.Gid},
		})
		expandedSummaries = append(expandedSummaries, engine.GotItemDisplay{
			Gid:           "0" + s.Gid, //VX:TODO here is the "0 prefix on the gid."
			Alias:         s.Alias,
			NumberGo:      num,
			Title:         s.Title,
			Path:          s.Path,
			SummaryObj:    s.SummaryObj,
			HasTNote:      s.HasTNote,
			Deadline:      s.Deadline,
			DeadlineToken: s.DeadlineToken,
		})
	}

	err := e.NumberGoStore.AssignNumberPairs(pairs)
	if err != nil {
		return nil, err
	}

	//the summaries injected dont have number go assigned.
	res := engine.GotFetchResult{Result: expandedSummaries}
	return &res, nil
}

func (e *EngineBullet) MarkActive(lookup engine.GidLookup) (*engine.NodeId, error) {
	var newState engine.GotState = engine.Active
	return nil, e.updateState(lookup, newState)

}

func (e *EngineBullet) MarkAsNote(lookup engine.GidLookup) (*engine.NodeId, error) {
	var newState engine.GotState = engine.Note
	return nil, e.updateState(lookup, newState)
}

func (e *EngineBullet) updateState(lookup engine.GidLookup, newState engine.GotState) error {
	gid, err := e.GidLookup.InputToGid(&lookup)
	if err != nil {
		return err
	}
	if gid == nil {
		return nil
	}
	summaryId := engine.SummaryId(gid.IntValue)
	ids := []engine.SummaryId{summaryId}
	res, err := e.SummaryStore.Fetch(ids)
	if err != nil {
		return err
	}
	if res == nil {
		return errors.New("missing summary")
	}
	summary, ok := res[summaryId]
	if !ok {
		return errors.New("no summary for this id")
	}
	oldState := summary.State
	if oldState == nil {
		return errors.New("cant resolve an item without a state")
	}

	ancestorResult, err := e.AncestorList.FetchAncestorsOf(*gid)
	if err != nil {
		return errors.New("error fetching ancestors")
	}
	var summaryIds []engine.SummaryId
	for _, id := range ancestorResult.Ids {
		summaryIds = append(summaryIds, engine.SummaryId(id.IntValue))

	}
	thisNode, err := e.Summary(&lookup)
	if err != nil {
		return err
	}
	if thisNode == nil {
		return errors.New("missing summary")
	}

	event := StateChangeEvent{
		Id:       summaryId,
		OldState: *oldState,
		NewState: newState,
		Ancestry: summaryIds, //VX:TODO fetch?
	}
	return e.publishStateChangeEvent(event)
}

func (e *EngineBullet) EditTitle(lookup engine.GidLookup, newHeading string) error {

	gid, err := e.GidLookup.InputToGid(&lookup)
	if err != nil {
		return err
	}
	if gid == nil {
		return nil
	}
	return e.TitleStore.UpsertItem(gid.IntValue, newHeading)
}

func (e *EngineBullet) MarkResolved(lookup []engine.GidLookup) error {
	for _, lookup := range lookup {
		var newState engine.GotState = engine.Complete

		//If one fails, we stop and return that error.
		err := e.updateState(lookup, newState)
		if err != nil {
			return err
		}
	}
	return nil

}

func (e *EngineBullet) Delete(lookup engine.GidLookup) (*engine.NodeId, error) {
	//check if the gid is an exact match for an item id
	//check int32 parse, check its length is the right length

	//aliases can't start with a number.
	return nil, errors.New("not impl")
}

func (e *EngineBullet) Move(lookup engine.GidLookup, newParent engine.GidLookup) (*engine.NodeId, error) {
	return nil, errors.New("not impl")
}

func (e *EngineBullet) CreateBuck(parent *engine.GidLookup, date *engine.DateLookup, completable bool, heading string) (*engine.NodeId, error) {
	newId, err := e.NextId()

	if err != nil {
		return nil, err
	}
	stringId, err := bullet_stl.BulletIdIntToaasci(int64(newId))
	if err != nil {
		return nil, err
	}
	gotId := engine.GotId{
		AasciValue: stringId,
		IntValue:   newId,
	}

	var deadline *engine.Deadline = nil

	if date != nil {
		deadlineTime, err := console.ParseRelativeDate(date.UserInput, time.Now())
		if err != nil {
			return nil, err
		}
		formatted := deadlineTime.Format("Mon 2 Jan 2006")
		fmt.Printf("VX: Deadline date it %s", formatted)
		dateJsonByes, err := deadlineTime.MarshalJSON()
		if err != nil {
			return nil, err
		}
		deadline = &engine.Deadline{Date: string(dateJsonByes)}
	}

	var parentGotId *engine.GotId = nil
	if parent != nil {
		fmt.Printf("Looking up parent %s\n", parent.Input)
		fetchedParent, err := e.GidLookup.InputToGid(parent)

		if err != nil {
			return nil, err
		}
		if fetchedParent == nil {
			return nil, errors.New("could not find parent")
		}
		parentGotId = fetchedParent
	}

	//add item to ancestry
	ancestry, err := e.AncestorList.AddItem(gotId, parentGotId)
	if err != nil {
		return nil, err
	}

	// if the heading is a valid alias, we just create the alias
	// and dont add it as a heading.

	var headingToStore = heading
	if engine.IsValidAlias(heading) {
		//headingToStore = "" //VX:Note I've decided against nulling the title because if you unalias, the meaning of this thing is totally gone.
		_, err := e.AliasStore.Alias(stringId, heading)
		if err != nil {
			return nil, err
		}
	}

	//add item heading to depot
	err = e.TitleStore.UpsertItem(newId, headingToStore)
	if err != nil {
		return nil, err
	}

	var summaryIds []engine.SummaryId
	if ancestry != nil {
		for _, a := range ancestry.Ids {
			summaryIds = append(summaryIds, engine.SummaryId(a.IntValue))
		}
	}

	var newState engine.GotState = engine.Note
	if completable {
		newState = engine.Active
	}
	e.publishAddEvent(AddItemEvent{
		Id:       engine.SummaryId(newId),
		State:    newState,
		Ancestry: summaryIds,
		Deadline: deadline,
	})

	return &engine.NodeId{
		Gid: engine.Gid{
			Id: stringId,
		},
		Title: heading,
		Alias: "",
	}, nil
}
func (e *EngineBullet) publishAddEvent(event AddItemEvent) error {
	for _, l := range e.EventListeners {
		err := l.ItemAdded(event)
		if err != nil {
			fmt.Printf("VX: Listner error was %s\n", err.Error())
			fmt.Printf("VX:TODO listener had an error and I dont think it shoudl stop anything so I'm ignoring it")
		}
	}
	return nil
}

func (e *EngineBullet) publishStateChangeEvent(event StateChangeEvent) error {
	for _, l := range e.EventListeners {
		err := l.ItemStateChanged(event)
		if err != nil {
			fmt.Printf("VX:state change  Listner error was %s\n", err.Error())
			fmt.Printf("VX:TODO listener had an error and I dont think it shoudl stop anything so I'm ignoring it")
		}
	}
	return nil
}

func (e *EngineBullet) Lookup(alias string) (*engine.GotId, error) {
	return e.AliasStore.Lookup(alias)

}
func (e *EngineBullet) Unalias(alias string) (*engine.GotId, error) {
	return e.AliasStore.Unalias(alias)
}

func (e *EngineBullet) LookupAliasForGid(gid string) (*string, error) {
	return e.AliasStore.LookupAliasForGid(gid)

}

func (e *EngineBullet) LookupAliasForMany(gid []string) (map[string]*string, error) {
	return e.AliasStore.LookupAliasForMany(gid)
}

func (e *EngineBullet) Alias(gid string, alias string) (bool, error) {
	//confirm the gid exists.
	lookup, err := e.Summary(&engine.GidLookup{
		Input: gid,
	})
	if err != nil {
		return false, err
	}

	if lookup == nil {
		return false, errors.New("can't alias a gid that doesn't exist")
	}
	return e.AliasStore.Alias(lookup.Gid, alias)
}

// VX:TODO this is used in Summary, but can be deleted and replaced with  ancestorPathFor
func (e *EngineBullet) ancestorPathFrom(ancestors *AncestorLookupResult) (*engine.GotPath, error) {
	var items []engine.PathItem
	//VX:TODO are they sorted by ancestry?
	//I'm confused. I think there is always 1 item in here

	var gids []string
	for _, gid := range ancestors.Ids {
		gids = append(gids, gid.AasciValue)
	}

	res, err := e.AliasStore.LookupAliasForMany(gids)
	if err != nil {
		return nil, nil
	}
	for _, id := range ancestors.Ids {
		var alias *string
		if res != nil { //if there are aliases to inspect.
			matchedAlias, ok := res[id.AasciValue]
			if ok {
				alias = matchedAlias
			}
		}

		items = append(items, engine.PathItem{
			Id:    id.AasciValue,
			Alias: alias,
		})
	}
	return &engine.GotPath{
		Ancestry: items,
	}, nil
}

// VX:TODO use this one, delete ancestorPathFrom
func ancestorPathFor(ancestors *AncestorLookupResult, aliases map[string]*string) *engine.GotPath {
	var items []engine.PathItem
	for _, id := range ancestors.Ids {
		var alias *string
		if aliases != nil { //if there are aliases to inspect.
			matchedAlias, ok := aliases[id.AasciValue]
			if ok {
				alias = matchedAlias
			}
		}

		items = append(items, engine.PathItem{
			Id:    id.AasciValue,
			Alias: alias,
		})
	}
	return &engine.GotPath{
		Ancestry: items,
	}

}

func NewEngineBullet(client bullet_interface.BulletClientInterface) (*EngineBullet, error) {
	ancestorList, err := NewAncestorList(client, "anc", ancestorBucket, ":", ">", "<")

	if err != nil {
		return nil, err
	}
	codec := &JSONCodec[engine.Summary]{}
	aggStore, err := NewBulletSummaryStore(codec, client, aggregateNamespace)
	if err != nil {
		return nil, err
	}

	titleStore, err := NewBulletTitleStore(client, titleBucket)
	if err != nil {
		return nil, err
	}
	longFormStore, err := NewBulletLongFormStore(client, longFormBucket)
	if err != nil {
		return nil, err
	}

	numberGoStore, err := NewBulletNumberGoStore(client, numberGoBucket)
	if err != nil {
		return nil, err
	}

	aliasStore, err := NewBulletAliasStore(client, aliasBucket)
	if err != nil {
		return nil, err
	}

	gidLookup, err := NewBulletGidLookup(aliasStore, numberGoStore)
	if err != nil {
		return nil, err

	}

	var listeners []EventListenerInterface
	aggregator, err := NewAggregator(aggStore)
	if err != nil {
		return nil, err
	}

	listeners = append(listeners, aggregator)

	return &EngineBullet{
		Client:         client,
		AncestorList:   ancestorList,
		TitleStore:     titleStore,
		GidLookup:      gidLookup,
		AliasStore:     aliasStore,
		NumberGoStore:  numberGoStore,
		SummaryStore:   aggStore,
		LongFormStore:  longFormStore,
		EventListeners: listeners,
	}, nil
}
