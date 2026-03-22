package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"vixac.com/got/console"
)

func buildTCommand(deps RootDependencies) *cobra.Command {
	var doneCmd = &cobra.Command{
		Use:   "t",
		Short: "edit content",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				deps.Printer.Error(console.Message{Message: "Expected at least one lookup as input"})
				return
			}
			fmt.Printf("VX: open then timesptamp is paused..")
			//			lookup := engine.GidLookup{Input: args[0]}

			/*	err := deps.Engine.OpenThenTimestamp(lookup)
				if err != nil {
					deps.Printer.Print(console.Message{Message: "error: " + err.Error()})
				} else {
					deps.Printer.Print(console.Message{Message: "<>"})
				}
			*/

			//msg := "Success: " + args[0] + " is marked complete."
			//deps.Printer.Print(console.Message{Message: msg})
		},
	}
	return doneCmd
}
