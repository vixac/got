package bullet_engine

import (
	"fmt"
	"testing"

	"gotest.tools/assert"
	"vixac.com/got/engine"
)

func TestAggregatorJustAliceAndBob(t *testing.T) {

	store := MakeMockSummaryStore()
	agg, err := NewAggregator(&store)
	assert.NilError(t, err)

	var aliceId = engine.SummaryId(10)
	var bob = engine.SummaryId(11)

	var activeState = engine.GotState(engine.Active)
	var completeState = engine.GotState(engine.Complete)

	//create alice -> bob
	err = agg.ItemAdded(AddItemEvent{
		Id:       aliceId,
		State:    engine.Active,
		Ancestry: []engine.SummaryId{},
		Deadline: nil,
	})
	assert.NilError(t, err)

	assert.Assert(t, store.aggs != nil)

	assert.Equal(t, len(store.aggs), 1)

	agg.ItemAdded(AddItemEvent{
		Id:       bob,
		State:    engine.Active,
		Ancestry: []engine.SummaryId{aliceId},
		Deadline: nil,
	})

	//fetch both
	assert.Equal(t, len(store.aggs), 2)
	fetchBoth, err := agg.summaryStore.Fetch([]engine.SummaryId{
		aliceId,
		bob,
	})
	assert.NilError(t, err)
	fetchedAlice, ok := fetchBoth[aliceId]
	assert.Equal(t, ok, true)
	fetchedBob, ok := fetchBoth[bob]
	//fmt.Printf("k, bobcounts a, %+v\n", fetchedBob.Counts)
	assert.Equal(t, ok, true)
	assert.Assert(t, fetchedAlice.Counts != nil)
	assert.Assert(t, fetchedBob.Counts == nil)

	assert.Equal(t, fetchedAlice.Counts.Active, 1)
	assert.Equal(t, fetchedAlice.Counts.Complete, 0)
	assert.Equal(t, fetchedAlice.Counts.Notes, 0)
	var nilState *engine.GotState = nil
	assert.Equal(t, fetchedAlice.State, nilState)

	//complete bob, expect alice to become active
	err = agg.ItemStateChanged(StateChangeEvent{
		Id:       bob,
		OldState: engine.Active,
		NewState: &completeState,
		Ancestry: []engine.SummaryId{aliceId},
	})
	assert.NilError(t, err)

	//refetch
	fetchBoth, err = agg.summaryStore.Fetch([]engine.SummaryId{
		aliceId,
		bob,
	})

	assert.NilError(t, err)
	fetchedAlice, ok = fetchBoth[aliceId]
	assert.Equal(t, ok, true)
	fetchedBob, ok = fetchBoth[bob]
	assert.Equal(t, ok, true)
	assert.Assert(t, fetchedAlice.Counts != nil)
	assert.Assert(t, fetchedBob.Counts == nil)

	assert.Equal(t, fetchedAlice.Counts.Active, 0)
	assert.Equal(t, fetchedAlice.Counts.Complete, 1)
	assert.Equal(t, fetchedAlice.Counts.Notes, 0)
	assert.Equal(t, *fetchedAlice.State, activeState)

	//now complete alice:

	//complete alice
	err = agg.ItemStateChanged(StateChangeEvent{
		Id:       aliceId,
		OldState: engine.Active,
		NewState: &completeState,
		Ancestry: []engine.SummaryId{},
	})
	assert.NilError(t, err)

	fetchBoth, err = agg.summaryStore.Fetch([]engine.SummaryId{
		aliceId,
		bob,
	})
	assert.NilError(t, err)

	fetchedAlice, ok = fetchBoth[aliceId]
	assert.Equal(t, ok, true)
	assert.Equal(t, *fetchedAlice.State, completeState)

	assert.Equal(t, fetchedAlice.Counts.Active, 0)
	assert.Equal(t, fetchedAlice.Counts.Complete, 1)
	assert.Equal(t, fetchedAlice.Counts.Notes, 0)
}

//VX:TODO test that notes under a complete group don't render
//VX:TODO summary of the top item

