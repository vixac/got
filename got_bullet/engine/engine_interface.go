package engine

// / This is the machine that takes the commands, changes the backend state and returns wahts requested.
type GotEngine interface {
	Summary(gid GidLookup) (*GotSummary, error)
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
