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

			states := []engine.GotState{engine.Active, engine.Note}
			options := bullet_engine.TableRenderOptions{
				FlatPaths:          false,
				ShowCreatedColumn:  true,
				ShowUpdatedColumn:  true,
				SortByPath:         true,
				HideUnderCollapsed: true,
			}
			renderTable(lookup, states, options, deps)
		},
	}
	return jobsCmd

}
