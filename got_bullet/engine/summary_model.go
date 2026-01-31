package engine

import (
	"encoding/json"
	"fmt"
	"time"

	"vixac.com/got/console"
)

// First pass of the kinds of things we'll count
type AggCount struct {
	Complete int `json:"c,omitempty"`
	Active   int `json:"a,omitempty"`
	Notes    int `json:"n,omitempty"`
}

// VX:TODO its either state OR its counts.
// deadline is separate. Maybe it doesn't belong here but we'll see.
type Summary struct {
	State       *GotState `json:"s,omitempty"`
	Counts      *AggCount `json:"c"`
	Deadline    *DateTime `json:"d"`
	CreatedDate *DateTime `json:"cr,omitempty"`
	UpdatedDate *DateTime `json:"u,omitempty"`
	Tags        []Tag     `json:"t,omitempty"`
}

type TagLiteral struct {
	Display string `json:"d,omitempty"`
	Token   string `json:"t,omitempty"`
}

// a tag can be represented either as a literal tag or as an identifier to a tag description
type Tag struct {
	Identifier *string     `json:"identifier,omitempty"`
	Literal    *TagLiteral `json:"tagLiteral,omitempty"`
}

func (lhs Tag) EqualTo(rhs Tag) bool {
	if lhs.Identifier != nil && rhs.Identifier != nil {
		return lhs.Identifier == rhs.Identifier
	}
	return lhs.Literal.Display == rhs.Literal.Display
}

func NewSummary(state GotState, deadline *DateTime, created *DateTime, tags []Tag) Summary {
	return Summary{
		State:       &state,
		Counts:      nil,
		Deadline:    deadline,
		CreatedDate: created,
		UpdatedDate: created,
		Tags:        tags,
	}
}

type DateTime struct {
	Special string `json:"s,omitempty"` //denotes a non-date datetime, such as "now"
	Date    string `json:"d,omitempty"`
	Millis  int64  `json:"e,omitempty"`
}

func (d *DateTime) EpochMillis() int64 {
	if d == nil || d.Special != "" {
		return 0
	}
	if d.Millis != 0 {
		return d.Millis
	}
	date, err := d.ToDate()
	if err != nil || date == nil {
		return 0
	}
	startOfDay := time.Time(*date)
	return startOfDay.UnixMilli()
}
func (d *DateTime) ToDate() (*console.RFC3339Time, error) {
	if d == nil || d.Special != "" {
		return nil, nil
	}
	var date console.RFC3339Time
	dateBytes := []byte(d.Date)
	err := json.Unmarshal(dateBytes, &date)
	if err != nil {
		return nil, err
	}
	return &date, nil
}
func (d *DateTime) JsonDateToReadable() (string, error) {
	if d == nil || d.Special != "" {
		return "", nil
	}

	var date console.RFC3339Time
	dateBytes := []byte(d.Date)
	err := json.Unmarshal(dateBytes, &date)
	if err != nil {
		fmt.Printf("VXL ERROR parsing is %s", err)
		return "", err
	}

	dateStr := console.DayFormat(time.Time(date))
	return dateStr, nil
}

func NewDateTime(time time.Time) (DateTime, error) {

	dateJsonByes, err := time.MarshalJSON()
	if err != nil {
		return DateTime{}, err
	}
	return DateTime{Date: string(dateJsonByes), Millis: time.UnixMilli()}, nil
}

func NewDeadlineFromDateLookup(inputString string, now time.Time) (DateTime, error) {
	if inputString == "<now>" {
		return DateTime{
			Special: "n",
		}, nil
	}

	deadlineTime, err := console.ParseRelativeDate(inputString, now)
	if err != nil {
		return DateTime{}, err
	}
	return NewDateTime((deadlineTime))
}

type DatedTask struct {
	Deadline DateTime  `json:"d"`
	Id       SummaryId `json:"i,omitempty"`
}

// something which can be combined and chained to form a single agg
type AggregateCountChange struct {
	NoteInc     int
	ActiveInt   int
	CompleteInc int
}

func (a *Summary) ApplyChange(change AggregateCountChange) {
	var count = AggCount{}
	if a.Counts != nil {
		count = *a.Counts
	}
	count.Active += change.ActiveInt
	count.Complete += change.CompleteInc
	count.Notes += change.NoteInc
	a.Counts = &count

}

func NewCountChange(state GotState, inc bool) AggregateCountChange {

	var change = 1
	if !inc {
		change = -1
	}
	if state == Active {
		return AggregateCountChange{
			ActiveInt: change,
		}
	}
	if state == Complete {
		return AggregateCountChange{
			CompleteInc: change,
		}
	}
	return AggregateCountChange{
		NoteInc: change,
	}
}

func (lhs AggregateCountChange) Combine(rhs AggregateCountChange) AggregateCountChange {
	return AggregateCountChange{
		ActiveInt:   lhs.ActiveInt + rhs.ActiveInt,
		NoteInc:     lhs.NoteInc + rhs.NoteInc,
		CompleteInc: lhs.CompleteInc + rhs.CompleteInc,
	}
}

// no count, no deadline for some reason
func NewLeafSummary(state GotState, deadline *DateTime, now time.Time, tags []Tag) Summary {
	dateTime, _ := NewDateTime(now)
	return NewSummary(state, deadline, &dateTime, tags)
}

func (c AggCount) ChangeState(state GotState, inc int) AggCount {
	comp := c.Complete
	active := c.Active
	notes := c.Notes
	if state == Active {
		active += inc
	} else if state == Complete {
		comp += inc
	} else if state == Note {
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

func (a *Summary) UpdatedCount(newCount AggCount) Summary {
	return Summary{
		State:       a.State,
		Counts:      &newCount,
		Deadline:    a.Deadline,
		CreatedDate: a.CreatedDate,
		UpdatedDate: a.UpdatedDate,
		Tags:        a.Tags,
	}
}
