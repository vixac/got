package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"vixac.com/got/console"
	"vixac.com/got/engine"
)

func buildScheduleCommand(deps RootDependencies) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "schedule <id> <date>",
		Short: "Create a note with no parent and no date",
		//Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				err := errors.New("missing args")
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err
			}

			id := args[0]
			dateStr := args[1]

			err := deps.Engine.ScheduleItem(engine.GidLookup{Input: id}, engine.DateLookup{UserInput: dateStr})
			if err != nil {
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err
			}
			return nil

		},
	}
	return cmd
}
