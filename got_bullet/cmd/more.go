package cmd

import (
	"github.com/spf13/cobra"
	"vixac.com/got/console"
)

// VX:TODO test
func buildMoreCommand(deps RootDependencies) *cobra.Command {

	var underGid string
	var jobsCmd = &cobra.Command{
		Use:   "more",
		Short: "Explain a gid",
		Run: func(cmd *cobra.Command, args []string) {
			if underGid == "" {
				print("VX: list jobs from the top lvl")
			}

			if underGid == "" {
				deps.Printer.Error(console.Message{Message: "no gid provided"})
				return
			}

			res, err := deps.Engine.Summary(nil)
			if err != nil {
				deps.Printer.Error(console.Message{Message: err.Error()})
				return
			}
			if res == nil {

				deps.Printer.Print(console.Message{Message: "no items found"})
				return
			}
			var msg = "Found job with gid: "
			msg += res.Gid
			msg += " , and title '"
			msg += res.Title
			msg += "'."
			deps.Printer.Print(console.Message{Message: msg})
		},
	}
	jobsCmd.Flags().StringVarP(&underGid, "under", "u", "", "The parent item")
	return jobsCmd

}
