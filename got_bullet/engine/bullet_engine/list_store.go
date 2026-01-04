package bullet_engine

import (
	"strconv"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/containers"
	"vixac.com/got/engine"
)

type ListLookupResult struct {
	Ids []engine.GotId
}
type ListsForGidResult struct {
	Lists []ListId
}

type ListId struct {
	Id int32
}
type ManyListLookupResult struct {
	Ids map[ListId]ListLookupResult
}

// The engine interface for whatever is going to store ancestor and descendant trees
type ListInterface interface {
	AddItem(id engine.GotId, list ListId) error
	RemoveItem(id engine.GotId, list ListId) error
	FetchListMembers(list ListId) (*ManyListLookupResult, error)
	FetchListsContaining(id engine.GotId) (*ListsForGidResult, error)
}

type BulletListStore struct {
	SubjectSeparator string
	Client           bullet_interface.TrackClientInterface
	Mesh             bullet_stl.Mesh
}

func NewListStore(client bullet_interface.TrackClientInterface, listStoreName string, bucketId int32, subjectSeparor string, forwardSeparator string, backwardSeparator string) (ListInterface, error) {
	mesh, err := bullet_stl.NewBulletMesh(client, bucketId, listStoreName, forwardSeparator, backwardSeparator)
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
	return bullet_stl.ListSubject{
		Value: strconv.Itoa(int(list.Id)),
	}
}

func toObject(id engine.GotId) bullet_stl.ListObject {
	return bullet_stl.ListObject{Value: id.AasciValue}
}

func (s *BulletListStore) AddItem(id engine.GotId, list ListId) error {

	//VX:TODO check existence of list?
	//we insert this item to the root node.

	parent := toSubject(list)

	object := toObject(id)
	pairs := []bullet_stl.ManyToManyPair{
		{Subject: parent, Object: object}, //list -> newItem
	}

	return s.Mesh.AppendPairs(pairs)
}

func (s *BulletListStore) RemoveItem(id engine.GotId, list ListId) error {
	manyPairs := []bullet_stl.ManyToManyPair{
		{
			Subject: toSubject(list),
			Object:  toObject(id),
		},
	}
	return s.Mesh.RemovePairs(manyPairs)
}

func (s *BulletListStore) FetchListMembers(list ListId) (*ManyListLookupResult, error) {
	//get the subject key for this id, and then use it as a prefix.
	subject := toSubject(list)
	allPairs, err := s.Mesh.AllPairsForSubject(subject)
	if err != nil || allPairs == nil {
		return nil, err
	}

	var ids []engine.GotId

	for _, pair := range allPairs.Pairs {
		id, err := engine.NewGotId(pair.Object.Value)
		if err != nil {
			return nil, err
		}
		ids = append(ids, *id)
	}

	lookupResult := ListLookupResult{Ids: ids}
	result := make(map[ListId]ListLookupResult)
	result[list] = lookupResult
	return &ManyListLookupResult{
		Ids: result,
	}, nil
}

func (s *BulletListStore) FetchListsContaining(id engine.GotId) (*ListsForGidResult, error) {
	object := toObject(id)
	allPairs, err := s.Mesh.AllPairsForObject(object)
	if err != nil || allPairs == nil {
		return nil, err
	}

	var lists []ListId
	for _, pair := range allPairs.Pairs {
		list, err := strconv.Atoi(pair.Subject.Value)
		if err != nil {
			return nil, err
		}
		lists = append(lists, ListId{Id: int32(list)}) //VX:TODO silent int32 cast
	}

	return &ListsForGidResult{
		Lists: lists,
	}, nil
}
