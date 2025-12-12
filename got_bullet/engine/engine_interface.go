package engine

import (
	"errors"
	"math"

	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/ids"
)

type SummaryId int32

// / This is the machine that takes the commands, changes the backend state and returns wahts requested.
type GotEngine interface {
	Summary(lookup *GidLookup) (*GotItemDisplay, error)

	EditTitle(lookup GidLookup, newHeading string) error
	//state changes
	MarkResolved(lookup []GidLookup) error
	MarkActive(lookup GidLookup) (*NodeId, error)
	MarkAsNote(lookup GidLookup) (*NodeId, error)

	Delete(lookup GidLookup) (*NodeId, error)

	Move(lookup GidLookup, newParent GidLookup) (*NodeId, error) //returns the oldParents id

	GotAliasInterface
	GotCreateItemInterface
	GotFetchInterface
}

type GotCreateItemInterface interface {
	CreateBuck(parent *GidLookup, date *DateLookup, completable bool, heading string) (*NodeId, error)
}

// descendant types
const (
	AllDescendants           = 0
	LeafNodesOnly            = 1
	ImmediateDescendantsOnly = 2
)

type GotState int

const (
	CompleteChar = "✔"
	ActiveChar   = "•"
	//bulletChar = "!"
	NoteChar = "~"
)

func (g GotState) ToStr() string {
	if g == Active {
		return ActiveChar
	}
	if g == Note {
		return NoteChar
	}
	if g == Complete {
		return CompleteChar
	}
	return "<?>"
}

// states
const (
	Active   = 0
	Note     = 8000
	Complete = 16000
)

type GotFetchResult struct {
	Result []GotItemDisplay
}

// All the lookup stuff
type GotFetchInterface interface {
	FetchItemsBelow(lookup *GidLookup, descendantType int, states []int) (*GotFetchResult, error)
}

// The interface for all aliasing functionality
type GotAliasInterface interface {
	Lookup(alias string) (*GotId, error)
	LookupAliasForGid(gid string) (*string, error)
	LookupAliasForMany(gid []string) (map[string]*string, error)
	Unalias(alias string) (*GotId, error)
	Alias(gid string, alias string) (bool, error)
}

type DateLookup struct {
	UserInput string
}

// User entered best guess at a gid. Might be the Gid, might be an alias. Might be the title
type GidLookup struct {
	Input string
}

type GotItemDisplay struct {
	Gid        string
	Title      string
	Alias      string
	Deadline   string
	SummaryObj *Summary
	Path       *GotPath
	NumberGo   int
}

func (i *GotItemDisplay) IsNote() bool {
	if i.SummaryObj != nil && i.SummaryObj.State != nil && *i.SummaryObj.State == Note {
		return true
	}
	return false
}
func (i *GotItemDisplay) IsActive() bool {
	if i.SummaryObj != nil && i.SummaryObj.State != nil && *i.SummaryObj.State == Active {
		return true
	}
	return false
}
func (i *GotItemDisplay) IsComplete() bool {
	if i.SummaryObj != nil && i.SummaryObj.State != nil && *i.SummaryObj.State == Complete {
		return true
	}
	return false
}

// VX:TODO not tested, used for sorting the items.
func (i *GotItemDisplay) FullPathString() string {
	var path = ""
	for _, p := range i.Path.Ancestry {
		s, _ := p.Shortcut()
		path += "/" + s
	}
	s, _ := i.Shortcut()
	return path + "/" + s
}

// either alias or gid, and true for alias, false for gid
func (i *PathItem) Shortcut() (string, bool) {
	if i.Alias != nil {
		return *i.Alias, true
	} else {
		return i.Id, false
	}

}

// either alias or gid, and true for alias, false for gid
func (i *GotItemDisplay) Shortcut() (string, bool) {
	if i.Alias != "" {
		return i.Alias, true
	} else {
		return i.Gid, false
	}

}

// VX:TODO replace with GotId
type Gid struct {
	Id string
}

// VX:TODO RM?
type NodeId struct {
	Gid   Gid
	Title string
	Alias string
}
type PathItem struct {
	Id    string
	Alias *string
}

type GotPath struct {
	Ancestry []PathItem
}

// VX:TODO consider using BUlletId semantics to lazy compute these
type GotId struct {
	AasciValue string
	IntValue   int32
}

func NewCompleteId(aasci string, intValue int32) GotId {
	return GotId{
		AasciValue: aasci,
		IntValue:   intValue,
	}
}

func FitsInInt32(v int64) bool {
	return v >= math.MinInt32 && v <= math.MaxInt32
}

// this is basically a wrapper for BulletId
func NewGotId(aasci string) (*GotId, error) {
	intVal, err := bullet_stl.AasciBulletIdToInt(aasci)
	if !FitsInInt32(intVal) {
		return nil, errors.New("id is too big")
	}

	if err != nil {
		return nil, err
	}
	return &GotId{
		AasciValue: aasci,
		IntValue:   int32(intVal),
	}, nil
}

//Lets be super clear about the ides
