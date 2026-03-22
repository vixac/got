package engine_util

import (
	"errors"
	"fmt"
	"time"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/containers"
	bullet_id "github.com/vixac/firbolg_clients/bullet/bullet_stl/ids"
	"vixac.com/got/engine"
)

type BulletLongFormStore struct {
	Collection bullet_stl.Collection
}

func NewBulletLongFormStore(bucketId int32, track bullet_interface.TrackClientInterface, depot bullet_interface.DepotClientInterface) (engine.LongFormStoreInterface, error) {
	coll := bullet_stl.NewBulletCollection(bucketId, track, depot)
	return &BulletLongFormStore{Collection: coll}, nil
}

func highestIdInside(collection map[bullet_stl.CollectionId]bullet_stl.CollectionItem) (*bullet_id.BulletId, error) {
	if len(collection) == 0 {
		return nil, nil
	}
	var highestIntValue int64 = 0
	for k, _ := range collection {

		longFormKey, err := engine.NewLongFormKeyFromString(k.Key)
		if err != nil {
			return nil, err
		}
		if longFormKey.NoteId.IntValue > highestIntValue {
			highestIntValue = longFormKey.NoteId.IntValue
		}
	}
	return bullet_id.NewBulletIdFromInt(highestIntValue)
}

func (s *BulletLongFormStore) AppendNote(id engine.GotId, block engine.LongFormBlock) error {
	idStr := idToStr(id)
	existing, err := s.Collection.AllItemsUnderPrefix(idStr)
	if err != nil {
		return err
	}

	highestExistingId, err := highestIdInside(existing)
	if err != nil {
		return err
	}
	if highestExistingId == nil { //this is the first note for this gotid
		first := engine.FirstNoteId()
		highestExistingId = &first
	}

	now := time.Now()

	newLongFormId := engine.LongFormKey{
		NoteId:      highestExistingId.Next(),
		GotId:       id,
		CreatedTime: now,
	}

	newLongFormNoteStringId := newLongFormId.ToString()
	collId, err := s.Collection.CreateItemUnder(newLongFormNoteStringId, block.Content, &now)
	if err != nil {
		return err
	}

	fmt.Printf("VX: Note created under colelction Id %s \n", collId.Key)
	return nil
}

/*
// VX:TODO RM?

	func getTheOnlyKeyOrError(result map[bullet_stl.CollectionId]bullet_stl.CollectionItem) (*bullet_stl.CollectionId, error) {
		if len(result) != 1 {
			return nil, errors.New("upserting to a key that is not unique")
		}
		var collId *bullet_stl.CollectionId = nil
		for k := range result {
			collId = &k
		}
		return collId, nil
	}
*/
func collectionToLongFormMap(collection map[bullet_stl.CollectionId]bullet_stl.CollectionItem) (map[engine.GotId]engine.LongFormBlockResult, error) {
	idsToBlocks := make(map[engine.GotId][]engine.LongFormBlock)
	for k, v := range collection {
		longformId, err := engine.NewLongFormKeyFromString(k.Key)
		if err != nil || longformId == nil {
			return nil, err
		}
		edited, err := engine.NewDateTime(v.Updated)
		if err != nil {
			return nil, err
		}
		newBlock := engine.LongFormBlock{
			Id:      *longformId,
			Content: v.Payload,
			Edited:  edited,
		}
		gotId := longformId.GotId
		existing, ok := idsToBlocks[gotId]
		if !ok {
			idsToBlocks[gotId] = []engine.LongFormBlock{newBlock}
		} else {
			idsToBlocks[gotId] = append(existing, newBlock)
		}
	}

	result := make(map[engine.GotId]engine.LongFormBlockResult)
	for k, v := range idsToBlocks {
		result[k] = engine.LongFormBlockResult{
			Blocks: v,
		}
	}
	fmt.Printf("VX:TODO sort the arrays by edited time or created time")
	return result, nil
}

func (s *BulletLongFormStore) LongFormForMany(ids []engine.GotId) (map[engine.GotId]engine.LongFormBlockResult, error) {
	var idStrings []string
	for _, id := range ids {
		idStrings = append(idStrings, idToStr(id))
	}
	items, err := s.Collection.AllItemsUnderPrefixes(idStrings)
	if err != nil {
		return nil, err
	}
	return collectionToLongFormMap(items)
}

func (s *BulletLongFormStore) LongFormNotesFor(id engine.GotId) (*engine.LongFormBlockResult, error) {
	idStr := idToStr(id)
	res, err := s.Collection.AllItemsUnderPrefix(idStr)
	if err != nil || len(res) == 0 {
		return nil, err
	}

	idMap, err := collectionToLongFormMap(res)
	if err != nil {
		return nil, err
	}
	//no notes for this id
	if len(idMap) == 0 {
		return nil, nil
	}
	if len(idMap) != 1 {
		return nil, errors.New("Too many gotIds in response for notes for id")
	}

	blockResult, ok := idMap[id]
	if !ok {
		return nil, errors.New("wrong id returned")
	}
	return &blockResult, nil
	/*
		var blocks []engine.LongFormBlock
		fmt.Printf("VX:WARN: collections dont return created and updated times so longform doesnt have it.")
		for collId, value := range res {
			edited, err := engine.NewDateTime(value.Updated)
			if err != nil {
				return nil, err
			}
			block := engine.LongFormBlock{
				Id: engine.LongFormId{
					String: collId.Key,
				},
				ParentID: id,
				Content:  value.Payload,
				Created:  engine.DateTime{}, //VX:TODO
				Edited:   edited,
			}
			blocks = append(blocks, block)
		}
		blockResult := engine.LongFormBlockResult{
			Blocks: blocks,
		}
	*/
}

func (s *BulletLongFormStore) RemoveAllItemsFromLongStoreUnder(id engine.GotId) error {

	items, err := s.Collection.AllItemsUnderPrefix(id.AasciValue)
	if err != nil {
		return err
	}

	var collIds []bullet_stl.CollectionId
	for k := range items {
		collIds = append(collIds, k)
	}
	return s.Collection.DeleteItems(collIds)
}
