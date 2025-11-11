package cmd

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
	"vixac.com/got/console"
	"vixac.com/got/engine"
)

func buildUnderCommand(deps RootDependencies) *cobra.Command {
	var parentAlias string

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
			_, err := deps.Engine.CreateBuck(&engine.GidLookup{Input: parentAlias},
				nil,
				true,
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