func TestAggregatorTopAliceAndBob(t *testing.T) {

	store := MakeMockSummaryStore()
	agg, err := NewAggregator(&store)
	assert.NilError(t, err)

	var top = engine.SummaryId(9)
	var aliceId = engine.SummaryId(10)
	var bob = engine.SummaryId(11)

	var activeState = engine.GotState(engine.Active)
	var completeState = engine.GotState(engine.Complete)

	err = agg.ItemAdded(AddItemEvent{
		Id:       top,
		State:    engine.Active,
		Ancestry: []engine.SummaryId{},
		Deadline: nil,
	})

	assert.NilError(t, err)

	//create zob ->alice -> bob
	err = agg.ItemAdded(AddItemEvent{
		Id:       aliceId,
		State:    engine.Active,
		Ancestry: []engine.SummaryId{top},
		Deadline: nil,
	})
	assert.NilError(t, err)

	assert.Assert(t, store.aggs != nil)

	assert.Equal(t, len(store.aggs), 2)

	agg.ItemAdded(AddItemEvent{
		Id:       bob,
		State:    engine.Active,
		Ancestry: []engine.SummaryId{top, aliceId},
		Deadline: nil,
	})

	//fetch all
	assert.Equal(t, len(store.aggs), 3)
	fetchAll, err := agg.summaryStore.Fetch([]engine.SummaryId{
		top,
		aliceId,
		bob,
	})
	assert.NilError(t, err)
	fetchedAlice, ok := fetchAll[aliceId]
	assert.Equal(t, ok, true)
	fetchedBob, ok := fetchAll[bob]
	assert.Equal(t, ok, true)

	fetchTop, ok := fetchAll[top]
	assert.Equal(t, ok, true)
	//fmt.Printf("k, bobcounts a, %+v\n", fetchedBob.Counts)
	assert.Equal(t, ok, true)
	assert.Assert(t, fetchTop.Counts != nil)
	assert.Assert(t, fetchedAlice.Counts != nil)
	assert.Assert(t, fetchedBob.Counts == nil)

	assert.Equal(t, fetchTop.Counts.Active, 1)
	assert.Equal(t, fetchTop.Counts.Complete, 0)
	assert.Equal(t, fetchTop.Counts.Notes, 0)

	assert.Equal(t, fetchedAlice.Counts.Active, 1)
	assert.Equal(t, fetchedAlice.Counts.Complete, 0)
	assert.Equal(t, fetchedAlice.Counts.Notes, 0)
	var nilState *engine.GotState = nil
	assert.Equal(t, fetchedAlice.State, nilState)

	//complete bob, expect alice to become active
	err = agg.ItemStateChanged(StateChangeEvent{
		Id:       bob,
		OldState: engine.Active,
		NewState: &completeState,
		Ancestry: []engine.SummaryId{top, aliceId},
	})
	assert.NilError(t, err)

	//refetch
	fetchAll, err = agg.summaryStore.Fetch([]engine.SummaryId{
		top,
		aliceId,
		bob,
	})

	assert.NilError(t, err)

	fetchedAlice, ok = fetchAll[aliceId]
	assert.Equal(t, ok, true)

	fetchedBob, ok = fetchAll[bob]
	assert.Equal(t, ok, true)

	fetchTop, ok = fetchAll[top]
	assert.Equal(t, ok, true)

	assert.Assert(t, fetchedAlice.Counts != nil)
	assert.Assert(t, fetchedBob.Counts == nil)

	assert.Equal(t, fetchedAlice.Counts.Active, 0)
	assert.Equal(t, fetchedAlice.Counts.Complete, 1)
	assert.Equal(t, fetchedAlice.Counts.Notes, 0)
	assert.Equal(t, *fetchedAlice.State, activeState)

	assert.Equal(t, fetchTop.Counts.Active, 1)
	assert.Equal(t, fetchTop.Counts.Complete, 1)
	assert.Equal(t, fetchTop.Counts.Notes, 0)

	//now complete alice:

	//complete alice
	err = agg.ItemStateChanged(StateChangeEvent{
		Id:       aliceId,
		OldState: engine.Active,
		NewState: &completeState,
		Ancestry: []engine.SummaryId{},
	})
	assert.NilError(t, err)

	fetchAll, err = agg.summaryStore.Fetch([]engine.SummaryId{
		top,
		aliceId,
		bob,
	})
	assert.NilError(t, err)

	fetchedAlice, ok = fetchAll[aliceId]
	assert.Equal(t, ok, true)
	assert.Equal(t, *fetchedAlice.State, completeState)

	assert.Equal(t, fetchedAlice.Counts.Active, 0)
	assert.Equal(t, fetchedAlice.Counts.Complete, 1)
	assert.Equal(t, fetchedAlice.Counts.Notes, 0)
}

