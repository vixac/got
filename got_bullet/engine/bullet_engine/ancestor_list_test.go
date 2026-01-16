package bullet_engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/containers"
	"vixac.com/got/engine"
)

func buildTestAncestorList(t *testing.T, meshName string) *BulletAncestorList {
	mesh, err := bullet_stl.NewBulletMesh(
		BuildTestClient(),
		42,
		meshName,
		">",
		"<",
	)
	assert.NoError(t, err)

	return &BulletAncestorList{
		SubjectSeparator: ":",
		Client:           BuildTestClient(),
		Mesh:             mesh,
	}
}

func TestAncestorList_AddItemAtRoot(t *testing.T) {
	list := buildTestAncestorList(t, "ancestor_add_root")

	alice, err := engine.NewGotId("alice")
	assert.NoError(t, err)

	// Add alice under root (nil parent)
	ancestry, err := list.AddItem(*alice, nil)
	assert.NoError(t, err)
	assert.Nil(t, ancestry) // No ancestors when adding under root

	// Fetch ancestors of alice
	res, err := list.FetchAncestorsOf(*alice)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 1, len(res.Ids))
	assert.Equal(t, TheRootNode.Value, res.Ids[0].AasciValue)
}

func TestAncestorList_AddItemUnderParent(t *testing.T) {
	list := buildTestAncestorList(t, "ancestor_add_under_parent")

	alice, _ := engine.NewGotId("alice")
	bob, _ := engine.NewGotId("bob")

	// Add alice under root
	_, err := list.AddItem(*alice, nil)
	assert.NoError(t, err)

	// Add bob under alice
	ancestry, err := list.AddItem(*bob, alice)
	assert.NoError(t, err)
	assert.NotNil(t, ancestry)
	assert.Equal(t, 2, len(ancestry.Ids))
	assert.Equal(t, TheRootNode.Value, ancestry.Ids[0].AasciValue)
	assert.Equal(t, "alice", ancestry.Ids[1].AasciValue)

	// Fetch ancestors of bob
	res, err := list.FetchAncestorsOf(*bob)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 2, len(res.Ids))
	assert.Equal(t, TheRootNode.Value, res.Ids[0].AasciValue)
	assert.Equal(t, "alice", res.Ids[1].AasciValue)
}

func TestAncestorList_AddItemDeepHierarchy(t *testing.T) {
	list := buildTestAncestorList(t, "ancestor_deep_hierarchy")

	alice, _ := engine.NewGotId("alice")
	bob, _ := engine.NewGotId("bob")
	charlie, _ := engine.NewGotId("c")

	// Build a hierarchy: root -> alice -> bob -> charlie
	list.AddItem(*alice, nil)
	list.AddItem(*bob, alice)
	list.AddItem(*charlie, bob)

	// Fetch ancestors of charlie
	res, err := list.FetchAncestorsOf(*charlie)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 3, len(res.Ids))
	assert.Equal(t, TheRootNode.Value, res.Ids[0].AasciValue)
	assert.Equal(t, "alice", res.Ids[1].AasciValue)
	assert.Equal(t, "bob", res.Ids[2].AasciValue)
}

func TestAncestorList_AddDuplicateItem_Error(t *testing.T) {
	list := buildTestAncestorList(t, "ancestor_duplicate")

	alice, _ := engine.NewGotId("alice")

	_, err := list.AddItem(*alice, nil)
	assert.NoError(t, err)

	// Try to add alice again
	_, err = list.AddItem(*alice, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "existing id")
}

func TestAncestorList_AddRootNode_Error(t *testing.T) {
	list := buildTestAncestorList(t, "ancestor_add_root_node")

	rootId, _ := engine.NewGotId(TheRootNode.Value)

	_, err := list.AddItem(*rootId, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "root node")
}

func TestAncestorList_FetchAncestorsOfMany_Empty(t *testing.T) {
	list := buildTestAncestorList(t, "ancestor_many_empty")

	res, err := list.FetchAncestorsOfMany([]engine.GotId{})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 0, len(res.Ids))
}

func TestAncestorList_FetchAncestorsOfMany_SingleItem(t *testing.T) {
	list := buildTestAncestorList(t, "ancestor_many_single")

	alice, _ := engine.NewGotId("alice")
	bob, _ := engine.NewGotId("bob")

	list.AddItem(*alice, nil)
	list.AddItem(*bob, alice)

	res, err := list.FetchAncestorsOfMany([]engine.GotId{*bob})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 1, len(res.Ids))

	bobAncestors, ok := res.Ids[*bob]
	assert.True(t, ok)
	assert.Equal(t, 2, len(bobAncestors.Ids))
	assert.Equal(t, TheRootNode.Value, bobAncestors.Ids[0].AasciValue)
	assert.Equal(t, "alice", bobAncestors.Ids[1].AasciValue)
}

