package engine_util

import (
	"strconv"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	"vixac.com/got/engine"
)

type BlockStore struct {
	Namespace int32
	Depot     bullet_interface.DepotClientInterface
	TreeId    string
	Grove     bullet_interface.GroveClientInterface
}

func NewBlockStore(Namespace int32, treeId string, grove bullet_interface.GroveClientInterface, depot bullet_interface.DepotClientInterface) (engine.LongFormStoreInterface, error) {

	return &BlockStore{
		Namespace: Namespace,
		TreeId:    treeId,
		Grove:     grove,
		Depot:     depot,
	}, nil

}

// checks a parent node exists and if it doesnt, create it at the top before returning thus ensuringParentExists.
func (b *BlockStore) ensureParentexists(id int32) error {
	//ask grove if the parent exists and create it if it doesnt.
	nodeId := bullet_interface.NodeID(id)
	req := bullet_interface.GroveExistsRequest{
		NodeID: nodeId,
		TreeID: bullet_interface.TreeID(b.TreeId),
	}
	res, err := b.Grove.GroveExists(req)
	if err != nil {
		return err
	}

	if res.Exists {
		return nil
	}
	//create the node at top level

	createReq := bullet_interface.GroveCreateNodeRequest{
		NodeID:   nodeId,
		TreeID:   bullet_interface.TreeID(b.TreeId),
		Position: nil,
		Metadata: nil,
		Parent:   nil,
	}
	return b.Grove.GroveCreateNode(createReq)
}

// VX:TODO rm the id from this.
// wtf is child position. the epoch i guess.
func (b *BlockStore) UpsertItem(id int32, block engine.LongFormBlock) error {

	err := b.ensureParentexists(block.ParentID)
	if err != nil {
		return err
	}

	var childPosition bullet_interface.ChildPosition = bullet_interface.ChildPosition(block.Edited.Millis)
	parentIdString := strconv.Itoa(int(block.ParentID))
	parentNodeId := bullet_interface.NodeID(parentIdString)
	req := bullet_interface.GroveCreateNodeRequest{
		NodeID:   bullet_interface.NodeID(block.Id.String),
		TreeID:   bullet_interface.TreeID(b.TreeId),
		Position: &childPosition,
		Metadata: nil,
		Parent:   &parentNodeId,
	}
	err = b.Grove.GroveCreateNode(req) //VX:TODO what if the parent doesn't exist.
	if err != nil {
		return err
	}

	//VX:TODO FUCKKKKK. depot is annoying as hell. I should add a layer to depot that allows keys
	//at this point we have the entry in grove but not in depot. So.. shit.

	//depotReq := bullet_interface.DepotRequest{}
	//b.Depot.DepotInsertOne()
	return nil
}

// VX:TODO bollocks.
func (b *BlockStore) LongFormFor(id int32) (*engine.LongFormBlockResult, error) {
	parentNodeId := bullet_interface.NodeID(id)
	req := bullet_interface.GroveGetChildrenRequest{
		NodeID:     parentNodeId,
		TreeID:     bullet_interface.TreeID(b.TreeId),
		Pagination: nil,
	}
	_, err := b.Grove.GroveGetChildren(req)
	if err != nil {
		return nil, err
	}
	/*
		var blocks []engine.LongFormBlock
		for _, c := range res.Children {

		}
	*/
	return nil, nil
}
func (b *BlockStore) LongFormForMany(ids []int32) (map[int32]engine.LongFormBlockResult, error) {
	return nil, nil
}
func (b *BlockStore) RemoveAllItemsFromLongStore(id int32) error {
	return nil

}
