package grove_engine

import (
	"fmt"
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

	// this is the individual contribution to aggregates for each node. There is no addition being done here.
	IndividualStateForMany(gotIds []engine.GotId) (map[engine.GotId]engine.GotState, error)

	//returns the aggregates of the descendants
	AggregatesOfDescendantsForMany(gotIds []engine.GotId) (map[engine.GotId]GotAggregate, error)
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

// VX:Note duplication here.
func (s *GroveGotStore) AggregatesOfDescendantsForMany(gotIds []engine.GotId) (map[engine.GotId]GotAggregate, error) {
	var nodeIds []bullet_interface.NodeID
	for _, id := range gotIds {
		nodeIds = append(nodeIds, nodeFrom(&id))
	}
	req := bullet_interface.GroveGetNodeWithDescendantsAggregatesBulkRequest{
		TreeID:  groveStoreTreeId,
		NodeIDs: nodeIds,
	}
	aggs, err := s.Grove.GroveGetNodeWithDescendantsAggregatesBulk(req)
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

func (s *GroveGotStore) IndividualStateForMany(gotIds []engine.GotId) (map[engine.GotId]engine.GotState, error) {
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
	result := make(map[engine.GotId]engine.GotState)
	for k, v := range aggs.Aggregates {
		gotId, err := gotIdFrom(k)
		if err != nil {
			return nil, err
		}
		groveAgg := NewAggregate(v)
		//VX:TODO for now, if its not active, its complete.
		if groveAgg.Active == 1 {
			result[*gotId] = engine.Active
		} else {
			result[*gotId] = engine.Complete
		}
	}
	//now we expect each aggregate to be 1 for one state and 0 for the rest, or no state at all.

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
func (s *GroveGotStore) createBuckAttempt(createBuckRequest GotStoreCreateRequest, numTries int) error {
	var parent *bullet_interface.NodeID = nil
	var parentIsRoot = false
	if createBuckRequest.Parent != nil {
		fmt.Printf("VX: poarent is %s\n", *createBuckRequest.Parent)
		parentVal := bullet_interface.NodeID(createBuckRequest.Parent.AasciValue)
		parent = &parentVal
	} else {
		parentIsRoot = true
		rootId, _ := engine.NewGotIdFromInt(0)
		rootParent := nodeFrom(rootId)
		parent = &rootParent
		fmt.Printf("VX: user entered no parent. Should we check the parent exists?")
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

		if parentIsRoot && numTries == 0 {
			fmt.Printf("VX: ok This must be the first node ver. Lets try to create the root node and then try again.")

			groveRootNodeReq := bullet_interface.GroveCreateNodeRequest{
				NodeID:   *parent,
				TreeID:   groveStoreTreeId,
				Parent:   nil,
				Position: nil,
				Metadata: nil,
			}
			err := s.Grove.GroveCreateNode(groveRootNodeReq)
			if err != nil {
				fmt.Printf("VX: Failed to create the root node.\n")
				return err
			}
			return s.createBuckAttempt(createBuckRequest, 1)
			//This might be the first buck entered, so we need to create the root node.

		}
		fmt.Printf("VX: grove create node failed.. should we check her?")
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
func (s *GroveGotStore) CreateBuck(createBuckRequest GotStoreCreateRequest) error {
	return s.createBuckAttempt(createBuckRequest, 0)

}
