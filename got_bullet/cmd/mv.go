package cmd

import (
	"github.com/spf13/cobra"
	"vixac.com/got/console"
	"vixac.com/got/engine"
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
			if targetGid == "" {
				deps.Printer.Error(console.Message{Message: "missing target"})
				return
			}
			if newParentGid == "" {
				deps.Printer.Error(console.Message{Message: "missing parent"})
				return
			}
			oldParent, err := deps.Engine.Move(engine.GidLookup{Input: targetGid}, engine.GidLookup{Input: newParentGid})
			if err != nil {
				deps.Printer.Error(console.Message{Message: err.Error()})
				return
			}

			var msg string
			if oldParent == nil {
				msg = "Success: " + targetGid + " moved to new parent " + newParentGid
			} else {
				msg = "Success: " + targetGid + " moved from old parent '" + oldParent.Title + "' to " + newParentGid
			}
			deps.Printer.Print(console.Message{Message: msg})
		},
	}

	moveCmd.Flags().StringVarP(&targetGid, "gid", "g", "", "Target item")
	moveCmd.Flags().StringVarP(&newParentGid, "destination", "d", "", "Destination Parent")
	return moveCmd
}
