package engine

import (
	"errors"
	"math"

	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/ids"
)

// / This is the machine that takes the commands, changes the backend state and returns wahts requested.
type GotEngine interface {
	Summary(lookup *GidLookup) (*GotSummary, error)
	Resolve(lookup GidLookup) (*NodeId, error)
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

// states
const (
	Active   = 0
	Note     = 8000
	Complete = 16000
)

type GotFetchResult struct {
	Result []GotSummary
}

// All the lookup stuff
type GotFetchInterface interface {
	FetchItemsBelow(lookup *GidLookup, descendantType int, states []int) (*GotFetchResult, error)
}

// The interface for all aliasing functionality
type GotAliasInterface interface {
	Lookup(alias string) (*GotId, error)
	LookupAliasForGid(gid string) (*string, error)
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

type GotSummary struct {
	Gid      string
	Title    string
	Alias    string
	Deadline string
	Path     *GotPath
	NumberGo int
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
	//VX:TODO maybe title in here?
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
