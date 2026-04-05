package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"vixac.com/got/console"
	"vixac.com/got/engine"
)

func buildMvCommand(deps RootDependencies) *cobra.Command {
	var newParent string
	var cmd = &cobra.Command{
		Use:   "move",
		Short: "move <leaf> <newParent>",
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) != 1 {
				err := errors.New("Invalid args. Just 1 please.")
				deps.Printer.Error(console.Message{Message: err.Error()})
				return
			}
			target := args[0]

			err := deps.Engine.Move(engine.GidLookup{Input: target}, engine.GidLookup{Input: newParent})
			if err != nil {
				deps.Printer.Error(console.Message{Message: err.Error()})
				return
			}

			var msg string
			msg = "Success: " + target + " moved to new parent " + newParent

			deps.Printer.Print(console.Message{Message: msg})
		},
	}
	cmd.Flags().StringVarP(&newParent, "under", "u", "", "The new parent")
	return cmd
}
