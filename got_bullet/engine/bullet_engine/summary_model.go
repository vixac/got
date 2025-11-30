package bullet_engine

import (
	"fmt"

	"vixac.com/got/engine"
)

// First pass of the kinds of things we'll count
type AggCount struct {
	Complete int `json:"c,omitempty"`
	Active   int `json:"a,omitempty"`
	Notes    int `json:"n,omitempty"`
}

type Deadline struct {
	Date string `json:"d,omitempty"`
}

type DatedTask struct {
	Deadline Deadline  `json:"d"`
	Id       SummaryId `json:"i,omitempty"`
}

// VX:TODO its either state OR its counts.
// deadline is separate. Maybe it doesn't belong here but we'll see.
type Summary struct {
	State    *engine.GotState `json:"s,omitempty"`
	Counts   *AggCount        `json:"c"`
	Deadline *Deadline        `json:"d"`
}

// something which can be combined and chained to form a single agg
type AggregateCountChange struct {
	NoteInc     int
	ActiveInt   int
	CompleteInc int
}

func (a Summary) ApplyChange(change AggregateCountChange) Summary {
	var count = AggCount{}
	if a.Counts != nil {
		count = *a.Counts
	}
	count.Active += change.ActiveInt
	count.Complete += change.CompleteInc
	count.Notes += change.NoteInc
	var old = ""
	if a.Counts != nil {
		old = fmt.Sprintf("%+v", *a.Counts)

	}

	fmt.Printf("VX: COUNT UP CHANGE CHECK IS from %+s ... to %+v\n", old, count)
	a.Counts = &count

	return a
}

func NewCountChange(state engine.GotState, inc bool) AggregateCountChange {

	var change = 1
	if !inc {
		change = -1
	}
	if state == engine.Active {
		return AggregateCountChange{
			ActiveInt: change,
		}
	}
	if state == engine.Complete {
		return AggregateCountChange{
			CompleteInc: change,
		}
	}
	return AggregateCountChange{
		NoteInc: change,
	}
}

/*
func (lhs AggregateCountChange) combine(rhs AggregateCountChange) AggregateCountChange {
	return AggregateCountChange{
		ActiveInt:   lhs.ActiveInt + rhs.ActiveInt,
		NoteInc:     lhs.NoteInc + rhs.NoteInc,
		CompleteInc: lhs.CompleteInc + rhs.CompleteInc,
	}
}
*/
// no count, no deadline for some reason
func NewLeafSummary(state engine.GotState, deadline *Deadline) Summary {
	return Summary{
		State:    &state,
		Deadline: deadline,
	}
}

func (c AggCount) ChangeState(state engine.GotState, inc int) AggCount {
	comp := c.Complete
	active := c.Active
	notes := c.Notes
	if state == engine.Active {
		active += inc
	} else if state == engine.Complete {
		comp += inc
	} else if state == engine.Note {
		notes += inc
	}
	return AggCount{
		Complete: comp,
		Active:   active,
		Notes:    notes,
	}
}

func (a Summary) IsLeaf() bool {
	return a.State != nil
}

func (c AggCount) changeActive(inc int) AggCount {
	return AggCount{
		c.Complete,
		c.Active + inc,
		c.Notes,
	}
}
func (c AggCount) changeNotes(inc int) AggCount {
	return AggCount{
		c.Complete,
		c.Active,
		c.Notes + inc,
	}
}
func (c AggCount) changeComplete(inc int) AggCount {
	return AggCount{
		c.Complete + inc,
		c.Active,
		c.Notes,
	}
}
func (a *Summary) UpdatedCount(newCount AggCount) Summary {
	return Summary{
		State:    a.State,
		Counts:   &newCount,
		Deadline: a.Deadline,
	}
}
