package grove_engine

import (
	"time"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	"vixac.com/got/engine"
)

const (
	groveStoreTreeId = "<g>"
)

// this handles eveyrthing the got store needs to do.
type GotStoreInterface interface {
	CreateBuck(req GotStoreCreateRequest) error
}

type GotStoreCreateRequest struct {
	//Title            string
	Id    engine.GotId
	State engine.GotState
	//Deadline         *engine.DateTime
	//	OverrideSettings *engine.CreateOverrideSettings
	Parent *engine.GotId
}

type GroveGotStore struct {
	Grove bullet_interface.GroveClientInterface
}

func NewGroveGotStore(grove bullet_interface.GroveClientInterface) (GotStoreInterface, error) {
	groveStore := GroveGotStore{
		Grove: grove,
	}

	return &groveStore, nil
}

func (s *GroveGotStore) CreateBuck(createBuckRequest GotStoreCreateRequest) error {
	var parent *bullet_interface.NodeID = nil
	if createBuckRequest.Parent != nil {
		parentVal := bullet_interface.NodeID(createBuckRequest.Parent.AasciValue)
		parent = &parentVal
	}

	nodeId := bullet_interface.NodeID(createBuckRequest.Id.AasciValue)
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
