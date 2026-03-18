package cmd

import (
	"github.com/spf13/cobra"
	"vixac.com/got/console"
)

func buildConsumeRestoreCommand(deps RootDependencies) *cobra.Command {

	var jobsCmd = &cobra.Command{
		Use:   "restore",
		Short: "reads a restore file and absorbs it.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				deps.Printer.Error(console.Message{Message: "path to restore file required."})
				return
			}
			filenmame := args[0]

			err := deps.Engine.RestoreFromFile(filenmame)
			if err != nil {
				deps.Printer.Error(console.Message{Message: err.Error()})
			}
		},
	}
	return jobsCmd

}
