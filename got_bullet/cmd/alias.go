package cmd

import (
	"github.com/spf13/cobra"
)

func buildAliasCommand(deps RootDependencies) *cobra.Command {
	var gid string
	var alias string
	var cmd = &cobra.Command{
		Use:   "alias",
		Short: "alias an item with a better name",
		Run: func(cmd *cobra.Command, args []string) {
			println(len(args))
			for _, v := range args {
				println("VX: done args are " + v)
			}
			if gid == "" {
				print("VX:TODO print to output: Error you didn't pass in a gid")
			}

			if alias == "" {
				print("VX:TODO print to output: Error you didn't pass in a gid")
			}

			println("VX: TODO complete.", gid)
		},
	}
	cmd.Flags().StringVarP(&gid, "gid", "g", "", "The item to alias")
	cmd.Flags().StringVarP(&alias, "alias", "a", "", "The alias")
	return cmd
}
