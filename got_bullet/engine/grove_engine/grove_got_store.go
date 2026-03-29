package grove_engine

import (
	"time"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	"vixac.com/got/engine"
)

const (
	groveStoreTreeId = "<g>"
	theRootNode      = bullet_interface.NodeID("0")
)

// this handles eveyrthing the got store needs to do.
type GotStoreInterface interface {
	CreateBuck(req GotStoreCreateRequest) error
	FetchBelow(id *engine.GotId) ([]GotIdWithDepth, error)
	FetchAncestorsForMany(gotIds []engine.GotId) ([]GotIdWithPath, error)
	AggregatesForMany(gotIds []engine.GotId) (map[engine.GotId]GotAggregate, error)
}

type GotIdWithDepth struct {
	Id                 engine.GotId
	DepthFromQueryNode int // Relative depth from query node (query node = 0, children = 1, etc.)
}
type GotIdWithPath struct {
	Id   engine.GotId
	Path []engine.GotId
}

type GotStoreCreateRequest struct {
	//Title            string
	Id    engine.GotId
	State engine.GotState
	//Deadline         *engine.DateTime
	//	OverrideSettings *engine.CreateOverrideSettings
	Parent *engine.GotId
}

type GotAggregate struct {
	Counts map[engine.GotState]int
}

type GroveGotStore struct {
	Grove bullet_interface.GroveClientInterface
}

// VX:TODO inject groveStoreTreeId
func NewGroveGotStore(grove bullet_interface.GroveClientInterface) (GotStoreInterface, error) {
	groveStore := GroveGotStore{
		Grove: grove,
	}

	return &groveStore, nil
}

func nodeFrom(id *engine.GotId) bullet_interface.NodeID {
	if id == nil {
		return theRootNode
	}
	return bullet_interface.NodeID(id.AasciValue)
}
func gotIdFrom(nodeId bullet_interface.NodeID) (*engine.GotId, error) {
	return engine.NewGotId(string(nodeId))
}

func (s *GroveGotStore) AggregatesForMany(gotIds []engine.GotId) (map[engine.GotId]GotAggregate, error) {
	var nodeIds []bullet_interface.NodeID
	for _, id := range gotIds {
		nodeIds = append(nodeIds, nodeFrom(&id))
	}
	req := bullet_interface.GroveGetNodeLocalAggregatesBulkRequest{
		TreeID:  groveStoreTreeId,
		NodeIDs: nodeIds,
	}
	aggs, err := s.Grove.GroveGetNodeLocalAggregatesBulk(req)
	if err != nil {
		return nil, err
	}
	result := make(map[engine.GotId]GotAggregate)
	for k, v := range aggs.Aggregates {
		gotId, err := gotIdFrom(k)
		if err != nil {
			return nil, err
		}
		groveAgg := NewAggregate(v)
		counts := make(map[engine.GotState]int)
		counts[engine.Active] = groveAgg.Active
		counts[engine.Complete] = groveAgg.Complete

		result[*gotId] = GotAggregate{
			Counts: counts,
		}
	}
	return result, nil
}

func (s *GroveGotStore) FetchAncestorsForMany(gotIds []engine.GotId) ([]GotIdWithPath, error) {
	var nodeIds []bullet_interface.NodeID
	for _, g := range gotIds {
		nodeIds = append(nodeIds, nodeFrom(&g))
	}
	req := bullet_interface.GroveGetAncestorsBulkRequest{
		TreeID:  groveStoreTreeId,
		NodeIDs: nodeIds,
	}
	res, err := s.Grove.GroveGetAncestorsBulk(req)
	if err != nil {
		return nil, err
	}

	var idsWithPaths []GotIdWithPath
	for k, v := range res.Ancestors {
		gotId, err := gotIdFrom(k)
		if err != nil {
			return nil, err
		}
		var pathIds []engine.GotId
		for _, ancestorNode := range v {
			ancestorId, err := gotIdFrom(ancestorNode)
			if err != nil {
				return nil, err
			}
			pathIds = append(pathIds, *ancestorId)
		}
		idsWithPaths = append(idsWithPaths, GotIdWithPath{
			Id:   *gotId,
			Path: pathIds,
		})

	}
	return idsWithPaths, nil
}

func (s *GroveGotStore) FetchBelow(id *engine.GotId) ([]GotIdWithDepth, error) {
	req := bullet_interface.GroveGetDescendantsRequest{
		NodeID:  nodeFrom(id),
		TreeID:  groveStoreTreeId,
		Options: nil,
	}
	res, err := s.Grove.GroveGetDescendants(req)
	if err != nil {
		return nil, err
	}
	var result []GotIdWithDepth

	for _, d := range res.Descendants {
		gotId, err := engine.NewGotId(string(d.NodeID))
		if err != nil || gotId == nil {
			return nil, err
		}
		result = append(result, GotIdWithDepth{
			Id:                 *gotId,
			DepthFromQueryNode: d.Depth,
		})
	}
	return result, nil
}

func (s *GroveGotStore) CreateBuck(createBuckRequest GotStoreCreateRequest) error {
	var parent *bullet_interface.NodeID = nil
	if createBuckRequest.Parent != nil {
		parentVal := bullet_interface.NodeID(createBuckRequest.Parent.AasciValue)
		parent = &parentVal
	}

	nodeId := nodeFrom(&createBuckRequest.Id)
	groveReq := bullet_interface.GroveCreateNodeRequest{
		NodeID:   nodeId,
		TreeID:   groveStoreTreeId,
		Parent:   parent,
		Position: nil,
		Metadata: nil,
	}
	err := s.Grove.GroveCreateNode(groveReq)
	if err != nil {
		return err
	}

	//now apply aggregates.
	mutationId := engine.TimeToMillisString(time.Now())

	deltas := NewMutationDelta(createBuckRequest.State)
	mutationReq := bullet_interface.GroveApplyAggregateMutationRequest{
		MutationID: bullet_interface.MutationID(mutationId),
		NodeID:     nodeId,
		TreeID:     groveStoreTreeId,
		Deltas:     deltas.ToGrove(),
	}
	err = s.Grove.GroveApplyAggregateMutation(mutationReq)

	return err
}
