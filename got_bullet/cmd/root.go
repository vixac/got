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

	rootCmd.AddCommand(buildDoneCommand(deps.Printer))
	rootCmd.AddCommand(buildJobsCommand(deps))
	rootCmd.AddCommand(buildMvCommand(deps))
	rootCmd.AddCommand(buildAliasCommand(deps))
	rootCmd.AddCommand(buildAddCommand(deps))
	return rootCmd
}
