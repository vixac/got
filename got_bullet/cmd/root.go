package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"vixac.com/got/console"
	"vixac.com/got/engine"
)

type RootDependencies struct {
	Printer console.Messenger
	Engine  engine.GotEngine
}

func Execute(deps RootDependencies) {
	rootCmd := NewRootCommand(deps)
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func NewRootCommand(deps RootDependencies) *cobra.Command {
	var rootCmd = &cobra.Command{

		Use:   "Got",
		Short: "Got is a command line todo list",
		Long:  `Got is a comamnd line tool for managing todo list in a folder structure`,
	}

	//these commands are passed into rootCmnd and also into repl, which is a root command too.
	var gotCommands []*cobra.Command

	gotCommands = append(gotCommands, buildCompleteCommand(deps))
	gotCommands = append(gotCommands, buildScheduleCommand(deps))
	gotCommands = append(gotCommands, buildTCommand(deps))
	gotCommands = append(gotCommands, buildEditCommand(deps))
	gotCommands = append(gotCommands, buildItCommand(deps))
	gotCommands = append(gotCommands, buildToCommand(deps))
	gotCommands = append(gotCommands, buildNoteCommand(deps))
	gotCommands = append(gotCommands, buildRemoveCommand(deps))

	gotCommands = append(gotCommands, buildEventCommand(deps))
	gotCommands = append(gotCommands, buildUnderCommand(deps))
	gotCommands = append(gotCommands, buildUnaliasCommand(deps))
	gotCommands = append(gotCommands, buildDoneCommand(deps))
	gotCommands = append(gotCommands, buildAliasCommand(deps))
	gotCommands = append(gotCommands, buildTillCommand(deps))

	gotCommands = append(gotCommands, buildMoreCommand(deps))

	gotCommands = append(gotCommands, buildMvCommand(deps))

	gotCommands = append(gotCommands, buildJobsCommand(deps))
	rootCmd.AddCommand(buildReplCommand(deps, gotCommands))
	for _, c := range gotCommands {
		rootCmd.AddCommand(c)
	}
	return rootCmd
}
