package cmd

import (
	"github.com/spf13/cobra"
)

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

func init() {
	rootCmd.AddCommand(moveCmd)
	moveCmd.Flags().StringVarP(&targetGid, "gid", "g", "", "Target item")
	moveCmd.Flags().StringVarP(&newParentGid, "destination", "d", "", "Destination Parent")

}