func TestAggregatorPreservesCompletes(t *testing.T) {

	store := MakeMockSummaryStore()
	agg, err := NewAggregator(&store)
	assert.NilError(t, err)

	var aliceId = engine.SummaryId(10)
	var bob = engine.SummaryId(11)
	var carolId = engine.SummaryId(12)

	var activeState = engine.GotState(engine.Active)
	var completeState = engine.GotState(engine.Complete)
	//create alice -> bob
	err = agg.ItemAdded(AddItemEvent{
		Id:       aliceId,
		State:    engine.Active,
		Ancestry: []engine.SummaryId{},
		Deadline: nil,
	})
	assert.NilError(t, err)

	assert.Assert(t, store.aggs != nil)

	assert.Equal(t, len(store.aggs), 1)

	agg.ItemAdded(AddItemEvent{
		Id:       bob,
		State:    engine.Active,
		Ancestry: []engine.SummaryId{aliceId},
		Deadline: nil,
	})

	//fetch both
	assert.Equal(t, len(store.aggs), 2)
	fetchBoth, err := agg.summaryStore.Fetch([]engine.SummaryId{
		aliceId,
		bob,
	})
	assert.NilError(t, err)
	fetchedAlice, ok := fetchBoth[aliceId]
	assert.Equal(t, ok, true)
	fetchedBob, ok := fetchBoth[bob]
	//fmt.Printf("k, bobcounts a, %+v\n", fetchedBob.Counts)
	assert.Equal(t, ok, true)
	assert.Assert(t, fetchedAlice.Counts != nil)
	assert.Assert(t, fetchedBob.Counts == nil)

	assert.Equal(t, fetchedAlice.Counts.Active, 1)
	assert.Equal(t, fetchedAlice.Counts.Complete, 0)
	assert.Equal(t, fetchedAlice.Counts.Notes, 0)
	var nilState *engine.GotState = nil
	assert.Equal(t, fetchedAlice.State, nilState)

	//complete bob, expect alice to become active
	err = agg.ItemStateChanged(StateChangeEvent{
		Id:       bob,
		OldState: engine.Active,
		NewState: &completeState,
		Ancestry: []engine.SummaryId{aliceId},
	})
	assert.NilError(t, err)

	//refetch
	fetchBoth, err = agg.summaryStore.Fetch([]engine.SummaryId{
		aliceId,
		bob,
	})

	assert.NilError(t, err)
	fetchedAlice, ok = fetchBoth[aliceId]
	assert.Equal(t, ok, true)
	fetchedBob, ok = fetchBoth[bob]
	assert.Equal(t, ok, true)
	assert.Assert(t, fetchedAlice.Counts != nil)
	assert.Assert(t, fetchedBob.Counts == nil)

	assert.Equal(t, fetchedAlice.Counts.Active, 0)
	assert.Equal(t, fetchedAlice.Counts.Complete, 1)
	assert.Equal(t, fetchedAlice.Counts.Notes, 0)
	assert.Equal(t, *fetchedAlice.State, activeState)

	//now add an active item under alice
	fmt.Println("adding carol..")
	err = agg.ItemAdded(AddItemEvent{
		Id:       carolId,
		State:    engine.Active,
		Ancestry: []engine.SummaryId{aliceId},
		Deadline: nil,
	})

	assert.NilError(t, err)

	fetchBoth, err = agg.summaryStore.Fetch([]engine.SummaryId{
		aliceId,
		carolId,
	})
	assert.NilError(t, err)

	fetchedAlice, ok = fetchBoth[aliceId]
	assert.Equal(t, ok, true)
	assert.Assert(t, fetchedAlice.State == nil)

	assert.Equal(t, fetchedAlice.Counts.Active, 1)   //carol
	assert.Equal(t, fetchedAlice.Counts.Complete, 1) //bob
	assert.Equal(t, fetchedAlice.Counts.Notes, 0)

	//complete carol
	err = agg.ItemStateChanged(StateChangeEvent{
		Id:       carolId,
		OldState: engine.Active,
		NewState: &completeState,
		Ancestry: []engine.SummaryId{aliceId},
	})
	assert.NilError(t, err)

	fetchBoth, err = agg.summaryStore.Fetch([]engine.SummaryId{
		aliceId,
		carolId,
	})
	assert.NilError(t, err)

	fetchedAlice, ok = fetchBoth[aliceId]
	assert.Equal(t, ok, true)

	assert.Equal(t, fetchedAlice.Counts.Active, 0)
	assert.Equal(t, fetchedAlice.Counts.Complete, 2)
	assert.Equal(t, fetchedAlice.Counts.Notes, 0)
	assert.Equal(t, *fetchedAlice.State, activeState) //alice has complete bob and complete carol under it. So its now the active node.
}

