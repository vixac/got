package cmd

import (
	"fmt"
	"time"

	"vixac.com/got/console"
	"vixac.com/got/engine"
	"vixac.com/got/engine/bullet_engine"
)

func (e *EngineBullet) addDailyDividers(items []engine.GotItemDisplay) []engine.GotItemDisplay {

	var finalResult []engine.GotItemDisplay
	//VX:TODO go to midnight instead of now:

	midnight := LastMidnightUTC()
	today := midnight.Unix()
	yesterday := midnight.AddDate(0, 0, -1).Unix()
	lastWeek := midnight.AddDate(0, 0, -7).Unix()

	var lastWeekInPlace = false
	var yesterdayInPlace = false
	var todayInPlace = false
	for _, i := range items {
		dateOfThisItem, err := i.SummaryObj.UpdatedDate.ToDate()
		timeOfThisItem := time.Time(*dateOfThisItem).Unix()
		//itemDate := Date(dateOfThisItemStr)
		if err != nil {
			fmt.Printf("Error dating this item. Not adding daily dividers. This is a quiet failure.")
			return items
		}
		if !lastWeekInPlace {
			if timeOfThisItem > lastWeek && timeOfThisItem < yesterday {

				lastWeekInPlace = true

			}

		} else if !yesterdayInPlace {

		} else if !todayInPlace {

		}
		finalResult = append(finalResult, i)
	}
	return finalResult

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
		//VX:TODO
		var todayItems []engine.GotItemDisplay
		var yesterdayItems []engine.GotItemDisplay
		var lastWeekItems []engine.GotItemDisplay
		var theRestItems []engine.GotItemDisplay
		for _, r := range res.Result {
			r.SummaryObj.UpdatedDate

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
