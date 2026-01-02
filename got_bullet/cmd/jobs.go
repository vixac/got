package cmd

import (
	"github.com/spf13/cobra"
	"vixac.com/got/console"
	"vixac.com/got/engine"
	"vixac.com/got/engine/bullet_engine"
)

// VX:TODO test
func buildJobsCommand(deps RootDependencies) *cobra.Command {

	var jobsCmd = &cobra.Command{
		Use:   "jobs",
		Short: "fetch jobs under gid",
		Run: func(cmd *cobra.Command, args []string) {
			//VX:TODO no more -u, lets got got jobs <id> optinoal.
			if len(args) > 1 {
				deps.Printer.Error(console.Message{Message: "jobs takes an optional <lookup> and thats it"})
				return
			}

			var lookup *engine.GidLookup = nil
			if len(args) == 1 {
				lookup = &engine.GidLookup{
					Input: args[0],
				}
			} else {
				lookup = &engine.GidLookup{
					Input: "",
				}
			}

			states := []int{engine.Active}

			res, err := deps.Engine.FetchItemsBelow(lookup, engine.AllDescendants, states)
			if err != nil {
				deps.Printer.Error(console.Message{Message: err.Error()})
				return
			}
			if res == nil || len(res.Result) == 0 {
				deps.Printer.Print(console.Message{Message: "no items found"})
				return
			}

			deps.Printer.Print(console.Message{Message: "\n-----------------------------------------\n\n"})
			table, err := bullet_engine.NewTable(res.Result, false)
			if err != nil {
				deps.Printer.Error(console.Message{Message: err.Error()})
				return
			}
			table.Render(deps.Printer, &console.GotTheme{})

		},
	}
	return jobsCmd

}
