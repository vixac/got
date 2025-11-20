package cmd

import (
	"github.com/spf13/cobra"
	"vixac.com/got/console"
)

func buildJobsCommand(deps RootDependencies) *cobra.Command {

	var underGid string
	var jobsCmd = &cobra.Command{
		Use:   "jobs",
		Short: "List jobs",
		Run: func(cmd *cobra.Command, args []string) {
			println("jobs took arg count:", len(args))
			for _, v := range args {
				println("VX: jobs args are " + v)
			}
			if underGid == "" {
				print("VX: list jobs from the top lvl")
			}

			println("VX: TODO ls items.", underGid)
			//VX:TODO lookup optional?

			res, err := deps.Engine.Summary(nil)
			if err != nil {
				deps.Printer.Error(console.Message{Message: err.Error()})
			}
			print("VX: title is " + res.Title)
		},
	}
	jobsCmd.Flags().StringVarP(&underGid, "under", "u", "", "The parent item")
	return jobsCmd

}
