package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"vixac.com/got/console"
	"vixac.com/got/engine"
)

func buildNotesCommand(deps RootDependencies) *cobra.Command {

	var r = false
	var recurse *bool = &r
	var cmd = &cobra.Command{
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
			if *recurse == true {
				//VX:TODO fix recursion.
				fmt.Printf("VX: Apologies, recurse option is not ready. It request FetchItemsBelow to also return the parent node. Otherwise recurse doesnt include the node you just looked up.")
				return
			}
			renderNotesFor(lookup, *recurse, deps)
		},
	}
	cmd.Flags().BoolVarP(recurse, "recurse", "r", false, "Fetch all notes for nodes under this one too.")
	return cmd

}