func TestAggregatorHandlesDelete(t *testing.T) {

	store := MakeMockSummaryStore()
	agg, err := NewAggregator(&store)
	assert.NilError(t, err)

	var aliceId = engine.SummaryId(10)
	var bob = engine.SummaryId(11)
	var carolId = engine.SummaryId(12)

	var activeState = engine.GotState(engine.Active)
	var completeState = engine.GotState(engine.Complete)

	//create alice -> bob
	err = agg.ItemAdded(AddItemEvent{
		Id:       aliceId,
		State:    engine.Active,
		Ancestry: []engine.SummaryId{},
		Deadline: nil,
	})
	assert.NilError(t, err)

	assert.Assert(t, store.aggs != nil)

	assert.Equal(t, len(store.aggs), 1)

	agg.ItemAdded(AddItemEvent{
		Id:       bob,
		State:    engine.Active,
		Ancestry: []engine.SummaryId{aliceId},
		Deadline: nil,
	})

	//fetch both
	assert.Equal(t, len(store.aggs), 2)
	fetchBoth, err := agg.summaryStore.Fetch([]engine.SummaryId{
		aliceId,
		bob,
	})
	assert.NilError(t, err)
	fetchedAlice, ok := fetchBoth[aliceId]
	assert.Equal(t, ok, true)
	fetchedBob, ok := fetchBoth[bob]
	//fmt.Printf("k, bobcounts a, %+v\n", fetchedBob.Counts)
	assert.Equal(t, ok, true)
	assert.Assert(t, fetchedAlice.Counts != nil)
	assert.Assert(t, fetchedBob.Counts == nil)

	assert.Equal(t, fetchedAlice.Counts.Active, 1)
	assert.Equal(t, fetchedAlice.Counts.Complete, 0)
	assert.Equal(t, fetchedAlice.Counts.Notes, 0)
	var nilState *engine.GotState = nil
	assert.Equal(t, fetchedAlice.State, nilState)

	//complete bob, expect alice to become active
	err = agg.ItemStateChanged(StateChangeEvent{
		Id:       bob,
		OldState: engine.Active,
		NewState: &completeState,
		Ancestry: []engine.SummaryId{aliceId},
	})
	assert.NilError(t, err)

	//refetch
	fetchBoth, err = agg.summaryStore.Fetch([]engine.SummaryId{
		aliceId,
		bob,
	})

	assert.NilError(t, err)
	fetchedAlice, ok = fetchBoth[aliceId]
	assert.Equal(t, ok, true)
	fetchedBob, ok = fetchBoth[bob]
	assert.Equal(t, ok, true)
	assert.Assert(t, fetchedAlice.Counts != nil)
	assert.Assert(t, fetchedBob.Counts == nil)

	assert.Equal(t, fetchedAlice.Counts.Active, 0)
	assert.Equal(t, fetchedAlice.Counts.Complete, 1)
	assert.Equal(t, fetchedAlice.Counts.Notes, 0)

	assert.Equal(t, *fetchedAlice.State, activeState)

	//now add an active item under alice
	err = agg.ItemAdded(AddItemEvent{
		Id:       carolId,
		State:    engine.Active,
		Ancestry: []engine.SummaryId{aliceId},
		Deadline: nil,
	})

	assert.NilError(t, err)

	fetchBoth, err = agg.summaryStore.Fetch([]engine.SummaryId{
		aliceId,
		carolId,
	})
	assert.NilError(t, err)

	fetchedAlice, ok = fetchBoth[aliceId]
	assert.Equal(t, ok, true)
	assert.Assert(t, fetchedAlice.State == nil)

	assert.Equal(t, fetchedAlice.Counts.Active, 1)   //carol
	assert.Equal(t, fetchedAlice.Counts.Complete, 1) //bob
	assert.Equal(t, fetchedAlice.Counts.Notes, 0)

	//deleting carol carol
	err = agg.ItemDeleted(ItemDeletedEvent{
		Id:       carolId,
		Ancestry: []engine.SummaryId{aliceId},
	})
	assert.NilError(t, err)

	fetchBoth, err = agg.summaryStore.Fetch([]engine.SummaryId{
		aliceId,
		carolId,
	})
	assert.NilError(t, err)

	fetchedAlice, ok = fetchBoth[aliceId]
	assert.Equal(t, ok, true)

	assert.Equal(t, fetchedAlice.Counts.Active, 0)
	assert.Equal(t, fetchedAlice.Counts.Complete, 1)
	assert.Equal(t, fetchedAlice.Counts.Notes, 0)
	assert.Assert(t, fetchedAlice.State != nil)
	assert.Equal(t, *fetchedAlice.State, activeState) //alice has complete bob and complete carol under it. So its now the active node.
}

