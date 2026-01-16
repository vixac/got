package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"vixac.com/got/console"
	"vixac.com/got/engine"
)

func buildMvCommand(deps RootDependencies) *cobra.Command {
	var moveCmd = &cobra.Command{
		Use:   "move",
		Short: "move <leaf> <newParent>",
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) != 2 {
				err := errors.New("Invalid args. just 2 please.")
				deps.Printer.Error(console.Message{Message: err.Error()})
				return
			}
			target := args[0]
			newParent := args[1]

			oldParent, err := deps.Engine.Move(engine.GidLookup{Input: target}, engine.GidLookup{Input: newParent})
			if err != nil {
				deps.Printer.Error(console.Message{Message: err.Error()})
				return
			}

			var msg string
			if oldParent == nil {
				msg = "Success: " + target + " moved to new parent " + newParent
			} else {
				msg = "Success: " + target + " moved from old parent '" + oldParent.Title + "' to " + newParent
			}
			deps.Printer.Print(console.Message{Message: msg})
		},
	}
	return moveCmd
}
