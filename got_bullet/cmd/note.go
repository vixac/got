package cmd

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
	"vixac.com/got/console"
	"vixac.com/got/engine"
)

func buildJotCommand(deps RootDependencies) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "note <gid> <note>",
		Short: "Create a note under the given gid",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				err := errors.New("missing args")
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err
			}
			parentLookup := args[0]
			note := strings.Join(args[1:], " ")

			if parentLookup == "" {
				err := errors.New("missing alias")
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err
			}
			if note == "" {
				err := errors.New("missing note")
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err
			}
			_, err := deps.Engine.JotNote(engine.GidLookup{Input: parentLookup}, note)
			if err != nil {
				deps.Printer.Error(console.Message{Message: err.Error()})
				return err
			}
			return nil
		},
	}
	return cmd
}