func TestAggregatorMoveActiveItem(t *testing.T) {
	store := MakeMockSummaryStore()
	agg, err := NewAggregator(&store)
	assert.NilError(t, err)

	var aliceId = engine.SummaryId(10)
	var bobId = engine.SummaryId(11)
	var carolId = engine.SummaryId(12)

	// Create alice and bob as top-level groups with carol under alice
	// alice -> carol (active)
	// bob (empty)
	err = agg.ItemAdded(AddItemEvent{
		Id:       aliceId,
		State:    engine.Active,
		Ancestry: []engine.SummaryId{},
		Deadline: nil,
	})
	assert.NilError(t, err)

	err = agg.ItemAdded(AddItemEvent{
		Id:       bobId,
		State:    engine.Active,
		Ancestry: []engine.SummaryId{},
		Deadline: nil,
	})
	assert.NilError(t, err)

	err = agg.ItemAdded(AddItemEvent{
		Id:       carolId,
		State:    engine.Active,
		Ancestry: []engine.SummaryId{aliceId},
		Deadline: nil,
	})
	assert.NilError(t, err)

	// Verify initial state: alice has 1 active, bob has 0 (nil Counts = leaf)
	fetchAll, err := agg.summaryStore.Fetch([]engine.SummaryId{aliceId, bobId, carolId})
	assert.NilError(t, err)

	fetchedAlice := fetchAll[aliceId]
	fetchedBob := fetchAll[bobId]
	assert.Equal(t, fetchedAlice.Counts.Active, 1)
	assert.Assert(t, fetchedBob.Counts == nil) // bob is a leaf, no children

	// Move carol from alice to bob
	err = agg.ItemMoved(ItemMovedEvent{
		Id:          carolId,
		OldAncestry: []engine.SummaryId{aliceId},
		NewAncestry: []engine.SummaryId{bobId},
	})
	assert.NilError(t, err)

	// Refetch and verify: alice should have 0 active, bob should have 1 active
	fetchAll, err = agg.summaryStore.Fetch([]engine.SummaryId{aliceId, bobId, carolId})
	assert.NilError(t, err)

	fetchedAlice = fetchAll[aliceId]
	fetchedBob = fetchAll[bobId]

	assert.Equal(t, fetchedAlice.Counts.Active, 0)
	assert.Assert(t, fetchedBob.Counts != nil)
	assert.Equal(t, fetchedBob.Counts.Active, 1)
}

