package cmd

import (
	"fmt"
	"time"

	"vixac.com/got/console"
	"vixac.com/got/engine"
	"vixac.com/got/engine/engine_util"
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

func setionsByTimeframe(res *engine.GotFetchResult) ([]engine_util.TableSection, error) {
	var sections []engine_util.TableSection
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
		sections = append(sections, engine_util.TableSection{Name: "", Items: theRestItems})
	}
	if len(lastWeekItems) > 0 {
		sections = append(sections, engine_util.TableSection{Name: " Last Week ", Items: lastWeekItems})
	}
	if len(yesterdayItems) > 0 {
		sections = append(sections, engine_util.TableSection{Name: " Yesterday ", Items: yesterdayItems})
	}
	if len(todayItems) > 0 {
		sections = append(sections, engine_util.TableSection{Name: " Today ", Items: todayItems})
	}
	return sections, nil
}

// creates a section for each of the top level siblings.
func sectionsByTopLevelSiblings(res *engine.GotFetchResult) ([]engine_util.TableSection, error) {
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
	var sections []engine_util.TableSection
	//this is a flat response, so don't break it up into sections
	if minDepth == maxDepth {
		sections = append(sections, engine_util.TableSection{Name: "", Items: res.Result})
		return sections, nil

	}

	//everytime we reach an item of minDepth (that is a top level node relative to this search), we start a new section
	var currentSection []engine.GotItemDisplay
	for _, r := range res.Result {
		if r.Path.Depth() == minDepth && len(currentSection) > 0 {
			//flush the current section and start a new one
			sections = append(sections, engine_util.TableSection{Name: "", Items: currentSection})
			currentSection = []engine.GotItemDisplay{}
		}

		currentSection = append(currentSection, r)
	}
	if len(currentSection) > 0 {
		sections = append(sections, engine_util.TableSection{Name: "", Items: currentSection})
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
	var finalSections []engine_util.TableSection

	if len(leafSection) > 0 {
		finalSections = append(finalSections, engine_util.TableSection{Name: "", Items: leafSection})
	}
	for _, s := range squashedSections {
		finalSections = append(finalSections, engine_util.TableSection{Name: "", Items: s})
	}

	return finalSections, nil
}

func renderNotesFor(lookup *engine.GidLookup, recurse bool, deps RootDependencies) {
	notes, err := deps.Engine.NotesFor(lookup, recurse)
	if err != nil {
		deps.Printer.Error(console.Message{Message: err.Error()})
		return
	}
	if notes == nil {
		deps.Printer.Error(console.Message{Message: "No notes to render."})
		return
	}

	maxTitleLen := 90
	var displays []engine.GotItemDisplay

	var allFetchedGotIds = make(map[engine.GotId]bool)
	for _, block := range notes.Blocks {
		id := block.Id.GotId
		allFetchedGotIds[id] = true
	}
	var idStrings []string
	for k, _ := range allFetchedGotIds {
		idStrings = append(idStrings, k.AasciValue)
	}

	aliases, err := deps.Engine.LookupAliasForMany(idStrings)

	now := time.Now()
	for _, block := range notes.Blocks {
		var truncatedContent = ""
		blockTitleLen := len(block.Content)
		loopLen := maxTitleLen
		var showDotDotDot = false
		if blockTitleLen < maxTitleLen {
			loopLen = blockTitleLen
		}

		for j := 0; j < loopLen; j++ {
			char := block.Content[j]
			if char == '\n' {
				showDotDotDot = true
				continue
			}
			truncatedContent += string(block.Content[j])
		}
		if blockTitleLen > maxTitleLen || showDotDotDot {
			truncatedContent += " ..."
		}

		var theAlias = ""
		alias, ok := aliases[block.Id.GotId.AasciValue]
		if ok {
			theAlias = *alias
		}

		pathItem := engine.PathItem{
			Id:    block.Id.ToString(),
			Alias: nil,
		}
		relativeEditedDate, _ := console.HumanizeDate(block.Edited, now)
		if err != nil {
			deps.Printer.Error(console.Message{Message: err.Error()})
			return
		}
		//VX:TODO review the rendering of this. It's currently way off.
		display := engine.GotItemDisplay{
			GotId:         block.Id.GotId,
			DisplayGid:    block.Id.ToString(),
			Title:         truncatedContent,
			Path:          &engine.GotPath{Ancestry: []engine.PathItem{pathItem}}, //janky. not ancestry. not a path.
			Alias:         theAlias,
			SummaryObj:    nil,
			HasTNote:      false,
			Deadline:      "",
			DeadlineToken: console.TokenBrand{},
			Created:       block.Created().String(),
			Updated:       relativeEditedDate,
		}
		displays = append(displays, display)
	}
	section := engine_util.TableSection{
		Name:  "Notes",
		Items: displays,
	}

	table, err := engine_util.NewTable(&engine_util.GotTableSections{
		Sections: []engine_util.TableSection{section},
	}, engine_util.TableRenderOptions{
		FlatPaths:         true,
		HideNumberGo:      true,
		ShowUpdatedColumn: true,
	})
	if err != nil {
		deps.Printer.Error(console.Message{Message: err.Error()})
		return
	}
	table.Render(deps.Printer, &console.GotTheme{})

}

func renderTable(lookup *engine.GidLookup, states []engine.GotState, options engine_util.TableRenderOptions, deps RootDependencies) {
	res, err := deps.Engine.FetchItemsBelow(lookup, options.SortByPath, states, options.HideUnderCollapsed)
	if err != nil {
		deps.Printer.Error(console.Message{Message: err.Error()})
		return
	}
	if res == nil || (len(res.Result) == 0 && res.Parent == nil) {
		deps.Printer.Print(console.Message{Message: "No items found."})
		return
	}

	var sections []engine_util.TableSection
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
			sections = append(sections, engine_util.TableSection{Name: "", Items: []engine.GotItemDisplay{*res.Parent}})
		}
	}

	table, err := engine_util.NewTable(&engine_util.GotTableSections{
		Sections: sections,
	}, options)
	if err != nil {
		deps.Printer.Error(console.Message{Message: err.Error()})
		return
	}
	table.Render(deps.Printer, &console.GotTheme{})
}
