package engine_util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/containers"
	"vixac.com/got/engine"
)

func buildTestTitleStore(t *testing.T) *BulletTitleStore {
	client := engine.BuildTestClient()
	coll := bullet_stl.NewBulletCollection(0, client, client)
	return &BulletTitleStore{Collection: coll}
}

func TestTitleStore_UpsertAndFetch(t *testing.T) {
	store := buildTestTitleStore(t)

	err := store.UpsertItem(42, "Hello World")
	assert.NoError(t, err)

	title, err := store.TitleFor(42)
	assert.NoError(t, err)
	assert.NotNil(t, title)
	assert.Equal(t, "Hello World", *title)
}

func TestTitleStore_UpdateExisting(t *testing.T) {
	store := buildTestTitleStore(t)

	err := store.UpsertItem(1, "First Title")
	assert.NoError(t, err)

	err = store.UpsertItem(1, "Updated Title")
	assert.NoError(t, err)

	title, err := store.TitleFor(1)
	assert.NoError(t, err)
	assert.NotNil(t, title)
	assert.Equal(t, "Updated Title", *title)
}

func TestTitleStore_TitleForMany(t *testing.T) {
	store := buildTestTitleStore(t)

	assert.NoError(t, store.UpsertItem(10, "Ten"))
	assert.NoError(t, store.UpsertItem(20, "Twenty"))
	assert.NoError(t, store.UpsertItem(30, "Thirty"))

	result, err := store.TitleForMany([]int32{10, 20, 30})
	assert.NoError(t, err)
	assert.Equal(t, 3, len(result))
	assert.Equal(t, "Ten", result[10])
	assert.Equal(t, "Twenty", result[20])
	assert.Equal(t, "Thirty", result[30])
}

func TestTitleStore_TitleFor_Missing(t *testing.T) {
	store := buildTestTitleStore(t)

	title, err := store.TitleFor(999)
	assert.NoError(t, err)
	assert.Nil(t, title)
}

func TestTitleStore_RemoveItem(t *testing.T) {
	store := buildTestTitleStore(t)

	assert.NoError(t, store.UpsertItem(5, "To Remove"))

	title, err := store.TitleFor(5)
	assert.NoError(t, err)
	assert.NotNil(t, title)

	err = store.RemoveItem(5)
	assert.NoError(t, err)

	title, err = store.TitleFor(5)
	assert.NoError(t, err)
	assert.Nil(t, title)
}

func TestTitleStore_TitleForMany_PartialMatch(t *testing.T) {
	store := buildTestTitleStore(t)

	assert.NoError(t, store.UpsertItem(100, "Hundred"))

	result, err := store.TitleForMany([]int32{100, 200})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "Hundred", result[100])
}

func TestTitleStore_MultipleItemsAreIndependent(t *testing.T) {
	store := buildTestTitleStore(t)

	assert.NoError(t, store.UpsertItem(7, "Seven"))
	assert.NoError(t, store.UpsertItem(8, "Eight"))

	assert.NoError(t, store.UpsertItem(7, "Seven Updated"))

	title7, err := store.TitleFor(7)
	assert.NoError(t, err)
	assert.Equal(t, "Seven Updated", *title7)

	title8, err := store.TitleFor(8)
	assert.NoError(t, err)
	assert.Equal(t, "Eight", *title8)
}
