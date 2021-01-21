package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "<Short description here>",
	Long:  `<Long description here>`,
	Run: func(cmd *cobra.Command, args []string) {
		deleteJob()
	},
}

func deleteJob() {
	// retrieve a list of files that are on the server and return them here
}