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
			if gid == "" {
				deps.Printer.Error(console.Message{Message: "Missing gid"})
				return
			}

			if alias == "" {
				deps.Printer.Error(console.Message{Message: "Missing alias"})
				return
			}
			_, err := deps.Engine.Alias(gid, alias)
			if err != nil {
				deps.Printer.Error(console.Message{Message: err.Error()})
				return
			}
			msg := "Success: " + gid + " is now aliased to " + alias + "."
			deps.Printer.Print(console.Message{Message: msg})
		},
	}
	cmd.Flags().StringVarP(&gid, "gid", "g", "", "The item to alias")
	cmd.Flags().StringVarP(&alias, "alias", "a", "", "The alias")
	return cmd
}
