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
	Ids []engine.GotId
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
	AddItem(id engine.GotId, under *engine.GotId) error
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

//aah crap it needs to be a two way mesh.

func (a *BulletAncestorList) AddItem(id engine.GotId, under *engine.GotId) error {
	ancestors, err := a.Mesh.AllPairsForObject(bullet_stl.ListObject{Value: id.AasciValue})
	if err != nil {
		return err
	}

	//can't insert an item that exists.
	if ancestors != nil && len(ancestors.Pairs) != 0 {
		return errors.New("attempted to insert an existing id")
	}
	//we insert this item to the root node.
	var parent bullet_stl.ListSubject
	if under == nil {
		parent = TheRootNode
	} else {
		parent = bullet_stl.ListSubject{Value: under.AasciValue}

		//we attempt to delete theleafNode from this parent, which will succeed if this item is the first child.
		deletePairs := []bullet_stl.ManyToManyPair{{Subject: parent, Object: TheLeafChild}}
		err := a.Mesh.RemovePairs(deletePairs) //VX:TODO make sure deleting doesnt fail if theres nothing to delete.
		if err != nil {
			return err
		}
	}
	object := bullet_stl.ListObject{Value: id.AasciValue}

	pairs := []bullet_stl.ManyToManyPair{
		{Subject: parent, Object: object},                //parent -> newItem
		{Subject: object.Invert(), Object: TheLeafChild}, //newItem -> theLeafNode
	}
	return a.Mesh.AppendPairs(pairs)
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
		fmt.Printf("VX: pair is %s -> %s", pair.Subject.Value, pair.Object.Value)
	}

	return nil, errors.New("not impl")

}
func (a *BulletAncestorList) FetchImmediatelyUnder(id engine.GotId) (*DescendantLookupResult, error) {
	descendants, err := a.Mesh.AllPairsForSubject(bullet_stl.ListSubject{Value: id.AasciValue})
	if err != nil {
		return nil, err
	}
	//VX:TODO
	for _, pair := range descendants.Pairs {
		fmt.Printf("VX: pair is %s -> %s", pair.Subject.Value, pair.Object.Value)
	}
	return nil, errors.New("not impl")

}

func (a *BulletAncestorList) FetchAncestorsOf(id engine.GotId) (*AncestorLookupResult, error) {
	ancestors, err := a.Mesh.AllPairsForObject(bullet_stl.ListObject{Value: id.AasciValue})
	if err != nil {
		return nil, err
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
