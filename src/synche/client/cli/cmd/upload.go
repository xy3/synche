package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(uploadCmd)
}

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "<Short description here>",
	Long:  `<Long description here>`,
	Run: func(cmd *cobra.Command, args []string) {
		uploadJob()
	},
}

func uploadJob() {
	// retrieve a list of files that are on the server and return them here
	println("Upload command executed")
}