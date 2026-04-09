package engine

import "vixac.com/got/console"

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
	IsParent      bool
}

func (i *GotItemDisplay) IsCollapsed() bool {
	return i.SummaryObj != nil && i.SummaryObj.Flags != nil && i.SummaryObj.Flags["collapsed"] == true
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
