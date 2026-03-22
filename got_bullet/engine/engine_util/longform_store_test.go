package engine_util

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/containers"
	"vixac.com/got/engine"
)

func buildTestLongFormStore(t *testing.T) *BulletLongFormStore {
	client := engine.BuildTestClient()
	coll := bullet_stl.NewBulletCollection(0, client, client)
	return &BulletLongFormStore{Collection: coll}
}

func makeBlock(content string) engine.LongFormBlock {
	return engine.LongFormBlock{Content: content}
}

func sortedBlocks(blocks []engine.LongFormBlock) []engine.LongFormBlock {
	sorted := make([]engine.LongFormBlock, len(blocks))
	copy(sorted, blocks)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Id.NoteId.IntValue < sorted[j].Id.NoteId.IntValue
	})
	return sorted
}

func TestLongFormStore_NoNotes(t *testing.T) {
	store := buildTestLongFormStore(t)

	gotId, err := engine.NewGotIdFromInt(10)
	assert.NoError(t, err)

	result, err := store.LongFormNotesFor(*gotId)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestLongFormStore_AppendFirstNote(t *testing.T) {
	store := buildTestLongFormStore(t)

	gotId, err := engine.NewGotIdFromInt(10)
	assert.NoError(t, err)

	err = store.AppendNote(*gotId, makeBlock("first block"))
	assert.NoError(t, err)

	result, err := store.LongFormNotesFor(*gotId)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Blocks))
	assert.Equal(t, "first block", result.Blocks[0].Content)

	// First note's NoteId should be FirstNoteId().Next()
	firstNoteId := engine.FirstNoteId()
	expectedNoteId := firstNoteId.Next()
	assert.Equal(t, expectedNoteId.IntValue, result.Blocks[0].Id.NoteId.IntValue)
}

func TestLongFormStore_NoteIdIncrementsForSameGotId(t *testing.T) {
	store := buildTestLongFormStore(t)

	gotId, err := engine.NewGotIdFromInt(20)
	assert.NoError(t, err)

	assert.NoError(t, store.AppendNote(*gotId, makeBlock("block 1")))
	assert.NoError(t, store.AppendNote(*gotId, makeBlock("block 2")))
	assert.NoError(t, store.AppendNote(*gotId, makeBlock("block 3")))

	result, err := store.LongFormNotesFor(*gotId)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 3, len(result.Blocks))

	sorted := sortedBlocks(result.Blocks)
	firstNoteId := engine.FirstNoteId()

	// NoteIds should start at FirstNoteId+1 and increment by 1 each time
	for i, block := range sorted {
		expected := firstNoteId.IntValue + int64(i+1)
		assert.Equal(t, expected, block.Id.NoteId.IntValue, "block %d should have NoteId %d", i, expected)
	}
}

func TestLongFormStore_NoteIdsAreIndependentPerGotId(t *testing.T) {
	store := buildTestLongFormStore(t)

	gotId1, err := engine.NewGotIdFromInt(30)
	assert.NoError(t, err)
	gotId2, err := engine.NewGotIdFromInt(40)
	assert.NoError(t, err)

	assert.NoError(t, store.AppendNote(*gotId1, makeBlock("id1 block 1")))
	assert.NoError(t, store.AppendNote(*gotId1, makeBlock("id1 block 2")))
	assert.NoError(t, store.AppendNote(*gotId2, makeBlock("id2 block 1")))

	result1, err := store.LongFormNotesFor(*gotId1)
	assert.NoError(t, err)
	assert.NotNil(t, result1)
	assert.Equal(t, 2, len(result1.Blocks))

	result2, err := store.LongFormNotesFor(*gotId2)
	assert.NoError(t, err)
	assert.NotNil(t, result2)
	assert.Equal(t, 1, len(result2.Blocks))

	// Each gotId's sequence starts independently from FirstNoteId
	firstNoteId := engine.FirstNoteId()

	sorted1 := sortedBlocks(result1.Blocks)
	assert.Equal(t, firstNoteId.IntValue+1, sorted1[0].Id.NoteId.IntValue)
	assert.Equal(t, firstNoteId.IntValue+2, sorted1[1].Id.NoteId.IntValue)

	sorted2 := sortedBlocks(result2.Blocks)
	assert.Equal(t, firstNoteId.IntValue+1, sorted2[0].Id.NoteId.IntValue)
}

func TestLongFormStore_GotIdIsPreservedInBlock(t *testing.T) {
	store := buildTestLongFormStore(t)

	gotId, err := engine.NewGotIdFromInt(50)
	assert.NoError(t, err)

	assert.NoError(t, store.AppendNote(*gotId, makeBlock("content")))

	result, err := store.LongFormNotesFor(*gotId)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Blocks))
	assert.Equal(t, gotId.AasciValue, result.Blocks[0].Id.GotId.AasciValue)
}

func TestLongFormStore_LongFormForMany(t *testing.T) {
	store := buildTestLongFormStore(t)

	gotId1, err := engine.NewGotIdFromInt(60)
	assert.NoError(t, err)
	gotId2, err := engine.NewGotIdFromInt(70)
	assert.NoError(t, err)

	assert.NoError(t, store.AppendNote(*gotId1, makeBlock("id1 content")))
	assert.NoError(t, store.AppendNote(*gotId1, makeBlock("id1 content 2")))
	assert.NoError(t, store.AppendNote(*gotId2, makeBlock("id2 content")))

	results, err := store.LongFormForMany([]engine.GotId{*gotId1, *gotId2})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(results))
	assert.Equal(t, 2, len(results[*gotId1].Blocks))
	assert.Equal(t, 1, len(results[*gotId2].Blocks))
}

func TestLongFormStore_LongFormForMany_MissingIdNotInResult(t *testing.T) {
	store := buildTestLongFormStore(t)

	gotId1, err := engine.NewGotIdFromInt(80)
	assert.NoError(t, err)
	gotId2, err := engine.NewGotIdFromInt(90)
	assert.NoError(t, err)

	assert.NoError(t, store.AppendNote(*gotId1, makeBlock("only id1 has notes")))

	results, err := store.LongFormForMany([]engine.GotId{*gotId1, *gotId2})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, 1, len(results[*gotId1].Blocks))
}

func TestLongFormStore_RemoveAllItems(t *testing.T) {
	store := buildTestLongFormStore(t)

	gotId, err := engine.NewGotIdFromInt(100)
	assert.NoError(t, err)

	assert.NoError(t, store.AppendNote(*gotId, makeBlock("to be removed")))
	assert.NoError(t, store.AppendNote(*gotId, makeBlock("also removed")))

	result, err := store.LongFormNotesFor(*gotId)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result.Blocks))

	err = store.RemoveAllItemsFromLongStoreUnder(*gotId)
	assert.NoError(t, err)

	result, err = store.LongFormNotesFor(*gotId)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestLongFormStore_AppendAfterRemoveRestartSequence(t *testing.T) {
	store := buildTestLongFormStore(t)

	gotId, err := engine.NewGotIdFromInt(110)
	assert.NoError(t, err)

	assert.NoError(t, store.AppendNote(*gotId, makeBlock("original")))
	assert.NoError(t, store.RemoveAllItemsFromLongStoreUnder(*gotId))

	assert.NoError(t, store.AppendNote(*gotId, makeBlock("after removal")))

	result, err := store.LongFormNotesFor(*gotId)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Blocks))

	// After removal and re-append, NoteId sequence restarts from FirstNoteId+1
	firstNoteId := engine.FirstNoteId()
	assert.Equal(t, firstNoteId.IntValue+1, result.Blocks[0].Id.NoteId.IntValue)
}
