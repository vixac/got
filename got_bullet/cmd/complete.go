package cmd

import (
	"github.com/spf13/cobra"
	"vixac.com/got/console"
	"vixac.com/got/engine"
	"vixac.com/got/engine/bullet_engine"
)

// VX:TODO test
func buildCompleteCommand(deps RootDependencies) *cobra.Command {

	var jobsCmd = &cobra.Command{
		Use:   "completed",
		Short: "fetch completed jobs under gid",
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

			states := []engine.GotState{engine.Complete}
			options := bullet_engine.TableRenderOptions{
				FlatPaths:          true,
				ShowCreatedColumn:  true,
				ShowUpdatedColumn:  true,
				SortByPath:         false, //for now this means sort by updated, because we dont have more options yet.
				HideUnderCollapsed: false,
			}
			renderTable(lookup, states, options, deps)
		},
	}
	return jobsCmd

}
