package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"vixac.com/got/console"
)

var configPath string

var rootCmd = &cobra.Command{

	Use:   "Got",
	Short: "Got is a command line todo list",
	Long:  `Got is a comamnd line tool for managing todo list in a folder structure`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	printer := console.Printer{}
	// Global persistent flags
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "firbolg-ec2.yml", "Path to config file")
	done := buildDoneCommand(printer)
	rootCmd.AddCommand(done)
}
