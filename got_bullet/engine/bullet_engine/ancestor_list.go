package bullet_engine

import (
	"errors"
	"fmt"
	"strings"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/containers"
	"vixac.com/got/engine"
)

type AncestorLookupResult struct {
	Ids []engine.GotId
}

type DescendantLookupResult struct {
	//each decendant gid mapped to their AncestorLookupResult
	Ids map[string]AncestorLookupResult
}

// 0 is a bit like a null terminator character. Beacuse ancestor list is forward only
// duplicate objects aren't a problem. We'll give this value to all nodes that are infact a leaf node.
var (
	TheLeafChild = bullet_stl.ListObject{
		Value: "00",
	}
	//if no parent is provided, then the root ancestor is provided.
	TheRootNode = bullet_stl.ListSubject{
		Value: "0",
	}
)

// The engine interface for whatever is going to store ancestor and descendant trees
type AncestorListInterface interface {
	AddItem(id engine.GotId, under *engine.GotId) (*Ancestry, error)
	RemoveItem(id engine.GotId) error
	FetchAllItems(under engine.GotId) (*DescendantLookupResult, error)
	FetchImmediatelyUnder(id engine.GotId) (*DescendantLookupResult, error)
	FetchAncestorsOf(id engine.GotId) (*AncestorLookupResult, error)

	//VX:TODO move semantics oh dear. MoveItem(id engine.GotId, under *engine.GotId) error
}

// subject separator is used to divide the subject into a list of ancestors, eg a:b:c, and then the objcet separator
// is used to separate the subject and object inside the mesh, leading to myList>a:b:c>object
type BulletAncestorList struct {
	SubjectSeparator string
	//	ObjectSeparator  string
	Client bullet_interface.TrackClientInterface
	//	ListName         string
	//	BucketId         int
	Mesh bullet_stl.Mesh
}

func NewAncestorList(client bullet_interface.TrackClientInterface, listName string, bucketId int32, subjectSeparor string, forwardSeparator string, backwardSeparator string) (*BulletAncestorList, error) {
	mesh, err := bullet_stl.NewBulletMesh(client, bucketId, listName, forwardSeparator, backwardSeparator)
	if err != nil {
		return nil, err
	}

	return &BulletAncestorList{
		SubjectSeparator: subjectSeparor,
		Client:           client,
		Mesh:             mesh,
	}, nil
}

type Ancestry struct {
	Ids []engine.GotId
}

//aah crap it needs to be a two way mesh.

func (a *BulletAncestorList) AddItem(id engine.GotId, under *engine.GotId) (*Ancestry, error) {
	fmt.Printf("attempting to insert id %s into ancestry\n", id.AasciValue)
	if id.AasciValue == TheRootNode.Value {
		return nil, errors.New("inserting the root node is not permitted")
	}

	ancestorsOfNewItem, err := a.Mesh.AllPairsForObject(bullet_stl.ListObject{Value: id.AasciValue})
	if err != nil {
		return nil, err
	}

	//can't insert an item that exists, and all items have an acnestor besides the root node.
	if ancestorsOfNewItem != nil && len(ancestorsOfNewItem.Pairs) != 0 {
		return nil, errors.New("attempted to insert an existing id")
	}
	//we insert this item to the root node.
	var parent bullet_stl.ListSubject

	var ancestry *Ancestry = nil
	if under == nil {
		parent = TheRootNode
	} else {
		//now we construct the ancestor prefix
		ancestors, err := a.Mesh.AllPairsForObject(bullet_stl.ListObject{Value: under.AasciValue})
		if err != nil {
			return nil, err
		}
		if ancestors == nil {
			return nil, errors.New("every node bedies the root node should have ancestors")
		}
		fmt.Printf("VX: WWTFF+==================\n")

		//VX:TODO finish
		var ancestorGotIdList []engine.GotId
		var ancestorList []string
		if len(ancestors.Pairs) > 1 {
			fmt.Printf("VX: IM PRETTY SURE THIS NEVER HAPPENS\n")

		}
		//		ancestorKey := ancestors.Pairs[0]
		//VX:TODO just pick the 1st item.
		for _, a := range ancestors.Pairs { //there's always 1 ancestor here no? Wierd. I think i
			fmt.Printf("VX: This is an ancestor '%s'\n", a.Subject.Value)
			ancestorList = append(ancestorList, a.Subject.Value)
		}

		ancestorList = append(ancestorList, under.AasciValue)
		parentKey := strings.Join(ancestorList, a.SubjectSeparator)

		eachAncestor := strings.Split(parentKey, a.SubjectSeparator)

		//splitting the parent key at the end is a bit cheap but it works.
		for _, a := range eachAncestor {
			gotId, err := engine.NewGotId(a)
			if err != nil {
				return nil, err
			}
			ancestorGotIdList = append(ancestorGotIdList, *gotId)

		}
		ancestry = &Ancestry{
			Ids: ancestorGotIdList,
		}

		fmt.Printf("VX: complete parent key is %s\n", parentKey)
		parent = bullet_stl.ListSubject{Value: parentKey}

		//we attempt to delete theleafNode from this parent, which will succeed if this item is the first child.
		deletePairs := []bullet_stl.ManyToManyPair{{Subject: parent, Object: TheLeafChild}}
		err = a.Mesh.RemovePairs(deletePairs) //VX:TODO make sure deleting doesnt fail if theres nothing to delete.
		if err != nil {
			return nil, err
		}
	}
	object := bullet_stl.ListObject{Value: id.AasciValue}

	pairs := []bullet_stl.ManyToManyPair{
		{Subject: parent, Object: object},                //parent -> newItem
		{Subject: object.Invert(), Object: TheLeafChild}, //newItem -> theLeafNode
	}

	return ancestry, a.Mesh.AppendPairs(pairs)
}

