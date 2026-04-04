package bullet_engine

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	"vixac.com/got/console"
	"vixac.com/got/engine"
	"vixac.com/got/engine/engine_util"
)

const (
	idGenBucket    int32 = 100
	aliasBucket    int32 = 1001
	nodeBucket     int32 = 1002
	ancestorBucket int32 = 1003
	numberGoBucket int32 = 1006
	titleBucket    int32 = 0 //backwards compatability
	longFormBucket int32 = 1005
)

const (
	aggregateNamespace int32 = 2000
	lastIdSymbol             = "0"
)

type EngineBullet struct {
	Client        bullet_interface.BulletClientInterface
	AliasStore    engine_util.AliasStoreInterface
	TitleStore    engine_util.TitleStoreInterface
	GidLookup     engine.GidLookupInterface
	NumberGoStore engine.NumberGoStoreInterface
	SummaryStore  engine.SummaryStoreInterface
	LongFormStore engine.LongFormStoreInterface
	IgGenerator   engine.IdGeneratorInterface

	EventListeners []engine.EventListenerInterface //these will listen to events broadcasted by engineBullet

	AncestorList AncestorListInterface
}

type IdAncestorPair struct {
	Id       engine.GotId
	Ancestry AncestorLookupResult
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

func itemDisplay(summary engine.Summary, now time.Time, gid engine.GotId, title string, alias string, path *engine.GotPath, hasToNote bool) (*engine.GotItemDisplay, error) {

	displayDeadline, deadlineToken, err := deadline(summary, now)
	if err != nil {
		return nil, err
	}

	createdStr, err := createdDateDisplayString(summary, now)
	if err != nil {
		return nil, err
	}
	updatedStr, err := updatedDateDisplayString(summary)
	if err != nil {
		return nil, err
	}
	return &engine.GotItemDisplay{
		GotId:         gid,
		DisplayGid:    "0" + gid.AasciValue,
		Title:         title,
		Path:          path,
		Alias:         alias,
		SummaryObj:    &summary,
		HasTNote:      hasToNote,
		Deadline:      displayDeadline,
		DeadlineToken: deadlineToken,
		Created:       createdStr,
		Updated:       updatedStr,
	}, nil
}

func deadline(summary engine.Summary, now time.Time) (string, console.Token, error) {

	var displayDeadline = ""
	var deadlineToken console.Token = console.TokenSecondary{}
	//VX:TODO get this date wrangling out. Its business logic	//if theres a deadline and either its a group or its an active job
	if summary.Deadline != nil && (summary.State == nil || (summary.State != nil && *summary.State == engine.Active)) {

		//this "n" is not strongly typed and I feel bad.
		//handle all the special cases
		if summary.Deadline.Special == "n" {
			return "---Now---", console.TokenNow{}, nil
		}

		//if its not special, its assumed to be a normal deadline

		deadlineDate, err := summary.Deadline.ToDate()
		if err != nil {
			return "", deadlineToken, err
		}
		if deadlineDate == nil {
			return "", deadlineToken, errors.New("Missing deadline date.")
		}

		deadStr, spaceTime := console.HumanizeDate(time.Time(*deadlineDate), now)
		displayDeadline = deadStr
		deadlineToken = ToToken(spaceTime)
		return displayDeadline, deadlineToken, nil
	}
	return "", deadlineToken, nil
}

func updatedDateDisplayString(summary engine.Summary) (string, error) {
	updatedStr, err := summary.UpdatedDate.JsonDateToReadable()
	if err != nil {
		return "", err
	}
	return updatedStr, nil

}

// VX:TODO move these functions out into a mapping from data to GotItemDisplay. One big func probably
func createdDateDisplayString(summary engine.Summary, now time.Time) (string, error) {
	createdDate, err := summary.CreatedDate.ToDate()
	if err != nil {
		return "", err
	}
	var createdStr = ""
	if createdDate != nil {
		createdStr, _ = console.HumanizeDate(time.Time(*createdDate), now)
	}
	return createdStr, nil
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

func (e *EngineBullet) performUpdateState(gid *engine.GotId, newState engine.GotState, ancestry *AncestorLookupResult) error {
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

	var summaryIds []engine.SummaryId
	if ancestry != nil {
		for _, id := range ancestry.Ids {
			summaryIds = append(summaryIds, engine.SummaryId(id.IntValue))

		}
	}
	event := engine.StateChangeEvent{
		Id:       summaryId,
		OldState: *oldState,
		NewState: &newState,
		Ancestry: summaryIds,
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
	_ = e.publishEditEvent(engine.EditItemEvent{Id: engine.SummaryId(gid.IntValue)})
	return e.TitleStore.UpsertItem(gid.IntValue, newHeading)
}

func (e *EngineBullet) fetchAndDepthSortAncestry(gids []engine.GotId) ([]IdAncestorPair, error) {
	ancestors, err := e.AncestorList.FetchAncestorsOfMany(gids)
	if err != nil || ancestors == nil {
		return nil, err
	}
	idMap := ancestors.Ids
	var sortablePairs []IdAncestorPair

	for id, ancestorResult := range idMap {
		sortablePairs = append(sortablePairs, IdAncestorPair{
			Id:       id,
			Ancestry: ancestorResult,
		})
	}

	//sorted for leaf nodes first.
	sort.Slice(sortablePairs, func(i, j int) bool {
		return len(sortablePairs[i].Ancestry.Ids) > len(sortablePairs[j].Ancestry.Ids)
	})
	return sortablePairs, nil
}

// VX:TODO move somewhere global

// resolves all gidlookups into gotids and then sorts them to deepest first.
func (e *EngineBullet) ResolveBulkLookupsReverseDepthSorted(lookups []engine.GidLookup) ([]IdAncestorPair, error) {
	var gids []engine.GotId
	for _, lookup := range lookups {
		gid, err := e.GidLookup.InputToGid(&lookup)
		if err != nil || gid == nil {
			return nil, err
		}
		gids = append(gids, *gid)
	}
	sortedPairs, err := e.fetchAndDepthSortAncestry(gids)
	if err != nil {
		return nil, err
	}
	return sortedPairs, nil
}

func (e *EngineBullet) MarkResolved(lookups []engine.GidLookup) error {

	sortedPairs, err := e.ResolveBulkLookupsReverseDepthSorted(lookups)
	if err != nil {
		return err
	}
	complete := engine.GotState(engine.Complete)
	for i, pair := range sortedPairs {
		err := e.performUpdateState(&pair.Id, complete, &pair.Ancestry)
		if err != nil {
			fmt.Printf("Warn: did not complete updating state. Only reached item %d of %d", i, len(sortedPairs))
			return err
		}
	}
	return nil
}

func (e *EngineBullet) Move(lookup engine.GidLookup, newParent engine.GidLookup) error {
	gid, err := e.GidLookup.InputToGid(&lookup)
	if err != nil || gid == nil {
		return err
	}
	parent, err := e.GidLookup.InputToGid(&newParent)
	if err != nil {
		return err
	}
	moveRes, err := e.AncestorList.MoveItem(*gid, parent)
	if err != nil {
		return err
	}
	if moveRes == nil {
		return errors.New("move returned nil result")
	}

	// Convert old ancestry to SummaryIds
	var oldAncestry []engine.SummaryId
	if moveRes.OldAncestry != nil {
		for _, id := range moveRes.OldAncestry.Ids {
			oldAncestry = append(oldAncestry, engine.SummaryId(id.IntValue))
		}
	}

	// Convert new ancestry to SummaryIds
	var newAncestry []engine.SummaryId
	if moveRes.NewAncestry != nil {
		for _, id := range moveRes.NewAncestry.Ids {
			newAncestry = append(newAncestry, engine.SummaryId(id.IntValue))
		}
	}

	e.publishMoveEvent(engine.ItemMovedEvent{
		Id:          engine.SummaryId(gid.IntValue),
		OldAncestry: oldAncestry,
		NewAncestry: newAncestry,
	})

	return nil
}

func (e *EngineBullet) publishMoveEvent(event engine.ItemMovedEvent) error {
	for _, l := range e.EventListeners {
		err := l.ItemMoved(event)
		if err != nil {
			fmt.Printf("VX: Listner error was %s\n", err.Error())
			fmt.Printf("VX:TODO listener had an error and I dont think it shoudl stop anything so I'm ignoring it")
		}
	}
	return nil
}

func (e *EngineBullet) publishAddEvent(event engine.AddItemEvent) error {
	for _, l := range e.EventListeners {
		err := l.ItemAdded(event)
		if err != nil {
			fmt.Printf("VX: Listner error was %s\n", err.Error())
			fmt.Printf("VX:TODO listener had an error and I dont think it shoudl stop anything so I'm ignoring it")
		}
	}
	return nil
}

func (e *EngineBullet) publishItemDeletedEvent(event engine.ItemDeletedEvent) error {
	for _, l := range e.EventListeners {
		err := l.ItemDeleted(event)
		if err != nil {
			fmt.Printf("VX:state change  Listner error was %s\n", err.Error())
			fmt.Printf("VX:TODO listener had an error and I dont think it shoudl stop anything so I'm ignoring it")
		}
	}
	return nil
}

func (e *EngineBullet) publishEditEvent(event engine.EditItemEvent) error {
	for _, l := range e.EventListeners {
		err := l.ItemEdited(event)
		if err != nil {
			fmt.Printf("VX: Listner error was %s\n", err.Error())
			fmt.Printf("VX:TODO listener had an error and I dont think it shoudl stop anything so I'm ignoring it")
		}
	}
	return nil
}

func (e *EngineBullet) publishStateChangeEvent(event engine.StateChangeEvent) error {
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

func (e *EngineBullet) Alias(lookup engine.GidLookup, alias string) (bool, error) {

	gid, err := e.GidLookup.InputToGid(&lookup)
	if err != nil || gid == nil {
		return false, err
	}

	//confirm the gid exists.
	return e.AliasStore.Alias(*gid, alias)
}

// VX:TODO this is used in Summary, but can be deleted and replaced with  ancestorPathFor
func (e *EngineBullet) ancestorPathFrom(ancestors *AncestorLookupResult) (*engine.GotPath, error) {
	if ancestors == nil {
		return nil, nil
	}
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

// VX:TODO that "anc" prefix is totally unneeded because you have a bucket.
func NewEngineBullet(client bullet_interface.BulletClientInterface) (engine.GotEngineInterface, error) {
	ancestorList, err := NewAncestorList(client, "anc", ancestorBucket, ":", ">", "<")

	if err != nil {
		return nil, err
	}
	codec := &engine_util.JSONCodec[engine.Summary]{}
	aggStore, err := engine_util.NewBulletSummaryStore(aggregateNamespace, client, client, codec)
	if err != nil {
		return nil, err
	}

	titleStore := engine_util.NewBulletTitleStore(titleBucket, client, client)

	longFormStore, err := engine_util.NewBulletLongFormStore(longFormBucket, client, client)
	if err != nil {
		return nil, err
	}

	numberGoCodec := &engine_util.JSONCodec[engine_util.NumberGoBlock]{}
	numberGoStore, err := engine_util.NewBulletNumberGoStore(numberGoBucket, client, client, numberGoCodec)
	if err != nil {
		return nil, err
	}

	aliasStore, err := engine_util.NewBulletAliasStore(client, aliasBucket)
	if err != nil {
		return nil, err
	}
	idGenerator := engine_util.NewIdBulletGenerator(client, idGenBucket, "next-id-list", "", "latest")

	gidLookup, err := engine_util.NewBulletGidLookup(aliasStore, numberGoStore, idGenerator)
	if err != nil {
		return nil, err
	}

	var listeners []engine.EventListenerInterface
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
		IgGenerator:    idGenerator,
		EventListeners: listeners,
	}, nil
}
