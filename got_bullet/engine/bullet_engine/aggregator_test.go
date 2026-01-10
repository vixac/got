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
		NewState: engine.Complete,
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
	var activeState = engine.GotState(engine.Active)
	assert.Equal(t, *fetchedAlice.State, activeState)

	//now complete alice:

	//complete alice
	err = agg.ItemStateChanged(StateChangeEvent{
		Id:       aliceId,
		OldState: engine.Active,
		NewState: engine.Complete,
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
	var completeState = engine.GotState(engine.Complete)
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
		NewState: engine.Complete,
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
	var activeState = engine.GotState(engine.Active)
	assert.Equal(t, *fetchedAlice.State, activeState)

	assert.Equal(t, fetchTop.Counts.Active, 1)
	assert.Equal(t, fetchTop.Counts.Complete, 1)
	assert.Equal(t, fetchTop.Counts.Notes, 0)

	//now complete alice:

	//complete alice
	err = agg.ItemStateChanged(StateChangeEvent{
		Id:       aliceId,
		OldState: engine.Active,
		NewState: engine.Complete,
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
	var completeState = engine.GotState(engine.Complete)
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
		NewState: engine.Complete,
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
	var activeState = engine.GotState(engine.Active)
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
		NewState: engine.Complete,
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
