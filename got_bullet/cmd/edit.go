package cmd

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
	"vixac.com/got/console"
	"vixac.com/got/engine"
)

func buildEditCommand(deps RootDependencies) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "edit <lookup> <heading>",
		Short: "Edit a tasks title",
		//Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				err := errors.New("missing args")
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err
			}

			lookup := args[0]
			heading := strings.Join(args[1:], " ")
			if heading == "" {
				err := errors.New("missing heading")
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err

			}

			err := deps.Engine.EditTitle(engine.GidLookup{Input: lookup}, heading)
			if err != nil {
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err
			}
			return nil
		},
	}
	return cmd
}