func TestAncestorList_FetchAncestorsOfMany_MultipleItems(t *testing.T) {
	list := buildTestAncestorList(t, "ancestor_many_multi")

	alice, _ := engine.NewGotId("alice")
	bob, _ := engine.NewGotId("bob")
	charlie, _ := engine.NewGotId("c")

	// Build hierarchy: root -> alice -> bob, root -> alice -> charlie
	list.AddItem(*alice, nil)
	list.AddItem(*bob, alice)
	list.AddItem(*charlie, alice)

	res, err := list.FetchAncestorsOfMany([]engine.GotId{*bob, *charlie})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 2, len(res.Ids))

	// Check bob's ancestors
	bobAncestors, ok := res.Ids[*bob]
	assert.True(t, ok)
	assert.Equal(t, 2, len(bobAncestors.Ids))
	assert.Equal(t, TheRootNode.Value, bobAncestors.Ids[0].AasciValue)
	assert.Equal(t, "alice", bobAncestors.Ids[1].AasciValue)

	// Check charlie's ancestors
	charlieAncestors, ok := res.Ids[*charlie]
	assert.True(t, ok)
	assert.Equal(t, 2, len(charlieAncestors.Ids))
	assert.Equal(t, TheRootNode.Value, charlieAncestors.Ids[0].AasciValue)
	assert.Equal(t, "alice", charlieAncestors.Ids[1].AasciValue)
}

func TestAncestorList_FetchAncestorsOfMany_DifferentDepths(t *testing.T) {
	list := buildTestAncestorList(t, "ancestor_many_depths")

	alice, _ := engine.NewGotId("alice")
	bob, _ := engine.NewGotId("bob")
	charlie, _ := engine.NewGotId("c")

	// Build hierarchy: root -> alice, root -> alice -> bob -> charlie
	list.AddItem(*alice, nil)
	list.AddItem(*bob, alice)
	list.AddItem(*charlie, bob)

	res, err := list.FetchAncestorsOfMany([]engine.GotId{*alice, *charlie})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 2, len(res.Ids))

	// Check alice's ancestors (just root)
	aliceAncestors, ok := res.Ids[*alice]
	assert.True(t, ok)
	assert.Equal(t, 1, len(aliceAncestors.Ids))
	assert.Equal(t, TheRootNode.Value, aliceAncestors.Ids[0].AasciValue)

	// Check charlie's ancestors (root -> alice -> bob)
	charlieAncestors, ok := res.Ids[*charlie]
	assert.True(t, ok)
	assert.Equal(t, 3, len(charlieAncestors.Ids))
	assert.Equal(t, TheRootNode.Value, charlieAncestors.Ids[0].AasciValue)
	assert.Equal(t, "alice", charlieAncestors.Ids[1].AasciValue)
	assert.Equal(t, "bob", charlieAncestors.Ids[2].AasciValue)
}

func TestAncestorList_FetchImmediatelyUnder_Root(t *testing.T) {
	list := buildTestAncestorList(t, "ancestor_fetch_under_root")

	alice, _ := engine.NewGotId("alice")
	bob, _ := engine.NewGotId("bob")

	list.AddItem(*alice, nil)
	list.AddItem(*bob, nil)

	rootId, _ := engine.NewGotId(TheRootNode.Value)
	res, err := list.FetchImmediatelyUnder(*rootId)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 2, len(res.Ids))

	_, hasAlice := res.Ids["alice"]
	_, hasBob := res.Ids["bob"]
	assert.True(t, hasAlice)
	assert.True(t, hasBob)
}

func TestAncestorList_FetchImmediatelyUnder_Parent(t *testing.T) {
	list := buildTestAncestorList(t, "ancestor_fetch_under_parent")

	alice, _ := engine.NewGotId("alice")
	bob, _ := engine.NewGotId("bob")
	charlie, _ := engine.NewGotId("c")

	// root -> alice -> bob, root -> alice -> charlie
	list.AddItem(*alice, nil)
	list.AddItem(*bob, alice)
	list.AddItem(*charlie, alice)

	res, err := list.FetchImmediatelyUnder(*alice)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 2, len(res.Ids))

	_, hasBob := res.Ids["bob"]
	_, hasCharlie := res.Ids["c"]
	assert.True(t, hasBob)
	assert.True(t, hasCharlie)
}

