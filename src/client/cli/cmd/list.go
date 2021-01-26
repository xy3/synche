package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List files on the server",
	Long:  `<Long description here>`,
	Run: func(cmd *cobra.Command, args []string) {
		listJob()
	},
}

func listJob() {
	// retrieve a list of files that are on the server and return them here
	println("List command executed")
}