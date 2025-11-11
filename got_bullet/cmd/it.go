package cmd

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
	"vixac.com/got/console"
)

func buildItCommand(deps RootDependencies) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "it <heading>",
		Short: "Create a note with no parent and no date",
		//Args:  cobra.MinimumNArgs(2),
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
			_, err := deps.Engine.CreateBuck(nil,
				nil,
				false, //this is the only differece between to and it
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
