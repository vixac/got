package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"vixac.com/got/console"
	"vixac.com/got/engine"
)

// VX:TODO test
func buildJobsCommand(deps RootDependencies) *cobra.Command {

	var underLookup string
	var jobsCmd = &cobra.Command{
		Use:   "jobs",
		Short: "fetch jobs under gid",
		Run: func(cmd *cobra.Command, args []string) {
			var lookup *engine.GidLookup = nil
			if underLookup != "" {

				realLookup := engine.GidLookup{
					Input: underLookup,
				}
				lookup = &realLookup
				fmt.Printf("Jobs lookup si %s\n", underLookup)
			}

			states := []int{engine.Active}

			res, err := deps.Engine.FetchItemsBelow(lookup, engine.AllDescendants, states)
			if err != nil {
				deps.Printer.Error(console.Message{Message: err.Error()})
				return
			}
			if res == nil || len(res.Result) == 0 {
				deps.Printer.Print(console.Message{Message: "no items found"})
				return
			}

			for _, v := range res.Result {
				var msg = ""
				if v.NumberGo != 0 {
					numStr := strconv.Itoa(v.NumberGo)
					msg += (numStr + "<GO>")
				}
				msg += ", item Id: "
				msg += v.Gid
				msg += "title: '"
				msg += v.Title
				msg += "'."

				deps.Printer.Print(console.Message{Message: msg})
			}

		},
	}
	jobsCmd.Flags().StringVarP(&underLookup, "under", "u", "", "The parent item")
	return jobsCmd

}
