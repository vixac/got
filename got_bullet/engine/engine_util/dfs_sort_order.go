package engine_util

import (
	"fmt"
	"sort"

	"vixac.com/got/engine"
)

type PathPair struct {
	PathString string
	Item       engine.GotItemDisplay
}

// this converts the paths into strings and uses those strings to to lexographical sort. The outcome is that you get depth first printing of the a parent
// followed by all its descendants, recursively
func SortTheseIntoDFS(items []engine.GotItemDisplay) []engine.GotItemDisplay {
	var sortablePaths []PathPair
	for _, item := range items {
		fullPath := item.FullPathString()
		sortablePaths = append(sortablePaths, PathPair{
			PathString: fullPath.IdPath,
			Item:       item,
		})
	}
	sort.Slice(sortablePaths, func(i, j int) bool {
		return sortablePaths[i].PathString < sortablePaths[j].PathString //this is actually a bit of a hack. the length of the path is nearly right but the sort will on occasion be wrong based on lengths.
	})

	var backToItems []engine.GotItemDisplay
	for _, i := range sortablePaths {
		backToItems = append(backToItems, i.Item)
	}
	return backToItems

}

func SortByUpdated(items []engine.GotItemDisplay) []engine.GotItemDisplay {
	var sortableItems []engine.GotItemDisplay = items
	sort.Slice(sortableItems, func(i, j int) bool {
		return updatedSort(sortableItems[i].SummaryObj, sortableItems[j].SummaryObj)
	})
	return sortableItems
}

func updatedSort(lhsSummary *engine.Summary, rhsSummary *engine.Summary) bool {
	if lhsSummary == nil || rhsSummary == nil { //VX:Note this should not happen.
		return true
	}
	lhs := lhsSummary.UpdatedDate
	rhs := rhsSummary.UpdatedDate
	if lhs == nil && rhs == nil {
		return true
	}
	if lhs == nil {
		return true
	}
	if rhs == nil {
		return false
	}
	return lhs.EpochMillis() < rhs.EpochMillis()
}

func deadlineSort(lhsSummary *engine.Summary, rhsSummary *engine.Summary) bool {
	if lhsSummary == nil || rhsSummary == nil { //VX:Note this should not happen.
		return true
	}
	lhs := lhsSummary.Deadline
	rhs := rhsSummary.Deadline
	if lhs == nil && rhs == nil { //if there are no deadlines, sort by updated.
		return updatedSort(lhsSummary, rhsSummary)
	}
	if lhs == nil {
		return true
	}
	if rhs == nil {
		return false
	}
	if lhs.IsNow() && rhs.IsNow() {
		return true
	}
	if lhs.IsNow() && !rhs.IsNow() {
		return false
	}
	if !lhs.IsNow() && rhs.IsNow() {
		return true
	}
	return lhs.EpochMillis() > rhs.EpochMillis()
}

func SortByDeadline(items []engine.GotItemDisplay) []engine.GotItemDisplay {
	fmt.Printf("VX: Sorteing by deadline \n")
	var sortableItems []engine.GotItemDisplay = items

	sort.Slice(sortableItems, func(i, j int) bool {
		return deadlineSort(sortableItems[i].SummaryObj, sortableItems[j].SummaryObj)
	})
	return sortableItems
}
