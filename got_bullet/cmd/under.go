package cmd

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
	"vixac.com/got/console"
	"vixac.com/got/engine"
)

func buildUnderCommand(deps RootDependencies) *cobra.Command {

	var n = false
	var now *bool = &n
	cmd := &cobra.Command{
		Use:   "under <alias> <heading>",
		Short: "Create a task with no deadline, under provided parent",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				err := errors.New("missing args")
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err
			}
			parentAlias := args[0]
			heading := strings.Join(args[1:], " ")

			if parentAlias == "" {
				err := errors.New("missing alias")
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err
			}
			if heading == "" {
				err := errors.New("missing heading")
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err
			}
			var dateLookup *string = nil
			//VX:Note we should support other datelookups here too.
			if *now {
				d := engine.NowDateLookup()
				dateLookup = &d.UserInput
			}
			req := engine.CreateBuckRequest{
				Heading:             heading,
				ScheduleLookupInput: dateLookup,
				GidLookupInput:      &parentAlias,
			}
			id, err := deps.Engine.CreateBuck(req)
			if err != nil {
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err
			}
			deps.Printer.Print(console.Message{Message: id.AasciValue})

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
