package engine

import bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/ids"

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
	FetchItemsBelow(ookup *GidLookup, descendantType int, states []int) (*GotFetchResult, error)
}

// The interface for all aliasing functionality
type GotAliasInterface interface {
	Lookup(alias string) (*NodeId, error)
	Unalias(alias string) (*NodeId, error)
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
	Path     GotPath
}

// VX:TODO replace with GotId
type Gid struct {
	Id string
}
type NodeId struct {
	Gid   Gid
	Title string
	Alias string
}
type GotPath struct {
	Ancestry []NodeId
}

// VX:TODO consider using BUlletId semantics to lazy compute these
type GotId struct {
	AasciValue string
	IntValue   int64
}

func NewCompleteId(aasci string, intValue int64) GotId {
	return GotId{
		AasciValue: aasci,
		IntValue:   intValue,
	}
}

// this is basically a wrapper for BulletId
func NewGotId(aasci string) (*GotId, error) {
	intVal, err := bullet_stl.AasciBulletIdToInt(aasci)
	if err != nil {
		return nil, err
	}
	return &GotId{
		AasciValue: aasci,
		IntValue:   intVal,
	}, nil
}

//Lets be super clear about the ides
