package cmd

import (
	"github.com/spf13/cobra"
	"vixac.com/got/console"
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
				deps.Printer.Error(console.Message{Message: "Missing gid"})
				return
			}

			if alias == "" {
				deps.Printer.Error(console.Message{Message: "Missing alias"})
				return
			}

			println("VX: TODO complete.", gid)
			ok, err := deps.Engine.Alias(gid, alias)
			if err != nil {
				println("VX: error aliasing: ", err.Error())
			}
			print("VX: ok was ", ok)
		},
	}
	cmd.Flags().StringVarP(&gid, "gid", "g", "", "The item to alias")
	cmd.Flags().StringVarP(&alias, "alias", "a", "", "The alias")
	return cmd
}
