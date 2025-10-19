package cmd

import (
	"github.com/spf13/cobra"
)

func buildMvCommand(deps RootDependencies) *cobra.Command {
	var (
		targetGid    string
		newParentGid string
	)
	var moveCmd = &cobra.Command{
		Use:   "mv",
		Short: "move all items to a new gid",
		Run: func(cmd *cobra.Command, args []string) {
			println(len(args))
			for _, v := range args {
				println("VX: mv args are " + v)
			}

			println("VX: TODO mv items.")
		},
	}

	moveCmd.Flags().StringVarP(&targetGid, "gid", "g", "", "Target item")
	moveCmd.Flags().StringVarP(&newParentGid, "destination", "d", "", "Destination Parent")
	return moveCmd
}