func (a *BulletAncestorList) RemoveItem(id engine.GotId) error {
	return errors.New("not impl")
}
func (a *BulletAncestorList) FetchAllItems(under engine.GotId) (*DescendantLookupResult, error) {
	descendants, err := a.Mesh.AllPairsForPrefixSubject(bullet_stl.ListSubject{Value: under.AasciValue})
	if err != nil {
		return nil, err
	}
	for _, pair := range descendants.Pairs {
		fmt.Printf("VX: all pair is %s -> %s\n", pair.Subject.Value, pair.Object.Value)
	}
	return nil, errors.New("not impl")
}

func (a *BulletAncestorList) FetchImmediatelyUnder(id engine.GotId) (*DescendantLookupResult, error) {
	//get the subject key for this id, and then use it as a prefix.

	var ancestorKey = "" //this can be left blank for TheRootNote.
	if id.AasciValue != TheRootNode.Value {
		ancestorPairs, err := a.Mesh.AllPairsForObject(bullet_stl.ListObject{Value: id.AasciValue})
		if err != nil {
			return nil, err
		}
		if ancestorPairs == nil || len(ancestorPairs.Pairs) != 1 {
			return nil, errors.New("zero ancestors. Itemas are su")
		}
		//append the query to the ancestor Key, so id = c, fetches a:b, and we want to lookup everything prefixed with a:b:c
		ancestorKey = ancestorPairs.Pairs[0].Subject.Value + a.SubjectSeparator + id.AasciValue

	}

	everythingBelowAncestor, err := a.Mesh.AllPairsForPrefixSubject(bullet_stl.ListSubject{Value: ancestorKey})
	if err != nil {
		return nil, err
	}
	if everythingBelowAncestor == nil {
		return nil, nil
	}
	//VX:TODO FINISH
	ids := make(map[string]AncestorLookupResult)
	for _, pair := range everythingBelowAncestor.Pairs {
		ancestorsIndividualIds := strings.Split(pair.Subject.Value, a.SubjectSeparator)
		var gids []engine.GotId
		for _, ancestorId := range ancestorsIndividualIds {
			gid, err := engine.NewGotId(ancestorId)
			if err != nil {
				return nil, err
			}
			if gid == nil {
				return nil, errors.New("empty gid")
			}
			gids = append(gids, *gid)
		}
		fmt.Printf("VX: Found item under %s,  subject  %s -> object %s\n", id.AasciValue, pair.Subject.Value, pair.Object.Value)
		ids[pair.Object.Value] = AncestorLookupResult{
			Ids: gids,
		}

	}

	return &DescendantLookupResult{
		Ids: ids,
	}, nil
}

func (a *BulletAncestorList) FetchAncestorsOf(id engine.GotId) (*AncestorLookupResult, error) {

	ancestors, err := a.Mesh.AllPairsForObject(bullet_stl.ListObject{Value: id.AasciValue})
	if err != nil {
		return nil, err
	}
	for i, v := range ancestors.Pairs {
		fmt.Printf("VX: %d ancestor of %s is pair ancestor %s -> %s\n", i, id.AasciValue, v.Subject.Value, v.Object.Value)
	}
	if len(ancestors.Pairs) != 1 { //the ancestor list is in the form of the subject of the object, so its 1 key, containing all ancestors.
		return nil, errors.New("this id seems to appear more than once")
	}
	keyString := ancestors.Pairs[0].Subject.Value
	ancestorSplit := strings.Split(keyString, a.SubjectSeparator)

	var ids []engine.GotId
	for _, ancestor := range ancestorSplit {
		id, err := engine.NewGotId(ancestor)
		if err != nil {
			return nil, err
		}
		ids = append(ids, *id)
	}

	return &AncestorLookupResult{
		Ids: ids,
	}, nil
}
