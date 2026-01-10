package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"vixac.com/got/console"
	"vixac.com/got/engine"
)

// VX:TODO test
func buildMoreCommand(deps RootDependencies) *cobra.Command {

	var jobsCmd = &cobra.Command{
		Use:   "more",
		Short: "Explain a gid",
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) < 1 {
				err := errors.New("missing args")
				deps.Printer.Error(console.Message{Message: err.Error()})
				return
			}
			gid := args[0]
			if gid == "" {
				deps.Printer.Error(console.Message{Message: "no gid provided"})
				return
			}

			lookup := engine.GidLookup{Input: gid}
			res, err := deps.Engine.Summary(&lookup)
			if err != nil {
				deps.Printer.Error(console.Message{Message: err.Error()})
				return
			}
			if res == nil {
				deps.Printer.Print(console.Message{Message: "no items found"})
				return
			}
			var msg = "Found job with gid: "
			msg += res.DisplayGid
			msg += " , and title '"
			msg += res.Title
			msg += "'."
			deps.Printer.Print(console.Message{Message: msg})
		},
	}
	return jobsCmd

}
