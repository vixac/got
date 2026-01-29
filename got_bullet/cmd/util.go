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

func renderTable(lookup *engine.GidLookup, states []engine.GotState, options bullet_engine.TableRenderOptions, deps RootDependencies) {
	res, err := deps.Engine.FetchItemsBelow(lookup, options.SortByPath, states)
	if err != nil {
		deps.Printer.Error(console.Message{Message: err.Error()})
		return
	}
	if res == nil || len(res.Result) == 0 {
		deps.Printer.Print(console.Message{Message: "no items found"})
		return
	}

	var sections [][]engine.GotItemDisplay
	if options.GroupByTimeFrame {

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
				fmt.Printf("VX: Unhandled error parsing faulty date.")
				return
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
			sections = append(sections, theRestItems)
		}
		if len(lastWeekItems) > 0 {
			sections = append(sections, lastWeekItems)
		}
		if len(yesterdayItems) > 0 {
			sections = append(sections, yesterdayItems)
		}
		if len(todayItems) > 0 {
			sections = append(sections, todayItems)
		}

	} else {
		sections = append(sections, res.Result)
		if res.Parent != nil {
			sections = append(sections, []engine.GotItemDisplay{*res.Parent})
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
