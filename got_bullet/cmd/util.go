package cmd

import (
	"vixac.com/got/console"
	"vixac.com/got/engine"
	"vixac.com/got/engine/bullet_engine"
)

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
		//VX:TODO
		var today []engine.GotItemDisplay
		var yesterday []engine.GotItemDisplay
		var lastWeek []engine.GotItemDisplay
		var theRest []engine.GotItemDisplay
		for _, r := range res.Result {

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
