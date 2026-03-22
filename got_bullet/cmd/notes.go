package cmd

import (
	"github.com/spf13/cobra"
	"vixac.com/got/console"
	"vixac.com/got/engine"
)

func buildNotesCommand(deps RootDependencies) *cobra.Command {

	var jobsCmd = &cobra.Command{
		Use:   "notes",
		Short: "fetch notes under gid",
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
			}
			renderNotesFor(*lookup, deps)
			/*
				//VX:TODO only real ids for now.
				//VX:TODO to the printer
				res, err := deps.Engine.NotesFor(*lookup)
				if err != nil {
					deps.Printer.Error(console.Message{Message: err.Error()})
					return
				}
				if res == nil {
					deps.Printer.Print(console.Message{Message: "There were no notes."})
					return
				}
				for _, v := range res.Blocks {
					fmt.Printf("VX: its a block: '%s'\n", v.Content)
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
			*/
		},
	}
	return jobsCmd

}
