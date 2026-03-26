package grove_engine

import (
	"fmt"
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
	Id               engine.GotId
	State            engine.GotState
	Deadline         *engine.DateTime
	OverrideSettings *engine.CreateOverrideSettings
	Parent           *engine.GotId
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

func stringToTime(dateString string) (engine.DateTime, error) {
	createdDate, err := engine.NewTimeFromString(dateString)
	if err != nil {
		return engine.DateTime{}, err
	}
	return engine.NewDateTime(time.Time(*createdDate))

}
func createdTimeOrNil(override *engine.CreateOverrideSettings) (*engine.DateTime, error) {

	if override != nil {
		date, err := stringToTime(override.CreatedDate)
		return &date, err
	}
	return nil, nil
}
func updatedTimeOrNil(override *engine.CreateOverrideSettings) (*engine.DateTime, error) {

	if override != nil {
		date, err := stringToTime(override.UpdatedDate)
		return &date, err
	}
	return nil, nil
}

func (s *GroveGotStore) CreateBuck(createBuckRequest GotStoreCreateRequest) error {
	var parent *bullet_interface.NodeID = nil
	if createBuckRequest.Parent != nil {
		parentVal := bullet_interface.NodeID(createBuckRequest.Parent.AasciValue)
		parent = &parentVal
	}

	now, _ := engine.NewDateTime(time.Now())
	created, err := createdTimeOrNil(createBuckRequest.OverrideSettings)
	if err != nil {
		return err
	}
	updated, err := updatedTimeOrNil(createBuckRequest.OverrideSettings)
	if err != nil {
		return err
	}
	if created == nil {
		created = &now
	}
	if updated == nil {
		updated = &now
	}

	var childPosition float64 = float64(updated.EpochMillis())
	position := bullet_interface.ChildPosition(childPosition)

	var customFlags []string
	if createBuckRequest.OverrideSettings != nil && createBuckRequest.OverrideSettings.Flags != nil {
		customFlags = createBuckRequest.OverrideSettings.Flags
	}
	fmt.Printf("VX:TODO Scheduled parameter is not currently being stored in grove store")
	meta := NewGroveMetaData(customFlags, *created, nil)

	nodeId := bullet_interface.NodeID(createBuckRequest.Id.AasciValue)
	groveReq := bullet_interface.GroveCreateNodeRequest{
		NodeID:   nodeId,
		TreeID:   groveStoreTreeId,
		Parent:   parent,
		Position: &position,
		Metadata: meta.ToGrove(),
	}
	err = s.Grove.GroveCreateNode(groveReq)
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
