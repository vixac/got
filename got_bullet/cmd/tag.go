package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"vixac.com/got/console"
	"vixac.com/got/engine"
)

func buildTagCommand(deps RootDependencies) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "tag <id> <tag>",
		Short: "add a tag to an id.",
		//Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				err := errors.New("missing args")
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err
			}

			id := args[0]
			tagLookup := args[1]
			//VX:TODO convert to lookupType?
			err := deps.Engine.TagItem(engine.GidLookup{Input: id}, engine.TagLookup{Input: tagLookup})
			if err != nil {
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err
			}
			deps.Printer.Print(console.Message{Message: "Tagged it."})
			return nil

		},
	}
	return cmd
}
