package grove_engine

import (
	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	"vixac.com/got/engine"
)

type GroveDeltas struct {
	Active   int
	Complete int
}

/*
type AggregateKey string
type MutationID string
type AggregateValue int64
type AggregateDeltas map[AggregateKey]AggregateValue
*/

func NewMutationDelta(state engine.GotState) GroveDeltas {
	deltas := GroveDeltas{}
	if state == engine.Active {
		deltas.Active = 1
	} else if state == engine.Complete {
		deltas.Complete = 1
	}
	return deltas

}

func (g *GroveDeltas) ToGrove() bullet_interface.AggregateDeltas {
	res := make(map[bullet_interface.AggregateKey]bullet_interface.AggregateValue)
	res["active"] = bullet_interface.AggregateValue(g.Active)
	res["complete"] = bullet_interface.AggregateValue(g.Complete)
	return res
}
