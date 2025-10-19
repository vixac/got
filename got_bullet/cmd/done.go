package cmd

import (
	"github.com/spf13/cobra"
	"vixac.com/got/console"
)

func buildDoneCommand(messenger console.Messenger) *cobra.Command {
	var gid string
	var doneCmd = &cobra.Command{
		Use:   "done",
		Short: "Complete an item",
		Run: func(cmd *cobra.Command, args []string) {
			println(len(args))
			for _, v := range args {
				println("VX: done args are " + v)
			}
			if gid == "" {
				print("VX:TODO print to output: Error you didn't pass in a gid")
			}

			println("VX: TODO complete.", gid)
		},
	}
	doneCmd.Flags().StringVarP(&gid, "gid", "g", "", "The item to complete")
	return doneCmd
}
