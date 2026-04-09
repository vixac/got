package engine

import (
	"errors"
	"math"
	"time"

	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/ids"
)

type SummaryId int32

// / This is the machine that takes the commands, changes the backend state and returns wahts requested.
type GotEngineInterface interface {
	GotAliasInterface
	GotCreateItemInterface
	GotFetchInterface
	RestoreInterface
	NoteInterface
	GotEditInterface
	GotTreeInterface
}
type GotTreeInterface interface {
	DeleteMany(lookups []GidLookup) error
	Move(lookup GidLookup, newParent GidLookup) error //returns the oldParents id
}

type GotEditInterface interface {
	MarkResolved(lookup []GidLookup) error
	EditTitle(lookup GidLookup, newHeading string) error
	ScheduleItem(lookup GidLookup, dateLookup DateLookup) error
	TagItem(lookup GidLookup, tag TagLookup) error
	ToggleCollapse(lookup GidLookup, collapsed bool) error
}

type NoteInterface interface {
	JotNote(lookup GidLookup, note string) (LongFormKey, error)
	NotesFor(lookup *GidLookup, recurse bool) (*LongFormBlockResult, error)
	OpenThenTimestamp(lookup GidLookup) error
}
type RestoreInterface interface {
	CreateStoreFile() (string, error)
	RestoreFromFile(filename string) error
}

// All the lookup stuff
type GotFetchInterface interface {
	FetchItemsBelow(lookup *GidLookup, sortByPath bool, states []GotState, hideUnderCollapsed bool) (*GotFetchResult, error)
}

type GotCreateItemInterface interface {
	CreateBuck(request CreateBuckRequest) (*GotId, error)
}

// The interface for all aliasing functionality
type GotAliasInterface interface {
	LookupAliasForMany(gid []string) (map[string]*string, error)
	Unalias(alias string) (*GotId, error)
	Alias(lookup GidLookup, alias string) (bool, error)
}

// ///VX:TODO everything under here doesn't belong in this file, and probably belongs in util.
type IdGeneratorInterface interface {
	SetLastIdIfLower(newId int64) error //if we're overriding the ids, the last Id may be replaced with this one.
	LastId() (int64, error)             //fetches the last createdId
	NextId() (int64, error)             //creates a new id, stores it as the lastId, and returns it
}

// The store that holds on to the meanings of the number goes, so when user
// can use them async
type NumberGoStoreInterface interface {
	AssignNumberPairs(pairs []NumberGoPair) error
	GidFor(number int) (*GotId, error)
}

type NumberGoPair struct {
	Number int    `json:"n"`
	Gid    string `json:"g"`
}

type SummaryStoreInterface interface {
	UpsertSummary(id SummaryId, agg Summary) error
	UpsertManySummaries(aggs map[SummaryId]Summary) error
	Fetch(ids []SummaryId) (map[SummaryId]Summary, error)
	Delete(ids []SummaryId) error
}

type GidLookupInterface interface {
	InputToGid(lookup *GidLookup) (*GotId, error)
}

type LongFormBlockResult struct {
	Blocks []LongFormBlock
}

type LongFormBlock struct {
	Id      LongFormKey
	Content string
	Edited  time.Time
}

func (l *LongFormBlock) Created() time.Time {
	return l.Id.CreatedTime
}

type LongFormStoreInterface interface {
	AppendNote(id GotId, content string) (*LongFormKey, error) //creates and appends a block by using an incremented id, and the current timetstamp for edit time.
	InsertBlock(block LongFormBlock) error                     //similar to appendNote, except the block includes meta data
	LongFormNotesFor(id GotId) (*LongFormBlockResult, error)
	LongFormForMany(ids []GotId) (map[GotId]LongFormBlockResult, error)
	RemoveAllItemsFromLongStoreUnder(id GotId) error
}

// Contains the values for fields that would normally be populated by the engine
type CreateOverrideSettings struct {
	OverrideId   *int32                 `json:"g,omitempty"`
	UpdatedDate  string                 `json:"u,omitempty"`
	CreatedDate  string                 `json:"c,omitempty"`
	ScheduleDate *DateTime              `json:"d,omitempty"`
	Alias        *string                `json:"a,omitempty"`
	NoAlias      bool                   `json:"no,omitempty"` //no override isnt the same as explicitly no alias at all
	Tags         []Tag                  `json:"t,omitempty"`
	Flags        []string               `json:"f,omitempty"`
	LongForm     []LongFormRestoreBlock `json:"l,omitempty"`
}

type LongFormRestoreBlock struct {
	KeyString  string `json:"k,omitempty"`
	Content    string `json:"c,omitempty"`
	EditMillis string `json:"m,omitempty"`
}

func NewRestoreBlock(block LongFormBlock) LongFormRestoreBlock {
	restoreBlock := LongFormRestoreBlock{
		KeyString:  block.Id.ToString(),
		Content:    block.Content,
		EditMillis: TimeToMillisString(block.Edited),
	}
	return restoreBlock
}

// VX:if createBuckRequests have idempotency keys, then retrying a failed create buck might be permissible.
type CreateBuckRequest struct {
	GidLookupInput      *string                 `json:"lookupInput,omitempty"`
	ScheduleLookupInput *string                 `json:"scheduleInput,omitempty"`
	Heading             string                  `json:"heading,omitempty"`
	OverrideSettings    *CreateOverrideSettings `json:"overrideSettings,omitempty"`
	InitialState        GotState                `json:"state,omitempty"`
}

func (c *CreateBuckRequest) HasOverride() bool {
	return c.OverrideSettings != nil
}

func NewCreateBuckRequest(lookup *GidLookup, dateLookup *DateLookup, heading string, state GotState, overrides *CreateOverrideSettings) CreateBuckRequest {
	var gidLookupString *string = nil
	var scheduleLookup *string = nil
	if lookup != nil {
		gidLookupString = &lookup.Input
	}
	if dateLookup != nil {
		scheduleLookup = &dateLookup.UserInput
	}
	return CreateBuckRequest{
		GidLookupInput:      gidLookupString,
		ScheduleLookupInput: scheduleLookup,
		Heading:             heading,
		OverrideSettings:    overrides,
		InitialState:        state,
	}
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

type PathString struct {
	DisplayPath string //either ID or alias for each node
	IdPath      string //path but based on ids only
}

type Shortcut struct {
	Display string //alias or id
	Id      string //just the id
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

// VX:TODO replace with GotId
type Gid struct {
	Id string
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

func (g *GotId) DisplayAasci() string {
	return "0" + g.AasciValue
}

func NewGotIdFromInt(intValue int32) (*GotId, error) {
	strVal, err := bullet_stl.BulletIdIntToAasci(int64(intValue))
	if err != nil {
		return nil, err
	}
	gotId := NewCompleteId(strVal, intValue)
	return &gotId, nil

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
