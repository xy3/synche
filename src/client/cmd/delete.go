package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a file on the server",
	Long:  `Sends a request to the server to delete a file by a specified file ID.`,
	Run: func(cmd *cobra.Command, args []string) {
		deleteJob()
	},
}

func deleteJob() {
	// retrieve a list of files that are on the server and return them here
	println("Delete command executed")
}
