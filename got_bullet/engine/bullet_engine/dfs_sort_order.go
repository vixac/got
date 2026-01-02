package bullet_engine

import (
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
			PathString: fullPath,
			Item:       item,
		})
	}
	sort.Slice(sortablePaths, func(i, j int) bool {
		return sortablePaths[i].PathString < sortablePaths[j].PathString
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
		return sortableItems[i].Updated < sortableItems[j].Updated
	})
	return sortableItems
}
