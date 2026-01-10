package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"vixac.com/got/console"
	"vixac.com/got/engine"
)

func buildAliasCommand(deps RootDependencies) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "alias",
		Short: "alias an item with a better name. <alias> <gid>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				err := errors.New("missing args. Pass in an alias and a gid")
				deps.Printer.Error(console.Message{Message: err.Error()})
				return
			}
			alias := args[0]
			gid := args[1]

			if alias == "" {
				deps.Printer.Error(console.Message{Message: "Missing alias"})
				return
			}

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
	return cmd
}
