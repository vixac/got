package cmd

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
	"vixac.com/got/console"
	"vixac.com/got/engine"
)

func buildEventCommand(deps RootDependencies) *cobra.Command {
	var parentAlias string

	cmd := &cobra.Command{
		Use:   "event <date> [-for <alias>] <heading>",
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
			var lookup *engine.GidLookup = nil
			if parentAlias != "" {
				lookup = &engine.GidLookup{Input: parentAlias}
			}
			_, err := deps.Engine.CreateBuck(lookup,
				&engine.DateLookup{UserInput: date},
				false, //this is the only difference btween till and event
				heading,
			)
			if err != nil {
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&parentAlias, "for", "f", "", "Alias to assign the task under")
	return cmd
}
