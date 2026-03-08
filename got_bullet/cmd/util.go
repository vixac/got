package cmd

import (
	"fmt"
	"time"

	"vixac.com/got/console"
	"vixac.com/got/engine"
	"vixac.com/got/engine/bullet_engine"
)

func LastMidnightUTC() time.Time {
	now := time.Now().UTC()
	return time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		time.UTC,
	)
}

func setionsByTimeframe(res *engine.GotFetchResult) ([]bullet_engine.TableSection, error) {
	var sections []bullet_engine.TableSection
	midnight := LastMidnightUTC()
	todayTime := midnight.Unix()
	yesterday := midnight.AddDate(0, 0, -1).Unix()
	lastWeek := midnight.AddDate(0, 0, -7).Unix()
	var todayItems []engine.GotItemDisplay
	var yesterdayItems []engine.GotItemDisplay
	var lastWeekItems []engine.GotItemDisplay
	var theRestItems []engine.GotItemDisplay
	for _, r := range res.Result {
		dateOfThisItem, err := r.SummaryObj.UpdatedDate.ToDate()
		if dateOfThisItem == nil {
			theRestItems = append(theRestItems, r)
			continue
		}
		timeOfThisItem := time.Time(*dateOfThisItem).Unix()
		if err != nil {
			return nil, err
		}
		if timeOfThisItem > todayTime {
			todayItems = append(todayItems, r)
		} else if timeOfThisItem > yesterday {
			yesterdayItems = append(yesterdayItems, r)
		} else if timeOfThisItem > lastWeek {
			lastWeekItems = append(lastWeekItems, r)
		} else {
			theRestItems = append(theRestItems, r)
		}
	}
	if len(theRestItems) > 0 {
		sections = append(sections, bullet_engine.TableSection{Name: "", Items: theRestItems})
	}
	if len(lastWeekItems) > 0 {
		sections = append(sections, bullet_engine.TableSection{Name: " Last Week ", Items: lastWeekItems})
	}
	if len(yesterdayItems) > 0 {
		sections = append(sections, bullet_engine.TableSection{Name: " Yesterday ", Items: yesterdayItems})
	}
	if len(todayItems) > 0 {
		sections = append(sections, bullet_engine.TableSection{Name: " Today ", Items: todayItems})
	}
	return sections, nil
}

// creates a section for each of the top level siblings.
func sectionsByTopLevelSiblings(res *engine.GotFetchResult) ([]bullet_engine.TableSection, error) {
	const MaxUint = ^uint(0)
	var minDepth int = int(MaxUint >> 1)
	var maxDepth int = 0
	for _, r := range res.Result {
		depth := r.Path.Depth()
		if minDepth > depth {
			minDepth = depth
		}
		if maxDepth < depth {
			maxDepth = depth
		}
	}
	var sections []bullet_engine.TableSection
	//this is a flat response, so don't break it up into sections
	if minDepth == maxDepth {
		sections = append(sections, bullet_engine.TableSection{Name: "", Items: res.Result})
		return sections, nil

	}

	//everytime we reach an item of minDepth (that is a top level node relative to this search), we start a new section
	var currentSection []engine.GotItemDisplay
	for _, r := range res.Result {
		if r.Path.Depth() == minDepth && len(currentSection) > 0 {
			//flush the current section and start a new one
			sections = append(sections, bullet_engine.TableSection{Name: "", Items: currentSection})
			currentSection = []engine.GotItemDisplay{}
		}

		currentSection = append(currentSection, r)
	}
	if len(currentSection) > 0 {
		sections = append(sections, bullet_engine.TableSection{Name: "", Items: currentSection})
	}

	//if theres a sequence of sections with only 1 item, we group them. This should act as though all leaf nodes are put into a single section.
	var squashedSections [][]engine.GotItemDisplay
	var leafSection []engine.GotItemDisplay
	for _, s := range sections {
		if len(s.Items) == 1 {
			leafSection = append(leafSection, s.Items[0])
		} else {
			squashedSections = append(squashedSections, s.Items)
		}
	}

	//here we create a new sections array and put the leaf nodes at the top. This is not performant.
	var finalSections []bullet_engine.TableSection

	if len(leafSection) > 0 {
		finalSections = append(finalSections, bullet_engine.TableSection{Name: "", Items: leafSection})
	}
	for _, s := range squashedSections {
		finalSections = append(finalSections, bullet_engine.TableSection{Name: "", Items: s})
	}

	return finalSections, nil
}

func renderTable(lookup *engine.GidLookup, states []engine.GotState, options bullet_engine.TableRenderOptions, deps RootDependencies) {
	res, err := deps.Engine.FetchItemsBelow(lookup, options.SortByPath, states, options.HideUnderCollapsed)
	if err != nil {
		deps.Printer.Error(console.Message{Message: err.Error()})
		return
	}
	if res == nil || (len(res.Result) == 0 && res.Parent == nil) {
		deps.Printer.Print(console.Message{Message: "No items found."})
		return
	}

	var sections []bullet_engine.TableSection
	if options.GroupByTimeFrame {
		timeFrameSections, err := setionsByTimeframe(res)
		if err != nil {
			fmt.Printf("VX: Unhandled error parsing faulty date.  %s", err.Error())
			return
		}
		sections = timeFrameSections
	} else {
		siblingSections, err := sectionsByTopLevelSiblings(res)
		if err != nil {
			fmt.Printf("VX: Unhandled error creating sections by siblings %s", err.Error())
			return
		}
		sections = siblingSections
		if res.Parent != nil {
			sections = append(sections, bullet_engine.TableSection{Name: "", Items: []engine.GotItemDisplay{*res.Parent}})
		}
	}

	table, err := bullet_engine.NewTable(&bullet_engine.GotTableSections{
		Sections: sections,
	}, options)
	if err != nil {
		deps.Printer.Error(console.Message{Message: err.Error()})
		return
	}
	table.Render(deps.Printer, &console.GotTheme{})
}