func TestAncestorList_FetchImmediatelyUnder_Empty(t *testing.T) {
	list := buildTestAncestorList(t, "ancestor_fetch_under_empty")

	alice, _ := engine.NewGotId("alice")
	list.AddItem(*alice, nil)

	// alice has no children
	res, err := list.FetchImmediatelyUnder(*alice)
	assert.NoError(t, err)
	assert.Nil(t, res)
}

func TestAncestorList_RemoveItem(t *testing.T) {
	list := buildTestAncestorList(t, "ancestor_remove")

	alice, _ := engine.NewGotId("alice")

	list.AddItem(*alice, nil)

	// Verify alice exists
	res, err := list.FetchAncestorsOf(*alice)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	// Remove alice
	err = list.RemoveItem(*alice)
	assert.NoError(t, err)

	// Verify alice no longer exists
	res, err = list.FetchAncestorsOf(*alice)
	assert.NoError(t, err)
	assert.Nil(t, res)
}

func TestAncestorList_FetchAncestors_NonExistent(t *testing.T) {
	list := buildTestAncestorList(t, "ancestor_nonexistent")

	alice, _ := engine.NewGotId("alice")

	res, err := list.FetchAncestorsOf(*alice)
	assert.NoError(t, err)
	assert.Nil(t, res)
}

func TestAncestorList_MoveItem_LeafToDifferentParent(t *testing.T) {
	list := buildTestAncestorList(t, "ancestor_move_leaf")

	alice, _ := engine.NewGotId("alice")
	bob, _ := engine.NewGotId("bob")
	charlie, _ := engine.NewGotId("c")

	// Build hierarchy: root -> alice -> charlie, root -> bob
	list.AddItem(*alice, nil)
	list.AddItem(*bob, nil)
	list.AddItem(*charlie, alice)

	// Verify charlie is under alice
	ancestors, err := list.FetchAncestorsOf(*charlie)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(ancestors.Ids))
	assert.Equal(t, "alice", ancestors.Ids[1].AasciValue)

	// Move charlie from under alice to under bob
	result, err := list.MoveItem(*charlie, bob)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify old ancestry
	assert.NotNil(t, result.OldAncestry)
	assert.Equal(t, 2, len(result.OldAncestry.Ids))
	assert.Equal(t, TheRootNode.Value, result.OldAncestry.Ids[0].AasciValue)
	assert.Equal(t, "alice", result.OldAncestry.Ids[1].AasciValue)

	// Verify new ancestry
	assert.NotNil(t, result.NewAncestry)
	assert.Equal(t, 2, len(result.NewAncestry.Ids))
	assert.Equal(t, TheRootNode.Value, result.NewAncestry.Ids[0].AasciValue)
	assert.Equal(t, "bob", result.NewAncestry.Ids[1].AasciValue)

	// Verify charlie is now under bob
	ancestors, err = list.FetchAncestorsOf(*charlie)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(ancestors.Ids))
	assert.Equal(t, "bob", ancestors.Ids[1].AasciValue)

	// Verify charlie is no longer under alice
	aliceChildren, err := list.FetchImmediatelyUnder(*alice)
	assert.NoError(t, err)
	assert.Nil(t, aliceChildren)

	// Verify charlie is now in bob's children
	bobChildren, err := list.FetchImmediatelyUnder(*bob)
	assert.NoError(t, err)
	assert.NotNil(t, bobChildren)
	_, hasCharlie := bobChildren.Ids["c"]
	assert.True(t, hasCharlie)
}

func TestAncestorList_MoveItem_LeafToRoot(t *testing.T) {
	list := buildTestAncestorList(t, "ancestor_move_to_root")

	alice, _ := engine.NewGotId("alice")
	bob, _ := engine.NewGotId("bob")

	// Build hierarchy: root -> alice -> bob
	list.AddItem(*alice, nil)
	list.AddItem(*bob, alice)

	// Verify bob is under alice
	ancestors, err := list.FetchAncestorsOf(*bob)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(ancestors.Ids))

	// Move bob to root (nil parent)
	result, err := list.MoveItem(*bob, nil)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify old ancestry
	assert.Equal(t, 2, len(result.OldAncestry.Ids))
	assert.Equal(t, "alice", result.OldAncestry.Ids[1].AasciValue)

	// Verify new ancestry (just root)
	assert.Equal(t, 1, len(result.NewAncestry.Ids))
	assert.Equal(t, TheRootNode.Value, result.NewAncestry.Ids[0].AasciValue)

	// Verify bob is now directly under root
	ancestors, err = list.FetchAncestorsOf(*bob)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(ancestors.Ids))
	assert.Equal(t, TheRootNode.Value, ancestors.Ids[0].AasciValue)
}

