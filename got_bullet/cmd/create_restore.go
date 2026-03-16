package cmd

import (
	"github.com/spf13/cobra"
	"vixac.com/got/console"
)

func buildRestoreCommand(deps RootDependencies) *cobra.Command {

	var jobsCmd = &cobra.Command{
		Use:   "create-restore",
		Short: "builds a restore file",
		Run: func(cmd *cobra.Command, args []string) {
			err := deps.Engine.CreateStoreFile()
			if err != nil {
				deps.Printer.Error(console.Message{Message: err.Error()})
			}
		},
	}
	return jobsCmd

}
