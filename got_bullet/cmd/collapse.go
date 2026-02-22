package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
	"vixac.com/got/console"
	"vixac.com/got/engine"
)

func buildCollapseCommand(deps RootDependencies) *cobra.Command {
	var doneCmd = &cobra.Command{
		Use:   "collapse",
		Short: "mark an item collapsed so we don't see its details",
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

			err := deps.Engine.ToggleCollapse(*lookup, true)
			if err != nil {
				deps.Printer.Error(console.Message{Message: err.Error()})
				return
			}
			msg := "Success: " + strconv.Itoa(len(args)) + " items is collapsed."
			deps.Printer.Print(console.Message{Message: msg})
		},
	}
	return doneCmd
}
