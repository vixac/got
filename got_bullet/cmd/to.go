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

			var dateLookup *engine.DateLookup = nil
			//VX:Note we should support other datelookups here too.
			if *now {
				d := engine.NowDateLookup()
				dateLookup = &engine.DateLookup{
					UserInput: d.UserInput,
				}
			}
			heading := strings.Join(args, " ")
			if heading == "" {
				err := errors.New("missing heading")
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err
			}
			req := engine.NewCreateBuckRequest(nil, dateLookup, heading, engine.Active, nil)
			id, err := deps.Engine.CreateBuck(req)
			if err != nil {
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err
			}
			if id == nil {
				deps.Printer.Print(console.Message{Message: "VX: dev error. Wheres the id"})
				return nil
			}

			deps.Printer.Print(console.Message{Message: id.DisplayAasci()})
			return nil
		},
	}
	cmd.Flags().BoolVarP(now, "now", "n", false, "Whether this is scheduled for now.")

	return cmd
}
