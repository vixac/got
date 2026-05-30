package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"vixac.com/got/console"
	"vixac.com/got/engine"
)

func buildAliasCommand(deps RootDependencies) *cobra.Command {
	var alias string
	var cmd = &cobra.Command{
		Use:   "alias",
		Short: "alias an item with a better name. <alias> <gid>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				err := errors.New("Invalid args. Just 1 please.")
				deps.Printer.Error(console.Message{Message: err.Error()})
				return
			}
			gid := args[0]

			if gid == "" {
				deps.Printer.Error(console.Message{Message: "Missing gid"})
				return
			}
			lookup := engine.GidLookup{Input: gid}
			_, err := deps.Engine.Alias(lookup, alias)
			if err != nil {
				deps.Printer.Error(console.Message{Message: err.Error()})
				return
			}
			msg := "Success: " + gid + " is now aliased to " + alias + "."
			deps.Printer.Print(console.Message{Message: msg})
		},
	}
	cmd.Flags().StringVarP(&alias, "to", "t", "", "The new alias")
	return cmd
}
