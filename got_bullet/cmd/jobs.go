package cmd

import (
	"github.com/spf13/cobra"
	"vixac.com/got/console"
	"vixac.com/got/engine"
	"vixac.com/got/engine/engine_util"
)

// VX:TODO test
func buildJobsCommand(deps RootDependencies) *cobra.Command {

	byDeadline := false
	var sortByDeadline *bool = &byDeadline
	var cmd = &cobra.Command{
		Use:   "jobs",
		Short: "fetch jobs under gid",
		Run: func(cmd *cobra.Command, args []string) {
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
			sortStyle := engine_util.SortByPath
			if *sortByDeadline {
				sortStyle = engine_util.SortByDeadlineDate
			}

			states := []engine.GotState{engine.Active, engine.Note}
			options := engine_util.TableRenderOptions{
				FlatPaths:          *sortByDeadline, //flat paths for deadline
				ShowCreatedColumn:  true,
				ShowUpdatedColumn:  false,
				SortStyle:          sortStyle,
				HideUnderCollapsed: true,
			}
			renderTable(lookup, states, options, deps)
		},
	}
	cmd.Flags().BoolVarP(sortByDeadline, "now", "n", false, "Whether to sort by now or not.")
	return cmd

}
