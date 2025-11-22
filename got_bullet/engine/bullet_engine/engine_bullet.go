package bullet_engine

import (
	"errors"
	"fmt"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/ids"

	"vixac.com/got/engine"
)

const (
	aliasBucket    int32 = 1001
	nodeBucket     int32 = 1002
	ancestorBucket int32 = 1003
)

type EngineBullet struct {
	Client       bullet_interface.BulletClientInterface
	AncestorList AncestorListInterface
	TitleStore   TitleStoreInterface
	GidLookup    GidLookupInterface
}

func NewEngineBullet(client bullet_interface.BulletClientInterface) (*EngineBullet, error) {
	ancestorList, err := NewAncestorList(client, "<anc>", ancestorBucket, ":", ">", "<")

	if err != nil {
		return nil, err
	}

	titleStore, err := NewBulletTitleStore(client)
	if err != nil {
		return nil, err
	}

	gidLookup, err := NewBulletGidLookup()
	if err != nil {
		return nil, err
	}
	return &EngineBullet{
		Client:       client,
		AncestorList: ancestorList,
		TitleStore:   titleStore,
		GidLookup:    gidLookup,
	}, nil
}

func (e *EngineBullet) Summary(lookup *engine.GidLookup) (*engine.GotSummary, error) {

	gid, err := e.GidLookup.InputToGid(lookup)
	if err != nil {
		return nil, err
	}
	if gid == nil {
		return nil, errors.New("no gid")
	}

	//descendants, err := e.AncestorList.FetchAncestorsOf(*gid)
	title, err := e.TitleStore.TitleFor(gid.IntValue)
	if err != nil {
		return nil, errors.New("no title")
	}

	//this is building the path.
	/*

		ancestorResult, err := e.AncestorList.FetchAncestorsOf(*gid)
		if err != nil {
			return nil, errors.New("error fetching ancestors")
		}

		var path string = ""
		if ancestorResult != nil {
			var stringIds []string
			for _, id := range ancestorResult.Ids {
				stringIds = append(stringIds, id.AasciValue)

			}
			path = strings.Join(stringIds, "->")

		}*/
	var theTitle = ""
	if title != nil {
		theTitle = *title
	}
	return &engine.GotSummary{
		Gid:   gid.AasciValue,
		Title: theTitle,
	}, nil

}

func (e *EngineBullet) Alias(gid string, alias string) (bool, error) {
	return false, errors.New(("not impl"))
}

func (e *EngineBullet) Resolve(lookup engine.GidLookup) (*engine.NodeId, error) {
	//check if the gid is an exact match for an item id
	//check int32 parse, check its length is the right length

	//aliases can't start with a number.
	return nil, errors.New("not impl")
}

func (e *EngineBullet) Delete(lookup engine.GidLookup) (*engine.NodeId, error) {
	//check if the gid is an exact match for an item id
	//check int32 parse, check its length is the right length

	//aliases can't start with a number.
	return nil, errors.New("not impl")
}

func (e *EngineBullet) Unalias(alias string) (*engine.NodeId, error) {
	//check if the gid is an exact match for an item id
	//check int32 parse, check its length is the right length

	//aliases can't start with a number.
	return nil, errors.New("not impl")
}

func (e *EngineBullet) Move(lookup engine.GidLookup, newParent engine.GidLookup) (*engine.NodeId, error) {
	return nil, errors.New("not impl")
}

func (e *EngineBullet) CreateBuck(parent *engine.GidLookup, date *engine.DateLookup, completable bool, heading string) (*engine.NodeId, error) {
	//VX:TODO this should hit both the keys and also hit depot too for the heading.

	newId, err := e.NextId()

	if err != nil {
		return nil, err
	}
	stringId, err := bullet_stl.BulletIdIntToaasci(newId)
	if err != nil {
		return nil, err
	}
	fmt.Printf("VX: newId is %s", stringId)
	gotId := engine.GotId{
		AasciValue: stringId,
		IntValue:   newId,
	}

	//VX:TODO lookup parent
	if parent != nil {
		return nil, errors.New("adding under parent not supported yet")
	}

	//add item to ancestry
	err = e.AncestorList.AddItem(gotId, nil)
	if err != nil {
		return nil, err
	}

	//add item heading to depot
	err = e.TitleStore.AddItem(newId, heading)
	if err != nil {
		return nil, err
	}

	fmt.Printf("VX: buck created.\n")
	return &engine.NodeId{
		Gid: engine.Gid{
			Id: stringId,
		},
		Title: heading,
		Alias: "",
	}, nil
}

func (e *EngineBullet) Lookup(alias string) (*engine.NodeId, error) {
	return nil, errors.New("not impl")

}

/**
The cache is going to work like this

asc(n)
and its only ascendants and descen

*/
/**
ok im goina do this the buck way..
it takes a wayfinder and uses it
and the business logic for using the wayfinder is reusable.
its a cool design.
*/
