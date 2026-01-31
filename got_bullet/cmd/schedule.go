package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"vixac.com/got/console"
	"vixac.com/got/engine"
)

func buildScheduleCommand(deps RootDependencies) *cobra.Command {

	var n = false
	var now *bool = &n
	cmd := &cobra.Command{
		Use:   "schedule <id> <date>",
		Short: "set the deadline for an id",
		//Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if now == nil && len(args) != 2 {
				err := errors.New("missing args")
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err
			}
			if now == nil && len(args) != 1 {
				err := errors.New("missing gid to schedule for now")
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err

			}

			id := args[0]

			var dateLookup *engine.DateLookup = nil
			if *now {
				d := engine.NowDateLookup()
				dateLookup = &d
			} else {
				dateStr := args[1]
				dateLookup = &engine.DateLookup{UserInput: dateStr}
			}

			err := deps.Engine.ScheduleItem(engine.GidLookup{Input: id}, *dateLookup)
			if err != nil {
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err
			}
			return nil

		},
	}
	cmd.Flags().BoolVarP(now, "now", "n", false, "Whether this is scheduled for now.")
	return cmd
}
