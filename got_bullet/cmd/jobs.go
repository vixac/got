package cmd

import (
	"github.com/spf13/cobra"
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
		},
	}
	jobsCmd.Flags().StringVarP(&underGid, "under", "u", "", "The parent item")
	return jobsCmd

}
