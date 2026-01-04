package bullet_engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/containers"
	"vixac.com/got/engine"
)

func buildTestListStore(t *testing.T, meshName string) *BulletListStore {
	mesh, err := bullet_stl.NewBulletMesh(
		BuildTestClient(),
		42,
		meshName,
		">",
		"<",
	)
	assert.NoError(t, err)

	return &BulletListStore{
		SubjectSeparator: ":",
		Mesh:             mesh,
	}
}

func TestListStore_AddAndFetchMembers(t *testing.T) {
	store := buildTestListStore(t, "liststore_basic")

	list := ListId{Id: 1}
	alice, err := engine.NewGotId("alice")
	assert.NoError(t, err)

	err = store.AddItem(*alice, list)
	assert.NoError(t, err)

	res, err := store.FetchListMembers(list)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	lookup, ok := res.Ids[list]
	assert.True(t, ok)
	assert.Equal(t, 1, len(lookup.Ids))
	assert.Equal(t, "alice", lookup.Ids[0].AasciValue)
}

func TestListStore_MultipleItemsSameList(t *testing.T) {
	store := buildTestListStore(t, "liststore_multi_items")

	list := ListId{Id: 2}

	alice, _ := engine.NewGotId("alice")
	bob, _ := engine.NewGotId("bob")

	assert.NoError(t, store.AddItem(*alice, list))
	assert.NoError(t, store.AddItem(*bob, list))

	res, err := store.FetchListMembers(list)
	assert.NoError(t, err)

	lookup := res.Ids[list]
	assert.Equal(t, 2, len(lookup.Ids))

	found := map[string]bool{}
	for _, id := range lookup.Ids {
		found[id.AasciValue] = true
	}

	assert.True(t, found["alice"])
	assert.True(t, found["bob"])
}

func TestListStore_FetchListsContaining(t *testing.T) {
	store := buildTestListStore(t, "liststore_reverse")

	alice, _ := engine.NewGotId("alice")

	listA := ListId{Id: 10}
	listB := ListId{Id: 20}

	assert.NoError(t, store.AddItem(*alice, listA))
	assert.NoError(t, store.AddItem(*alice, listB))

	res, err := store.FetchListsContaining(*alice)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.Equal(t, 2, len(res.Lists))

	found := map[int32]bool{}
	for _, l := range res.Lists {
		found[l.Id] = true
	}

	assert.True(t, found[10])
	assert.True(t, found[20])
}
func TestListStore_ListIsolation(t *testing.T) {
	store := buildTestListStore(t, "liststore_isolation")

	listA := ListId{Id: 100}
	listB := ListId{Id: 200}

	alice, _ := engine.NewGotId("alice")
	bob, _ := engine.NewGotId("bob")

	assert.NoError(t, store.AddItem(*alice, listA))
	assert.NoError(t, store.AddItem(*bob, listB))

	resA, err := store.FetchListMembers(listA)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resA.Ids[listA].Ids))
	assert.Equal(t, "alice", resA.Ids[listA].Ids[0].AasciValue)

	resB, err := store.FetchListMembers(listB)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resB.Ids[listB].Ids))
	assert.Equal(t, "bob", resB.Ids[listB].Ids[0].AasciValue)
}

func TestListStore_EmptyList(t *testing.T) {
	store := buildTestListStore(t, "liststore_empty")

	list := ListId{Id: 999}

	res, err := store.FetchListMembers(list)
	assert.NoError(t, err)
	assert.Nil(t, res)

}
