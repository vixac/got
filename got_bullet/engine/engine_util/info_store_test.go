package engine_util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/containers"
	"vixac.com/got/engine"
)

func buildTestBuckStore(t *testing.T) *BuckStore {
	client := engine.BuildTestClient()
	coll := bullet_stl.NewBulletCollection(0, client, client)
	return &BuckStore{
		Codec:      &JSONCodec[BuckInfo]{},
		Collection: coll,
	}
}

func makeId(t *testing.T, i int32) engine.GotId {
	id, err := engine.NewGotIdFromInt(i)
	assert.NoError(t, err)
	return *id
}

func TestBuckStore_UpsertAndFetch(t *testing.T) {
	store := buildTestBuckStore(t)
	id := makeId(t, 42)

	info := BuckInfo{Title: "Hello World"}
	err := store.UpsertInfo(id, info)
	assert.NoError(t, err)

	resp, err := store.InfoForMany([]engine.GotId{id})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, len(resp.InfoMap))
	assert.Equal(t, 0, len(resp.Missing))
	assert.Equal(t, "Hello World", resp.InfoMap[id].Title)
}

func TestBuckStore_UpdateExisting(t *testing.T) {
	store := buildTestBuckStore(t)
	id := makeId(t, 1)

	err := store.UpsertInfo(id, BuckInfo{Title: "First"})
	assert.NoError(t, err)

	err = store.UpsertInfo(id, BuckInfo{Title: "Updated"})
	assert.NoError(t, err)

	resp, err := store.InfoForMany([]engine.GotId{id})
	assert.NoError(t, err)
	assert.Equal(t, "Updated", resp.InfoMap[id].Title)
}

func TestBuckStore_InfoForMany(t *testing.T) {
	store := buildTestBuckStore(t)
	id10 := makeId(t, 10)
	id20 := makeId(t, 20)
	id30 := makeId(t, 30)

	assert.NoError(t, store.UpsertInfo(id10, BuckInfo{Title: "Ten"}))
	assert.NoError(t, store.UpsertInfo(id20, BuckInfo{Title: "Twenty"}))
	assert.NoError(t, store.UpsertInfo(id30, BuckInfo{Title: "Thirty"}))

	resp, err := store.InfoForMany([]engine.GotId{id10, id20, id30})
	assert.NoError(t, err)
	assert.Equal(t, 3, len(resp.InfoMap))
	assert.Equal(t, 0, len(resp.Missing))
	assert.Equal(t, "Ten", resp.InfoMap[id10].Title)
	assert.Equal(t, "Twenty", resp.InfoMap[id20].Title)
	assert.Equal(t, "Thirty", resp.InfoMap[id30].Title)
}

func TestBuckStore_InfoForMany_Missing(t *testing.T) {
	store := buildTestBuckStore(t)
	id := makeId(t, 999)

	resp, err := store.InfoForMany([]engine.GotId{id})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 0, len(resp.InfoMap))
	assert.Equal(t, 1, len(resp.Missing))
	assert.Equal(t, id, resp.Missing[0])
}

func TestBuckStore_InfoForMany_PartialMatch(t *testing.T) {
	store := buildTestBuckStore(t)
	id100 := makeId(t, 100)
	id200 := makeId(t, 200)

	assert.NoError(t, store.UpsertInfo(id100, BuckInfo{Title: "Hundred"}))

	resp, err := store.InfoForMany([]engine.GotId{id100, id200})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.InfoMap))
	assert.Equal(t, 1, len(resp.Missing))
	assert.Equal(t, "Hundred", resp.InfoMap[id100].Title)
	assert.Equal(t, id200, resp.Missing[0])
}

func TestBuckStore_DeleteInfoMany(t *testing.T) {
	store := buildTestBuckStore(t)
	id5 := makeId(t, 5)
	id6 := makeId(t, 6)

	assert.NoError(t, store.UpsertInfo(id5, BuckInfo{Title: "Five"}))
	assert.NoError(t, store.UpsertInfo(id6, BuckInfo{Title: "Six"}))

	err := store.DeleteInfoMany([]engine.GotId{id5, id6})
	assert.NoError(t, err)

	resp, err := store.InfoForMany([]engine.GotId{id5, id6})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(resp.InfoMap))
	assert.Equal(t, 2, len(resp.Missing))
}

func TestBuckStore_MultipleItemsAreIndependent(t *testing.T) {
	store := buildTestBuckStore(t)
	id7 := makeId(t, 7)
	id8 := makeId(t, 8)

	assert.NoError(t, store.UpsertInfo(id7, BuckInfo{Title: "Seven"}))
	assert.NoError(t, store.UpsertInfo(id8, BuckInfo{Title: "Eight"}))
	assert.NoError(t, store.UpsertInfo(id7, BuckInfo{Title: "Seven Updated"}))

	resp, err := store.InfoForMany([]engine.GotId{id7, id8})
	assert.NoError(t, err)
	assert.Equal(t, "Seven Updated", resp.InfoMap[id7].Title)
	assert.Equal(t, "Eight", resp.InfoMap[id8].Title)
}

func TestBuckStore_PreservesAllFields(t *testing.T) {
	store := buildTestBuckStore(t)
	id := makeId(t, 55)

	now := engine.DateTime{Millis: 1234567890}
	work := "work"
	tags := []engine.Tag{{Identifier: &work}}
	flags := map[string]bool{"collapsed": true}
	info := BuckInfo{
		Title:       "Full Info",
		CreatedDate: now,
		Tags:        tags,
		Flags:       flags,
	}

	assert.NoError(t, store.UpsertInfo(id, info))

	resp, err := store.InfoForMany([]engine.GotId{id})
	assert.NoError(t, err)
	got := resp.InfoMap[id]
	assert.Equal(t, "Full Info", got.Title)
	assert.Equal(t, now, got.CreatedDate)
	assert.Equal(t, tags, got.Tags)
	assert.Equal(t, flags, got.Flags)
}