func TestAggregatorMoveCompleteItem(t *testing.T) {
	store := MakeMockSummaryStore()
	agg, err := NewAggregator(&store)
	assert.NilError(t, err)

	var aliceId = engine.SummaryId(10)
	var bobId = engine.SummaryId(11)
	var carolId = engine.SummaryId(12)
	var completeState = engine.GotState(engine.Complete)

	// Create alice -> carol (active), bob (empty)
	err = agg.ItemAdded(AddItemEvent{
		Id:       aliceId,
		State:    engine.Active,
		Ancestry: []engine.SummaryId{},
		Deadline: nil,
	})
	assert.NilError(t, err)

	err = agg.ItemAdded(AddItemEvent{
		Id:       bobId,
		State:    engine.Active,
		Ancestry: []engine.SummaryId{},
		Deadline: nil,
	})
	assert.NilError(t, err)

	err = agg.ItemAdded(AddItemEvent{
		Id:       carolId,
		State:    engine.Active,
		Ancestry: []engine.SummaryId{aliceId},
		Deadline: nil,
	})
	assert.NilError(t, err)

	// Complete carol
	err = agg.ItemStateChanged(StateChangeEvent{
		Id:       carolId,
		OldState: engine.Active,
		NewState: &completeState,
		Ancestry: []engine.SummaryId{aliceId},
	})
	assert.NilError(t, err)

	// Verify: alice has 1 complete, bob has 0 (nil Counts = leaf)
	fetchAll, err := agg.summaryStore.Fetch([]engine.SummaryId{aliceId, bobId, carolId})
	assert.NilError(t, err)

	fetchedAlice := fetchAll[aliceId]
	fetchedBob := fetchAll[bobId]
	assert.Equal(t, fetchedAlice.Counts.Complete, 1)
	assert.Assert(t, fetchedBob.Counts == nil) // bob is a leaf, no children

	// Move carol from alice to bob
	err = agg.ItemMoved(ItemMovedEvent{
		Id:          carolId,
		OldAncestry: []engine.SummaryId{aliceId},
		NewAncestry: []engine.SummaryId{bobId},
	})
	assert.NilError(t, err)

	// Verify: alice has 0 complete, bob has 1 complete
	fetchAll, err = agg.summaryStore.Fetch([]engine.SummaryId{aliceId, bobId, carolId})
	assert.NilError(t, err)

	fetchedAlice = fetchAll[aliceId]
	fetchedBob = fetchAll[bobId]

	assert.Equal(t, fetchedAlice.Counts.Complete, 0)
	assert.Assert(t, fetchedBob.Counts != nil)
	assert.Equal(t, fetchedBob.Counts.Complete, 1)
}

