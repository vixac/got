package grove_engine

import (
	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	"vixac.com/got/engine"
)

const (
	activeKey   = "active"
	completeKey = "complete"
)

type GroveAggregate struct {
	Active   int
	Complete int
}

func NewAggregate(groveMap map[bullet_interface.AggregateKey]bullet_interface.AggregateValue) GroveAggregate {
	activeCount := 0
	completeCount := 0
	active, ok := groveMap[activeKey]
	if ok {
		activeCount = int(active)
	}
	complete, ok := groveMap[completeKey]
	if ok {
		completeCount = int(complete)
	}
	return GroveAggregate{
		Active:   activeCount,
		Complete: completeCount,
	}

}
func NewMutationDelta(state engine.GotState) GroveAggregate {
	deltas := GroveAggregate{}
	if state == engine.Active {
		deltas.Active = 1
	} else if state == engine.Complete {
		deltas.Complete = 1
	}
	return deltas

}

func (g *GroveAggregate) ToGrove() bullet_interface.AggregateDeltas {
	res := make(map[bullet_interface.AggregateKey]bullet_interface.AggregateValue)
	res[activeKey] = bullet_interface.AggregateValue(g.Active)
	res[completeKey] = bullet_interface.AggregateValue(g.Complete)
	return res
}
