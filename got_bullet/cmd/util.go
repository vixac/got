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

	//deps.Printer.Print(console.Message{Message: "\n-----------------------------------------\n\n"})

	table, err := bullet_engine.NewTable(res, options)
	if err != nil {
		deps.Printer.Error(console.Message{Message: err.Error()})
		return
	}
	table.Render(deps.Printer, &console.GotTheme{})

}
