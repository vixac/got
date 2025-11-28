package bullet_engine

import (
	"fmt"
)

/*
*
The aggregator is going to contain the business logic that maps events to
changes in the aggstore.
*/

type Aggregator struct {
	summaryStore SummaryStoreInterface
	// ancestorStore AncestorListInterface
}

func (a *Aggregator) ItemAdded(e AddItemEvent) error {
	ancestorAggs, err := a.summaryStore.Fetch(e.Ancestry)
	if err != nil {
		return err
	}
	//VX:TODO add the item

	//step 2, increment ancestry state
	if len(e.Ancestry) == 0 {
		//presumably this means we are adding to the root node?
		fmt.Printf("VX: THIS IS A ROOT NODE CHILD")
		return nil
	}
	//check if we need to update ancestors
	lastId := e.Ancestry[len(e.Ancestry)-1]
	last := ancestorAggs[lastId]

	//the changes to propagate

	var inc AggregateCountChange
	var idsToUpdate = e.Ancestry
	if last.IsLeaf() {
		//that means decrement the total by that state. We're deleting this
		inc = inc.combine(NewCountChange(*last.State, false))
		idsToUpdate = idsToUpdate[:len(idsToUpdate)-1] //we don't update the final item with this decrement

	}
	//ok time to finish this. 2 cases are adding to a leaf and not adding to a leaf
	//	and i suppose adding to root. 3 cases, each slightly different. lets code thema
	//	separately and test them separately. We should test this shit properly as its business logica
	//check if parent is a leaf. If so, we need to move it to group state, and delete it

	fmt.Printf("VX:TODO unhandled event ")
	return nil
}

func (a *Aggregator) ItemStateChanged(e StateChangeEvent) error {
	fmt.Printf("VX:TODO unhandled event ")
	return nil
}
func (a *Aggregator) ItemDeleted(e ItemDeletedEvent) error {
	fmt.Printf("VX:TODO unhandled event ")
	return nil
}

func (a *Aggregator) ItemMoved(e ItemMovedEvent) error {
	fmt.Printf("VX:TODO unhandled event ")
	return nil
}