func TestAncestorList_MoveItem_FromRootToParent(t *testing.T) {
	list := buildTestAncestorList(t, "ancestor_move_from_root")

	alice, _ := engine.NewGotId("alice")
	bob, _ := engine.NewGotId("bob")

	// Build hierarchy: root -> alice, root -> bob
	list.AddItem(*alice, nil)
	list.AddItem(*bob, nil)

	// Move bob under alice
	result, err := list.MoveItem(*bob, alice)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify old ancestry (just root)
	assert.Equal(t, 1, len(result.OldAncestry.Ids))
	assert.Equal(t, TheRootNode.Value, result.OldAncestry.Ids[0].AasciValue)

	// Verify new ancestry
	assert.Equal(t, 2, len(result.NewAncestry.Ids))
	assert.Equal(t, TheRootNode.Value, result.NewAncestry.Ids[0].AasciValue)
	assert.Equal(t, "alice", result.NewAncestry.Ids[1].AasciValue)

	// Verify bob is now under alice
	ancestors, err := list.FetchAncestorsOf(*bob)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(ancestors.Ids))
	assert.Equal(t, "alice", ancestors.Ids[1].AasciValue)
}

func TestAncestorList_MoveItem_NonLeaf_Error(t *testing.T) {
	list := buildTestAncestorList(t, "ancestor_move_nonleaf")

	alice, _ := engine.NewGotId("alice")
	bob, _ := engine.NewGotId("bob")
	charlie, _ := engine.NewGotId("c")

	// Build hierarchy: root -> alice -> bob, root -> charlie
	list.AddItem(*alice, nil)
	list.AddItem(*bob, alice)
	list.AddItem(*charlie, nil)

	// Try to move alice (which has bob as child) - should fail
	result, err := list.MoveItem(*alice, charlie)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "cannot move non-leaf node")
}

func TestAncestorList_MoveItem_NonExistent_Error(t *testing.T) {
	list := buildTestAncestorList(t, "ancestor_move_nonexistent")

	alice, _ := engine.NewGotId("alice")
	bob, _ := engine.NewGotId("bob")

	list.AddItem(*alice, nil)

	// Try to move bob (which doesn't exist)
	result, err := list.MoveItem(*bob, alice)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestAncestorList_MoveItem_RootNode_Error(t *testing.T) {
	list := buildTestAncestorList(t, "ancestor_move_root")

	alice, _ := engine.NewGotId("alice")
	rootId, _ := engine.NewGotId(TheRootNode.Value)

	list.AddItem(*alice, nil)

	// Try to move the root node - should fail
	result, err := list.MoveItem(*rootId, alice)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "root node")
}

func TestAncestorList_MoveItem_DeepHierarchy(t *testing.T) {
	list := buildTestAncestorList(t, "ancestor_move_deep")

	alice, _ := engine.NewGotId("alice")
	bob, _ := engine.NewGotId("bob")
	charlie, _ := engine.NewGotId("c")
	dave, _ := engine.NewGotId("dave")

	// Build hierarchy: root -> alice -> bob -> charlie, root -> dave
	list.AddItem(*alice, nil)
	list.AddItem(*bob, alice)
	list.AddItem(*charlie, bob)
	list.AddItem(*dave, nil)

	// Verify charlie's current ancestry (root -> alice -> bob)
	ancestors, err := list.FetchAncestorsOf(*charlie)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(ancestors.Ids))

	// Move charlie from bob to dave
	result, err := list.MoveItem(*charlie, dave)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify old ancestry
	assert.Equal(t, 3, len(result.OldAncestry.Ids))
	assert.Equal(t, TheRootNode.Value, result.OldAncestry.Ids[0].AasciValue)
	assert.Equal(t, "alice", result.OldAncestry.Ids[1].AasciValue)
	assert.Equal(t, "bob", result.OldAncestry.Ids[2].AasciValue)

	// Verify new ancestry
	assert.Equal(t, 2, len(result.NewAncestry.Ids))
	assert.Equal(t, TheRootNode.Value, result.NewAncestry.Ids[0].AasciValue)
	assert.Equal(t, "dave", result.NewAncestry.Ids[1].AasciValue)

	// Verify charlie is now under dave
	ancestors, err = list.FetchAncestorsOf(*charlie)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(ancestors.Ids))
	assert.Equal(t, "dave", ancestors.Ids[1].AasciValue)
}
