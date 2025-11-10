package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

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
