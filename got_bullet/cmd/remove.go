package cmd

import (
	"github.com/spf13/cobra"
	"vixac.com/got/console"
	"vixac.com/got/engine"
)

func buildRemoveCommand(deps RootDependencies) *cobra.Command {
	var doneCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove an item",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				deps.Printer.Error(console.Message{Message: "Expected the alias as input"})
				return
			}
			err := deps.Engine.Delete(engine.GidLookup{Input: args[0]})
			if err != nil {
				deps.Printer.Error(console.Message{Message: err.Error()})
				return
			}
			msg := "Success: " + args[0] + " is removed."
			deps.Printer.Print(console.Message{Message: msg})
		},
	}
	return doneCmd
}
