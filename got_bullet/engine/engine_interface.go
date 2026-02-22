package engine

import (
	"errors"
	"math"
	"strings"
	"unicode"

	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/ids"
	"vixac.com/got/console"
)

type SummaryId int32

// / This is the machine that takes the commands, changes the backend state and returns wahts requested.
type GotEngine interface {
	EditTitle(lookup GidLookup, newHeading string) error
	//state changes
	MarkResolved(lookup []GidLookup) error
	MarkActive(lookup GidLookup) (*NodeId, error)
	MarkAsNote(lookup GidLookup) (*NodeId, error)
	DeleteMany(lookups []GidLookup) error
	ToggleCollapse(lookup GidLookup, collapsed bool) error

	Move(lookup GidLookup, newParent GidLookup) (*NodeId, error) //returns the oldParents id
	OpenThenTimestamp(lookup GidLookup) error
	ScheduleItem(lookup GidLookup, dateLookup DateLookup) error
	TagItem(lookup GidLookup, tag TagLookup) error

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
	CompleteChar  = "✔"
	ActiveChar    = "⏺"
	NoteChar      = "~"
	TNoteChar     = "📎"
	CollapsedChar = "📁"
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
	Parent *GotItemDisplay
	Result []GotItemDisplay
}

// All the lookup stuff
type GotFetchInterface interface {
	FetchItemsBelow(lookup *GidLookup, sortByPath bool, states []GotState) (*GotFetchResult, error)
}

// The interface for all aliasing functionality
type GotAliasInterface interface {
	Lookup(alias string) (*GotId, error)
	LookupAliasForGid(gid string) (*string, error)
	LookupAliasForMany(gid []string) (map[string]*string, error)
	Unalias(alias string) (*GotId, error)
	Alias(lookup GidLookup, alias string) (bool, error)
}

func IsValidAlias(input string) bool {
	if len(input) == 0 {
		return false
	}
	spaces := strings.Contains(input, " ")
	if spaces {
		return false
	}
	bytes := []byte(input)
	firstCharIsNumber := CheckNumber([]byte{bytes[0]})
	if firstCharIsNumber {
		return false
	}
	return true

}

func CheckNumber(p []byte) bool {
	r := string(p)
	sep := 0
	for _, b := range r {
		if unicode.IsNumber(b) {
			continue
		}
		if b == rune('.') {
			if sep > 0 {
				return false
			}
			sep++
			continue
		}
		return false
	}
	return true
}

type DateLookup struct {
	UserInput string
}

func NowDateLookup() DateLookup {
	return DateLookup{
		UserInput: "<now>",
	}
}

// For now we'll treat this as a tag literal.
type TagLookup struct {
	Input string
}

// User entered best guess at a gid. Might be the Gid, might be an alias. Might be the title
type GidLookup struct {
	Input string
}

type GotItemDisplay struct {
	GotId         GotId
	DisplayGid    string
	Title         string
	Alias         string
	Deadline      string
	Created       string
	Updated       string
	DeadlineToken console.Token
	SummaryObj    *Summary
	Path          *GotPath
	NumberGo      int
	HasTNote      bool
}

func (i *GotItemDisplay) IsCollapsed() bool {
	return i.SummaryObj.Flags != nil && i.SummaryObj.Flags["collapsed"] == true
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

type PathString struct {
	DisplayPath string //either ID or alias for each node
	IdPath      string //path but based on ids only
}

type Shortcut struct {
	Display string //alias or id
	Id      string //just the id
}

// VX:TODO not tested, used for sorting the items.
func (i *GotItemDisplay) FullPathString() PathString {
	var displayPath = ""
	var idPath = ""
	delimiter := "/"
	for _, p := range i.Path.Ancestry {
		s, _ := p.Shortcut()
		displayPath += delimiter + s.Display
		idPath += delimiter + s.Id

	}
	s, _ := i.Shortcut()
	displayPath += delimiter + s.Display
	idPath += delimiter + s.Id
	return PathString{DisplayPath: displayPath, IdPath: idPath}
}

// either alias or gid, and true for alias, false for gid
func (i *PathItem) Shortcut() (Shortcut, bool) {
	var display = ""
	var aliased = false
	if i.Alias != nil {
		display = *i.Alias
		aliased = true
	} else {
		display = i.Id
	}
	return Shortcut{Display: display, Id: i.Id}, aliased
}

// either alias or gid, and true for alias, false for gid
func (i *GotItemDisplay) Shortcut() (Shortcut, bool) {
	var display = ""
	var aliased = false
	if i.Alias != "" {
		display = i.Alias
		aliased = true
	} else {
		display = i.DisplayGid //this one has a "0" prefix
	}
	return Shortcut{Display: display, Id: i.GotId.AasciValue}, aliased //i.GotId.Aasci value has no 0 prefix, useful for the actual path.

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

func (p *GotPath) Depth() int {
	if p == nil {
		return 0
	}
	return len(p.Ancestry)
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
