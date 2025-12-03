package bullet_engine

import (
	"errors"
	"fmt"
	"sort"
	"strconv"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/ids"

	"vixac.com/got/engine"
)

const (
	aliasBucket    int32 = 1001
	nodeBucket     int32 = 1002
	ancestorBucket int32 = 1003
	numberGoBucket int32 = 1004
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

	EventListeners []EventListenerInterface //these will listen to events broadcasted by engineBullet
}

func NewEngineBullet(client bullet_interface.BulletClientInterface) (*EngineBullet, error) {
	ancestorList, err := NewAncestorList(client, "anc", ancestorBucket, ":", ">", "<")

	if err != nil {
		return nil, err
	}
	codec := &JSONCodec[Summary]{}
	aggStore, err := NewBulletSummaryStore(codec, client, aggregateNamespace)
	if err != nil {
		return nil, err
	}

	titleStore, err := NewBulletTitleStore(client)
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
		EventListeners: listeners,
	}, nil
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

func (e *EngineBullet) ancestorPathFrom(ancestors *AncestorLookupResult) (*engine.GotPath, error) {
	var items []engine.PathItem
	//VX:TODO are they sorted by ancestry?
	//I'm confused. I think there is always 1 item in here

	fmt.Printf("VX: There are %d ancestor Ids here\n", len(ancestors.Ids))
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

// lets rewrite this maybe.
func (e *EngineBullet) FetchItemsBelow(lookup *engine.GidLookup, descendantType int, states []int) (*engine.GotFetchResult, error) {
	gid, err := e.GidLookup.InputToGid(lookup)
	if err != nil {
		return nil, err
	}
	if gid == nil {
		return nil, nil
	}
	fmt.Printf("VX: resolved gid is %s\n", gid.AasciValue)
	all, err := e.AncestorList.FetchImmediatelyUnder(*gid)
	if err != nil {
		return nil, err
	}
	if all == nil {
		return nil, nil
	}

	var intIds []int32
	ancestorPaths := make(map[int32]engine.GotPath)
	for id, ancestorLookup := range all.Ids {
		intId, err := bullet_stl.AasciBulletIdToInt(id)
		if err != nil {
			return nil, err
		}
		intIds = append(intIds, int32(intId))

		fmt.Printf("VX: building ancestor path for %s\n", id)
		path, err := e.ancestorPathFrom(&ancestorLookup)
		if err != nil {
			return nil, err
		}
		if path != nil {
			ancestorPaths[int32(intId)] = *path
		}

	}
	titles, err := e.TitleStore.TitleForMany(intIds)
	if err != nil {
		return nil, err
	}

	//get string ids of all items to do the alias lookup
	stringIds := make([]string, len(all.Ids))

	i := 0
	for k := range all.Ids {
		stringIds[i] = k
		i++
	}

	aliases, err := e.LookupAliasForMany(stringIds)
	if err != nil {
		return nil, err
	}

	var summaryIds []SummaryId
	for _, v := range intIds {
		summaryIds = append(summaryIds, SummaryId(v))
	}
	summaries, err := e.SummaryStore.Fetch(summaryIds)
	if err != nil {
		return nil, err
	}
	//VX:TODO change summaryId to gotId and then fetch it here. Reusing the word summary is no good.
	//summaries, err := e.SummaryStore.Fetch()
	//VX:TODO lookup many here.
	var itemDisplays []engine.GotItemDisplay
	for k, v := range titles {

		stringId, err := bullet_stl.BulletIdIntToaasci(int64(k))
		if err != nil {
			return nil, err
		}

		var alias string = ""
		found, ok := aliases[stringId]
		if ok {
			alias = *found
		}
		var path *engine.GotPath = nil
		if foundPath, ok := ancestorPaths[k]; ok {
			path = &foundPath
		} else {
			fmt.Printf("VX: NO PATH FOR '%s'\n", v)
		}
		summaryText := "["
		gotId, err := engine.NewGotId(stringId)
		if err != nil {
			return nil, err
		}
		summaryId := NewSummaryId(*gotId)

		summary, ok := summaries[summaryId]
		if ok {

			if summary.State != nil {
				summaryText += "Leaf (" + summary.State.ToStr() + ")"
			}
			if summary.Counts != nil {
				summaryText += " {Total: "
				if summary.Counts.Active != 0 {
					summaryText += "active :" + strconv.Itoa(summary.Counts.Active)
				}
				if summary.Counts.Complete != 0 {
					summaryText += "complete :" + strconv.Itoa(summary.Counts.Complete)
				}
				if summary.Counts.Notes != 0 {
					summaryText += "notes :" + strconv.Itoa(summary.Counts.Notes)
				}
				summaryText += "}"
			}
		}
		summaryText += "]"
		itemDisplays = append(itemDisplays, engine.GotItemDisplay{
			Gid:     stringId,
			Title:   v,
			Path:    path,
			Alias:   alias,
			Summary: summaryText,
		})

	}

	sort.Slice(itemDisplays, func(i, j int) bool {
		a := itemDisplays[i]
		b := itemDisplays[j]
		if a.Path == nil && b.Path == nil {
			return a.Gid < b.Gid
		}
		if a.Path == nil {
			return true
		}
		if b.Path == nil {
			return false
		}
		lenA := len(a.Path.Ancestry)
		lenB := len(a.Path.Ancestry)
		if lenA == lenB {
			//wierd choice, but we just go chronoligal sorting for siblings.
			return a.Gid < b.Gid
		}
		return lenA < lenB

	})
	return e.renderSummaries(itemDisplays)

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
			Gid:      s.Gid,
			Alias:    s.Alias,
			NumberGo: num,
			Title:    s.Title,
			Path:     s.Path,
			Summary:  s.Summary,
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
	fmt.Printf("VX: updaitng state..\n")
	gid, err := e.GidLookup.InputToGid(&lookup)
	if err != nil {
		return err
	}
	if gid == nil {
		return nil
	}
	summaryId := SummaryId(gid.IntValue)
	ids := []SummaryId{summaryId}
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
	var summaryIds []SummaryId
	for _, id := range ancestorResult.Ids {
		summaryIds = append(summaryIds, SummaryId(id.IntValue))

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

func (e *EngineBullet) MarkResolved(lookup engine.GidLookup) (*engine.NodeId, error) {
	var newState engine.GotState = engine.Complete
	return nil, e.updateState(lookup, newState)
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
	fmt.Printf("VX: newId is %s\n", stringId)
	gotId := engine.GotId{
		AasciValue: stringId,
		IntValue:   newId,
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

	//VX:TODO DELETE
	//add item to ancestry
	ancestry, err := e.AncestorList.AddItem(gotId, parentGotId)
	if err != nil {
		return nil, err
	}

	//add item heading to depot
	err = e.TitleStore.AddItem(newId, heading)
	if err != nil {
		return nil, err
	}

	var summaryIds []SummaryId
	if ancestry != nil {
		fmt.Printf("VX: buck created. %+v\n", *ancestry)
		for _, a := range ancestry.Ids {
			summaryIds = append(summaryIds, SummaryId(a.IntValue))
		}
	}

	var newState engine.GotState = engine.Note
	if completable {
		newState = engine.Active
	}
	e.publishAddEvent(AddItemEvent{
		Id:       SummaryId(newId),
		State:    newState,
		Ancestry: summaryIds,
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
