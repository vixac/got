package cmd

import (
	"github.com/spf13/cobra"
	"vixac.com/got/console"
)

func buildUnaliasCommand(deps RootDependencies) *cobra.Command {
	var doneCmd = &cobra.Command{
		Use:   "unalias",
		Short: "Remove the alias of a gid",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				deps.Printer.Error(console.Message{Message: "Expected the alias as input"})
				return
			}
			_, err := deps.Engine.Unalias(args[0])
			if err != nil {
				deps.Printer.Error(console.Message{Message: err.Error()})
				return
			}
			msg := "Success: " + args[0] + " is unaliased."
			deps.Printer.Print(console.Message{Message: msg})
		},
	}
	return doneCmd
}
