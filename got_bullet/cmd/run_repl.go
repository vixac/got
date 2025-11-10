package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// you execute this to run got in repl mode
func buildReplCommand(deps RootDependencies, commands []*cobra.Command) *cobra.Command {

	var rootCmd = &cobra.Command{

		Use:   "Got Interface",
		Short: "Got Repl is a command line todo list",
		Long:  `Got Repl is a comamnd line tool for managing todo list in a folder structure`,
	}

	for _, c := range commands {
		rootCmd.AddCommand(c)
	}
	var cmd = &cobra.Command{
		Use:   "repl",
		Short: "run go in repl mode",
		Run: func(cmd *cobra.Command, args []string) {
			println(len(args))

			RunRepl(rootCmd)
		},
	}

	return cmd
}

func RunRepl(root *cobra.Command) {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Admin REPL (type 'help' or 'exit')")
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if line == "exit" {
			break
		}

		// Parse input into args
		args := strings.Fields(line)

		// Feed into cobra
		root.SetArgs(args)
		if err := root.Execute(); err != nil {
			fmt.Println("Error:", err)
		}
	}
}
