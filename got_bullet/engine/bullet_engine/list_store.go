package bullet_engine

import (
	"errors"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/containers"
	"vixac.com/got/engine"
)

type ListLookupResult struct {
	Ids []engine.GotId
}

type ListId struct {
	Id int32
}
type ManyListLookupResult struct {
	//each decendant gid mapped to their ListLookupResult
	Ids map[ListId]ListLookupResult
}

// The engine interface for whatever is going to store ancestor and descendant trees
type ListInterface interface {
	AddItem(id engine.GotId, list ListId) error
	RemoveItem(id engine.GotId, list ListId) error
	FetchListMembers(id ListId) (*ManyListLookupResult, error)
	FetchListsContaining(id engine.GotId) (*ListLookupResult, error)

	//VX:TODO move semantics oh dear. MoveItem(id engine.GotId, under *engine.GotId) error
}

// subject separator is used to divide the subject into a list of ancestors, eg a:b:c, and then the objcet separator
// is used to separate the subject and object inside the mesh, leading to myList>a:b:c>object
type BulletListStore struct {
	SubjectSeparator string
	//	ObjectSeparator  string
	Client bullet_interface.TrackClientInterface
	//	ListName         string
	//	BucketId         int
	Mesh bullet_stl.Mesh
}

func NewListStore(client bullet_interface.TrackClientInterface, listName string, bucketId int32, subjectSeparor string, forwardSeparator string, backwardSeparator string) (ListInterface, error) {
	mesh, err := bullet_stl.NewBulletMesh(client, bucketId, listName, forwardSeparator, backwardSeparator)
	if err != nil {
		return nil, err
	}

	return &BulletListStore{
		SubjectSeparator: subjectSeparor,
		Client:           client,
		Mesh:             mesh,
	}, nil
}

func toSubject(list ListId) bullet_stl.ListSubject {
	aasci := string(list.Id)
	return bullet_stl.ListSubject{Value: aasci}
}

//aah crap it needs to be a two way mesh.

func (s *BulletListStore) AddItem(id engine.GotId, list ListId) error {

	//VX:TODO check existence of list?
	//we insert this item to the root node.

	parent := toSubject(list)

	object := bullet_stl.ListObject{Value: id.AasciValue}

	pairs := []bullet_stl.ManyToManyPair{
		{Subject: parent, Object: object}, //parent -> newItem
	}

	return s.Mesh.AppendPairs(pairs)
}

func (a *BulletListStore) RemoveItem(id engine.GotId, list ListId) error {
	return errors.New("not impl")
}

func (s *BulletListStore) FetchListMembers(id ListId) (*ManyListLookupResult, error) {
	//get the subject key for this id, and then use it as a prefix.
	return nil, errors.New("not impl")
}

func (s *BulletListStore) FetchListsContaining(id engine.GotId) (*ListLookupResult, error) {
	return nil, errors.New("not impl")
}
