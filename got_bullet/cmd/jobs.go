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
			//VX:TODO screw -u, lets got got jobs <id> optinoal.
			if len(args) != 0 {
				deps.Printer.Error(console.Message{Message: "jobs takes a -u and nothing else"})
				return
			}
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

			deps.Printer.Print(console.Message{Message: "-----------------------------------------\n\n"})

			for _, v := range res.Result {
				var msg = ""
				if v.NumberGo != 0 {
					numStr := strconv.Itoa(v.NumberGo)
					msg += (numStr + "<GO>")
				}
				msg += ", Path: "
				if v.Path != nil {
					for i, a := range v.Path.Ancestry {
						var id = a.Id
						if a.Alias != nil {
							id = *a.Alias
						}
						if i != 0 {
							msg += "->" + id
						} else {
							msg += id
						}

					}
				}

				msg += ", Gid = "
				if v.Alias != "" {
					msg += v.Alias + "("
				}
				msg += v.Gid
				if v.Alias != "" {
					msg += ")"
				}

				msg += ", Title = '"
				msg += v.Title
				msg += "'."

				deps.Printer.Print(console.Message{Message: msg})
			}

		},
	}
	jobsCmd.Flags().StringVarP(&underLookup, "under", "u", "", "The parent item")
	return jobsCmd

}
