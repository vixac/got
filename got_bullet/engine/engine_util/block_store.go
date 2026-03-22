package engine_util

//VX:TODO this is perhaps replaced by the collections implemtnation.
/*
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
func (b *BlockStore) ensureParenteExists(id int32) error {
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
func (b *BlockStore) AppendNote(id engine.GotId, block engine.LongFormBlock) error {

	err := b.ensureParenteExists(block.ParentID)
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
		//VX:TOO bollocks we need the track component.
	return nil, nil
}
func (b *BlockStore) LongFormForMany(ids []int32) (map[int32]engine.LongFormBlockResult, error) {
	return nil, nil
}
func (b *BlockStore) RemoveAllItemsFromLongStore(id int32) error {
	return nil

}
*/
