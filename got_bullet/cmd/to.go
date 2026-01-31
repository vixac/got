package cmd

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
	"vixac.com/got/console"
	"vixac.com/got/engine"
)

func buildToCommand(deps RootDependencies) *cobra.Command {

	var n = false
	var now *bool = &n
	cmd := &cobra.Command{
		Use:   "to <heading>",
		Short: "Create a task with no parent and no date",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				err := errors.New("missing args")
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err
			}

			heading := strings.Join(args, " ")
			if heading == "" {
				err := errors.New("missing heading")
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err
			}
			var dateLookup *engine.DateLookup = nil
			if *now {
				d := engine.NowDateLookup()
				dateLookup = &d
			}
			_, err := deps.Engine.CreateBuck(nil,
				dateLookup,
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
	cmd.Flags().BoolVarP(now, "now", "n", false, "Whether this is scheduled for now.")

	return cmd
}
