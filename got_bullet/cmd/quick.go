package cmd

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
	"vixac.com/got/console"
	"vixac.com/got/engine"
)

func buildQuickCommand(deps RootDependencies) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "quick <date> <heading>",
		Short: "Create a till task with a due date and optional alias",
		//Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				err := errors.New("missing args")
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err
			}
			date := args[0]
			heading := strings.Join(args[1:], " ")
			println("QUICK: date is " + date + " and heading is " + heading)

			if date == "" {
				err := errors.New("missing date")
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err
			}
			if heading == "" {
				err := errors.New("missing heading")
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err
			}
			_, err := deps.Engine.CreateBuck(nil,
				&engine.DateLookup{UserInput: date},
				true, //this is the only differece between til and event
				heading,
			)
			if err != nil {
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err
			}
			return nil
		},
	}
	return cmd
}