func TestAggregatorMoveItemDeepHierarchy(t *testing.T) {
	store := MakeMockSummaryStore()
	agg, err := NewAggregator(&store)
	assert.NilError(t, err)

	var topId = engine.SummaryId(9)
	var aliceId = engine.SummaryId(10)
	var bobId = engine.SummaryId(11)
	var carolId = engine.SummaryId(12)

	// Create hierarchy: top -> alice -> carol (active), top -> bob (empty)
	err = agg.ItemAdded(AddItemEvent{
		Id:       topId,
		State:    engine.Active,
		Ancestry: []engine.SummaryId{},
		Deadline: nil,
	})
	assert.NilError(t, err)

	err = agg.ItemAdded(AddItemEvent{
		Id:       aliceId,
		State:    engine.Active,
		Ancestry: []engine.SummaryId{topId},
		Deadline: nil,
	})
	assert.NilError(t, err)

	err = agg.ItemAdded(AddItemEvent{
		Id:       bobId,
		State:    engine.Active,
		Ancestry: []engine.SummaryId{topId},
		Deadline: nil,
	})
	assert.NilError(t, err)

	err = agg.ItemAdded(AddItemEvent{
		Id:       carolId,
		State:    engine.Active,
		Ancestry: []engine.SummaryId{topId, aliceId},
		Deadline: nil,
	})
	assert.NilError(t, err)

	// Verify initial counts:
	// - top has 2 active (bob as leaf + carol under alice)
	// - alice has 1 active (carol)
	// - bob is a leaf (nil Counts)
	fetchAll, err := agg.summaryStore.Fetch([]engine.SummaryId{topId, aliceId, bobId})
	assert.NilError(t, err)

	fetchedTop := fetchAll[topId]
	fetchedAlice := fetchAll[aliceId]
	fetchedBob := fetchAll[bobId]

	assert.Equal(t, fetchedTop.Counts.Active, 2)  // bob + carol
	assert.Equal(t, fetchedAlice.Counts.Active, 1) // carol
	assert.Assert(t, fetchedBob.Counts == nil)     // bob is a leaf

	// Move carol from alice to bob (both under top)
	// Old ancestry: [top, alice]
	// New ancestry: [top, bob]
	err = agg.ItemMoved(ItemMovedEvent{
		Id:          carolId,
		OldAncestry: []engine.SummaryId{topId, aliceId},
		NewAncestry: []engine.SummaryId{topId, bobId},
	})
	assert.NilError(t, err)

	// Verify: top still has 2 active, alice has 0 active, bob has 1 active
	fetchAll, err = agg.summaryStore.Fetch([]engine.SummaryId{topId, aliceId, bobId})
	assert.NilError(t, err)

	fetchedTop = fetchAll[topId]
	fetchedAlice = fetchAll[aliceId]
	fetchedBob = fetchAll[bobId]

	// top: decrement 1 (old) + increment 1 (new) = net 0 change, still 2
	assert.Equal(t, fetchedTop.Counts.Active, 2)
	// alice: decremented from 1 to 0
	assert.Equal(t, fetchedAlice.Counts.Active, 0)
	// bob: now has carol, so Counts.Active = 1
	assert.Assert(t, fetchedBob.Counts != nil)
	assert.Equal(t, fetchedBob.Counts.Active, 1)
}

func TestAggregatorMoveNoteItem(t *testing.T) {
	store := MakeMockSummaryStore()
	agg, err := NewAggregator(&store)
	assert.NilError(t, err)

	var aliceId = engine.SummaryId(10)
	var bobId = engine.SummaryId(11)
	var noteId = engine.SummaryId(12)

	// Create alice -> note, bob (empty)
	err = agg.ItemAdded(AddItemEvent{
		Id:       aliceId,
		State:    engine.Active,
		Ancestry: []engine.SummaryId{},
		Deadline: nil,
	})
	assert.NilError(t, err)

	err = agg.ItemAdded(AddItemEvent{
		Id:       bobId,
		State:    engine.Active,
		Ancestry: []engine.SummaryId{},
		Deadline: nil,
	})
	assert.NilError(t, err)

	err = agg.ItemAdded(AddItemEvent{
		Id:       noteId,
		State:    engine.Note,
		Ancestry: []engine.SummaryId{aliceId},
		Deadline: nil,
	})
	assert.NilError(t, err)

	// Verify initial: alice has 1 note, bob has 0 (nil Counts = leaf)
	fetchAll, err := agg.summaryStore.Fetch([]engine.SummaryId{aliceId, bobId})
	assert.NilError(t, err)

	fetchedAlice := fetchAll[aliceId]
	fetchedBob := fetchAll[bobId]
	assert.Equal(t, fetchedAlice.Counts.Notes, 1)
	assert.Assert(t, fetchedBob.Counts == nil) // bob is a leaf, no children

	// Move note from alice to bob
	err = agg.ItemMoved(ItemMovedEvent{
		Id:          noteId,
		OldAncestry: []engine.SummaryId{aliceId},
		NewAncestry: []engine.SummaryId{bobId},
	})
	assert.NilError(t, err)

	// Verify: alice has 0 notes, bob has 1 note
	fetchAll, err = agg.summaryStore.Fetch([]engine.SummaryId{aliceId, bobId})
	assert.NilError(t, err)

	fetchedAlice = fetchAll[aliceId]
	fetchedBob = fetchAll[bobId]

	assert.Equal(t, fetchedAlice.Counts.Notes, 0)
	assert.Assert(t, fetchedBob.Counts != nil)
	assert.Equal(t, fetchedBob.Counts.Notes, 1)
}
