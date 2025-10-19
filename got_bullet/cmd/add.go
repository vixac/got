package cmd

import (
	"github.com/spf13/cobra"
)

var (
	underGid string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Create a new item",
	Run: func(cmd *cobra.Command, args []string) {
		println(len(args))
		for _, v := range args {
			println("VX: add args are " + v)
		}
		if underGid == "" {
			print("VX: todo create item at the top level")
		}

		println("VX: TODO ls items.", underGid)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&underGid, "under", "u", "", "The parent item")
}
